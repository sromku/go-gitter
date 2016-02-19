// Package gitter is a Go client library for the Gitter API.
//
// Author: sromku
package gitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mreiferson/go-httpclient"
)

var gitterRESTAPI = "https://api.gitter.im/v1/"
var gitterStreamAPI = "https://stream.gitter.im/v1/"

type Gitter struct {
	config struct {
		token  string
		client *http.Client
	}
	debug     bool
	logWriter io.Writer
}

// New initializez the Gitter API client
//
// For example:
//  api := gitter.New("YOUR_ACCESS_TOKEN")
func New(token string) *Gitter {

	transport := &httpclient.Transport{
		ConnectTimeout:   5 * time.Second,
		ReadWriteTimeout: 40 * time.Second,
	}
	defer transport.Close()

	s := &Gitter{}
	s.config.token = token
	s.config.client = &http.Client{
		Transport: transport,
	}
	return s
}

// SetClient sets a custom http client. Can be useful in App Engine case.
func (gitter *Gitter) SetClient(client *http.Client) {
	gitter.config.client = client
}

// GetUser returns the current user
func (gitter *Gitter) GetUser() (*User, error) {

	var users []User
	response, err := gitter.get(gitterRESTAPI + "user")
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &users)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	if len(users) > 0 {
		return &users[0], nil
	}

	err = APIError{What: "Failed to retrieve current user"}
	gitter.log(err)
	return nil, err
}

// GetUserRooms returns a list of Rooms the user is part of
func (gitter *Gitter) GetUserRooms(userID string) ([]Room, error) {

	var rooms []Room
	response, err := gitter.get(gitterRESTAPI + "user/" + userID + "/rooms")
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &rooms)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	return rooms, nil
}

// GetRooms returns a list of rooms the current user is in
func (gitter *Gitter) GetRooms() ([]Room, error) {

	var rooms []Room
	response, err := gitter.get(gitterRESTAPI + "rooms")
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &rooms)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	return rooms, nil
}

// GetRoom returns a room with the passed id
func (gitter *Gitter) GetRoom(roomID string) (*Room, error) {

	var room Room
	response, err := gitter.get(gitterRESTAPI + "rooms/" + roomID)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &room)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	return &room, nil
}

// GetMessages returns a list of messages in a room.
// Pagination is optional. You can pass nil or specific pagination params.
func (gitter *Gitter) GetMessages(roomID string, params *Pagination) ([]Message, error) {

	var messages []Message
	url := gitterRESTAPI + "rooms/" + roomID + "/chatMessages"
	if params != nil {
		url += "?" + params.encode()
	}
	response, err := gitter.get(url)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &messages)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	return messages, nil
}

// GetMessage returns a message in a room.
func (gitter *Gitter) GetMessage(roomID, messageID string) (*Message, error) {

	var message Message
	response, err := gitter.get(gitterRESTAPI + "rooms/" + roomID + "/chatMessages/" + messageID)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	err = json.Unmarshal(response, &message)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	return &message, nil
}

// SendMessage sends a message to a room
func (gitter *Gitter) SendMessage(roomID, text string) error {

	message := Message{Text: text}
	body, _ := json.Marshal(message)
	err := gitter.post(gitterRESTAPI+"rooms/"+roomID+"/chatMessages", body)
	if err != nil {
		gitter.log(err)
		return err
	}

	return nil
}

// SetDebug traces errors if it's set to true.
func (gitter *Gitter) SetDebug(debug bool, logWriter io.Writer) {
	gitter.debug = debug
	gitter.logWriter = logWriter
}

// Pagination params
type Pagination struct {

	// Skip n messages
	Skip int

	// Get messages before beforeId
	BeforeID string

	// Get messages after afterId
	AfterID string

	// Maximum number of messages to return
	Limit int
}

func (messageParams *Pagination) encode() string {
	values := url.Values{}

	if messageParams.AfterID != "" {
		values.Add("afterId", messageParams.AfterID)
	}

	if messageParams.BeforeID != "" {
		values.Add("beforeId", messageParams.BeforeID)
	}

	if messageParams.Skip > 0 {
		values.Add("skip", strconv.Itoa(messageParams.Skip))
	}

	if messageParams.Limit > 0 {
		values.Add("limit", strconv.Itoa(messageParams.Limit))
	}

	return values.Encode()
}

func (gitter *Gitter) getResponse(url string) (*http.Response, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		gitter.log(err)
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+gitter.config.token)
	response, err := gitter.config.client.Do(r)
	if err != nil {
		gitter.log(err)
		return nil, err
	}
	return response, nil
}

func (gitter *Gitter) get(url string) ([]byte, error) {
	resp, err := gitter.getResponse(url)
	if err != nil {
		gitter.log(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		err = APIError{What: fmt.Sprintf("Status code: %v", resp.StatusCode)}
		gitter.log(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gitter.log(err)
		return nil, err
	}

	return body, nil
}

func (gitter *Gitter) post(url string, body []byte) error {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		gitter.log(err)
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+gitter.config.token)

	resp, err := gitter.config.client.Do(r)
	if err != nil {
		gitter.log(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		err = APIError{What: fmt.Sprintf("Status code: %v", resp.StatusCode)}
		gitter.log(err)
		return err
	}

	return nil
}

func (gitter *Gitter) log(a interface{}) {
	if gitter.debug {
		log.Println(a)
		if gitter.logWriter != nil {
			timestamp := time.Now().Format(time.RFC3339)
			msg := fmt.Sprintf("%v: %v", timestamp, a)
			fmt.Fprintln(gitter.logWriter, msg)
		}
	}
}

// APIError holds data of errors returned from the API.
type APIError struct {
	What string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.What)
}
