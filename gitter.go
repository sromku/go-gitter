package gitter

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var GITTER_API string = "https://api.gitter.im/v1/"

func New(token string) *Gitter {
	s := &Gitter{}
	s.config.token = token
	s.config.client = &http.Client{}
	return s
}

type Gitter struct {
	config struct {
		       token string
		       client  *http.Client
	       }
	debug  bool
}

// Set your own http client. Can be useful in App Engine case.
func (gitter *Gitter) SetClient(client *http.Client) {
	gitter.config.client = client
}

// List rooms the current user is in
func (gitter *Gitter) GetRooms() ([]Room, error) {
	var rooms []Room
	response, err := gitter.get(GITTER_API + "rooms")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (gitter *Gitter) get(url string) ([]byte, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer " + gitter.config.token)
	resp, err := gitter.config.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, GitterApiError{What: fmt.Sprintf("Status code: %v", resp.StatusCode) }
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type GitterApiError struct {
	What string
}

func (e GitterApiError) Error() string {
	return fmt.Sprintf("%v", e.What)
}