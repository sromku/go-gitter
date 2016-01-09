# gitter
Gitter API in Go

#### Install

`go get github.com/sromku/gitter`

- [Initialize](#initialize)
- [Users](#users)
- [Rooms](#rooms)
- [Messages](#messages)
- [App Engine](#app-engine)

##### Initialize 
``` Go
api := gitter.New("YOUR_ACCESS_TOKEN")
```

##### Users

- Get user
	``` Go
	user, err := api.GetUser()
	```

##### Rooms

- Get all rooms
	``` Go
	rooms, err := api.GetRooms()
	```

- Get room by id
	``` Go
	room, err := api.GetRoom("roomId")
	```

- Get rooms of some user
	``` Go
	rooms, err := api.GetRooms("userId")
	```

##### Messages

- Get messages of room
	``` Go
	messages, err := api.GetMessages("roomId", nil)
	```

- Get one message
	``` Go
	message, err := api.GetMessage("roomId", "messageId")
	```

- Send message
	``` Go
	err := api.SendMessage("roomId", "free chat text")
	```

##### App Engine

Initialize and continue as usual

``` Go
c := appengine.NewContext(r)
client := urlfetch.Client(c)

api := gitter.New("YOUR_ACCESS_TOKEN")
api.SetClient(client)
```

[Documentation](https://godoc.org/github.com/sromku/gitter)
