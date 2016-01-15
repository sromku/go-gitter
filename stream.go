package gitter
import (
	"bufio"
	"encoding/json"
)

// Initialize stream
func (gitter *Gitter) Stream(roomId string) *Stream {
	return &Stream{
		url: GITTER_STREAM_API + "rooms/" + roomId + "/chatMessages",
		GitterMessage: make(chan Message),
	}
}

// Start streaming api and listen to incoming messages
func (gitter *Gitter) Listen(stream *Stream) {
	res, _ := gitter.getResponse(stream.url)
	reader := bufio.NewReader(res.Body)
	var gitterMessage Message
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			gitter.log(err)
		}
		err = json.Unmarshal(line, &gitterMessage)
		if err == nil {
			stream.GitterMessage <- gitterMessage
		} else {
			gitter.log(err)
		}
	}
}

// Definition of stream
type Stream struct {
	url           string
	GitterMessage chan Message
}
