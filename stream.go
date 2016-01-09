package gitter
import (
	"bufio"
	"encoding/json"
	"log"
	"time"
)

// Initialize stream
func (room *Room) Stream() *Stream {
	return &Stream{
		Room: room,
		Url: GITTER_STREAM_API + "rooms/" + room.Id + "/chatMessages",
		GitterMessage: make(chan GitterMessage),
	}
}

// Start streaming api and listen to incoming messages
func (gitter *Gitter) Listen(stream *Stream) {
	res, _ := gitter.getResponse(stream.Url)
	reader := bufio.NewReader(res.Body)
	var gitterMessage GitterMessage
	for {
		line, _ := reader.ReadBytes('\n')
		err := json.Unmarshal(line, &gitterMessage)
		if err == nil {
			stream.GitterMessage <- gitterMessage
		} else if gitter.debug {
			log.Println(err)
		}
	}
}

type GitterMessage struct {

	// message id
	Id string `json:"id"`

	// message body
	Text string `json:"text"`

	// sent time
	Sent time.Time `json:"sent"`

	// from user
	From User `json:"fromUser"`

}

// Definition of stream
type Stream struct {
	Room          *Room
	Url           string
	GitterMessage chan GitterMessage
}
