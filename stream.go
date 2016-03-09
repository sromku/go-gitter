package gitter

import (
	"bufio"
	"encoding/json"
	"net/http"
	"time"
)

var defaultConnectionWaitTime time.Duration = 3000 // millis
var defaultConnectionMaxRetries = 5

// Stream initialize stream
func (gitter *Gitter) Stream(roomID string) *Stream {
	return &Stream{
		url:    gitterStreamAPI + "rooms/" + roomID + "/chatMessages",
		Event:  make(chan Event),
		gitter: gitter,
		streamConnection: gitter.newStreamConnection(
			defaultConnectionWaitTime,
			defaultConnectionMaxRetries),
	}
}

func (gitter *Gitter) Listen(stream *Stream) {

	defer stream.destroy()

	var reader *bufio.Reader
	var gitterMessage Message

	// connect
	stream.connect()

Loop:
	for {

		// if closed then stop trying
		if stream.isClosed() {
			stream.Event <- Event{
				Data: &GitterConnectionClosed{},
			}
			break Loop
		}

		reader = bufio.NewReader(stream.getResponse().Body)
		line, err := reader.ReadBytes('\n')
		if stream.isClosed() {
			continue
		} else if err != nil {
			stream.connect()
			continue
		}

		// unmarshal the streamed data
		err = json.Unmarshal(line, &gitterMessage)
		if err != nil {
			gitter.log(err)
			continue
		}

		// we are here, then we got the good message. pipe it forward.
		stream.Event <- Event{
			Data: &MessageReceived{
				Message: gitterMessage,
			},
		}
	}

	gitter.log("Listening was completed")
}

// Stream holds stream data.
type Stream struct {
	url              string
	Event            chan Event
	streamConnection *streamConnection
	gitter           *Gitter
}

func (stream *Stream) destroy() {
	close(stream.Event)
}

type Event struct {
	Data interface{}
}

type GitterConnectionClosed struct {
}

type MessageReceived struct {
	Message Message
}

// connect and try to reconnect with
func (stream *Stream) connect() {

	if stream.streamConnection.retries == stream.streamConnection.currentRetries {
		stream.Close()
		stream.gitter.log("Number of retries exceeded the max retries number, we are done here")
		return
	}

	res, err := stream.gitter.getResponse(stream.url)
	if err != nil || res.StatusCode != 200 {
		stream.gitter.log("Failed to get response, trying reconnect")
		stream.gitter.log(err)

		// sleep and wait
		stream.streamConnection.currentRetries++
		time.Sleep(time.Millisecond * stream.streamConnection.wait * time.Duration(stream.streamConnection.currentRetries))

		// connect again
		stream.Close()
		stream.connect()

	} else {
		stream.gitter.log("Response was received")
		stream.streamConnection.currentRetries = 0
		stream.streamConnection.closed = false
		stream.streamConnection.response = res
	}
}

type streamConnection struct {

	// connection was closed
	closed bool

	// wait time till next try
	wait time.Duration

	// max tries to recover
	retries int

	// current streamed response
	response *http.Response

	// current status
	currentRetries int
}

// Close the stream connection and stop receiving streamed data
func (stream *Stream) Close() {
	conn := stream.streamConnection
	conn.closed = true
	if conn.response != nil {
		stream.gitter.log("Stream connection was closed")
		defer conn.response.Body.Close()
	}
	conn.currentRetries = 0
}

func (stream *Stream) isClosed() bool {
	return stream.streamConnection.closed
}

func (stream *Stream) getResponse() *http.Response {
	return stream.streamConnection.response
}

// Optional, set stream connection properties
// wait - time in milliseconds of waiting between reconnections. Will grow exponentially.
// retries - number of reconnections retries before dropping the stream.
func (gitter *Gitter) newStreamConnection(wait time.Duration, retries int) *streamConnection {
	return &streamConnection{
		closed:  true,
		wait:    wait,
		retries: retries,
	}
}
