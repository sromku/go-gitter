package gitter
import (
	"bufio"
	"encoding/json"
	"log"
)

func (room *Room) Stream() *Stream {
	return &Stream{
		Room: room,
		GitterMessage: make(chan GitterMessage),
	}
}

func (gitter *Gitter) Listen(stream *Stream) {
	res, _ := gitter.getResponse(GITTER_STREAM_API + "rooms/" + stream.Room.Id + "/chatMessages")
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
	Id string `json:"id"`
}

type Stream struct {
	Room          *Room
	GitterMessage chan GitterMessage
}
