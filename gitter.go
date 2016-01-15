// Gitter API in Go.
//
// Author: sromku
package gitter

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strconv"
	"bytes"
	"io"
	"log"
	"time"
)

var GITTER_REST_API string = "https://api.gitter.im/v1/"
var GITTER_STREAM_API string = "https://stream.gitter.im/v1/"

type Gitter struct {
	config struct {
		       token  string
		       client *http.Client
	       }
	debug  bool
	logWriter io.Writer
}

// Initialize Gitter API
//
// For example:
//  api := gitter.New("YOUR_ACCESS_TOKEN")
func New(token string) *Gitter {
	s := &Gitter{}
	s.config.token = token
	s.config.client = &http.Client{}
	return s
}

// Set your own http client. Can be useful in App Engine case.
func (gitter *Gitter) SetClient(client *http.Client) {
	gitter.config.client = client
}

// Get current user
func (gitter *Gitter) GetUser() (*User, error) {

	var users []User
	response, err := gitter.get(GITTER_REST_API + "user")
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

	err = GitterApiError{What:"Failed to retrieve current user"}
	gitter.log(err)
	return nil, err
}

// List of Rooms the user is part of
func (gitter *Gitter) GetUserRooms(userId string) ([]Room, error) {

	var rooms []Room
	response, err := gitter.get(GITTER_REST_API + "user/" + userId + "/rooms")
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

// List rooms the current user is in
func (gitter *Gitter) GetRooms() ([]Room, error) {

	var rooms []Room
	response, err := gitter.get(GITTER_REST_API + "rooms")
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

// Get room by id
func (gitter *Gitter) GetRoom(roomId string) (*Room, error) {

	var room Room
	response, err := gitter.get(GITTER_REST_API + "rooms/" + roomId)
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

// List of messages in a room.
// Pagination is optional. You can pass nil or specific pagination params.
func (gitter *Gitter) GetMessages(roomId string, params *Pagination) ([]Message, error) {

	var messages []Message
	url := GITTER_REST_API + "rooms/" + roomId + "/chatMessages"
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

// Get message in a room by message id.
func (gitter *Gitter) GetMessage(roomId, messageId string) (*Message, error) {

	var message Message
	response, err := gitter.get(GITTER_REST_API + "rooms/" + roomId + "/chatMessages/" + messageId)
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

// Send a message to a room
func (gitter *Gitter) SendMessage(roomId, text string) error {

	message := Message{Text:text}
	body, _ := json.Marshal(message)
	err := gitter.post(GITTER_REST_API + "rooms/" + roomId + "/chatMessages", body)
	if err != nil {
		gitter.log(err)
		return err
	}

	return nil
}

// Set true if you want to trace errors
func (gitter *Gitter) SetDebug(debug bool, logWriter io.Writer) {
	gitter.debug = debug
	gitter.logWriter = logWriter
}

// Pagination params
type Pagination struct {

	// Skip n messages
	Skip     int

	// Get messages before beforeId
	BeforeId string

	// Get messages after afterId
	AfterId  string

	// Maximum number of messages to return
	Limit    int
}

func (messageParams *Pagination) encode() string {
	values := url.Values{}

	if messageParams.AfterId != "" {
		values.Add("afterId", messageParams.AfterId)
	}

	if messageParams.BeforeId != "" {
		values.Add("beforeId", messageParams.BeforeId)
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
	r.Header.Set("Authorization", "Bearer " + gitter.config.token)
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
		err = GitterApiError{What: fmt.Sprintf("Status code: %v", resp.StatusCode)}
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
	r.Header.Set("Authorization", "Bearer " + gitter.config.token)

	resp, err := gitter.config.client.Do(r)
	if err != nil {
		gitter.log(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		err = GitterApiError{What: fmt.Sprintf("Status code: %v", resp.StatusCode)}
		gitter.log(err)
		return err
	}

	return nil
}

func (gitter *Gitter) log(a interface{}) {
	if gitter.debug {
		if gitter.logWriter == nil {
			log.Println(a)
		} else {
			timestamp := time.Now().Format(time.RFC3339)
			msg := fmt.Sprintf("%v: %v", timestamp, a)
			fmt.Fprintln(gitter.logWriter, msg)
		}
	}
}

type GitterApiError struct {
	What string
}

func (e GitterApiError) Error() string {
	return fmt.Sprintf("%v", e.What)
}