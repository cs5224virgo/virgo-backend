package socket

type EventType string

const (
	EventTypeJoinRoom     EventType = "join-room"
	EventTypeUpdateUnread EventType = "update-unread"
)

type Event struct {
	EventType EventType   `json:"eventType"`
	Data      interface{} `json:"data"`
}

type JoinRoomEventResp struct {
	RoomCode string `json:"roomCode"`
}

type JoinRoomEventReq struct {
	Username string `json:"username"`
	RoomCode string `json:"roomCode"`
}

// Message struct to hold message data
type Message struct {
	Type     string `json:"type"`
	RoomCode string `json:"roomCode"`
	Username string `json:"username"`
	Content  string `json:"content"`
}
