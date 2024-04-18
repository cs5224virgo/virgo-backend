package socket

import (
	"encoding/json"

	"github.com/cs5224virgo/virgo/internal/datalayer"
)

type DataLayer interface {
	SetUnreadCount(username string, roomCode string, unread int) error
	NewMessage(content string, roomCode string, username string) (datalayer.DetailedMessage, datalayer.DetailedRoom, error)
}

type EventType string

const (
	EventTypeResumeRoomReq   EventType = "req-resume-room"
	EventTypeJoinRoomReq     EventType = "req-join-room"
	EventTypeJoinRoomResp    EventType = "resp-join-room"
	EventTypeLeaveRoomResp   EventType = "resp-leave-room"
	EventTypeNewRoomResp     EventType = "resp-new-room"
	EventTypeUpdateUnreadReq EventType = "req-update-unread"
	EventTypeNewMessageReq   EventType = "req-new-message"
	EventTypeNewMessageResp  EventType = "resp-new-message"
)

type Event struct {
	EventType EventType `json:"eventType"`
	Data      any       `json:"data"`
}

type ResumeRoomEventReq struct {
	Username string `json:"username"`
	RoomCode string `json:"roomCode"`
}

type JoinRoomEventResp struct {
	User     datalayer.DetailedUser `json:"user"`
	RoomCode string                 `json:"roomCode"`
}

type NewRoomEventResp struct {
	Room datalayer.DetailedRoom `json:"room"`
}

type LeaveRoomEventResp struct {
	User     datalayer.DetailedUser `json:"user"`
	RoomCode string                 `json:"roomCode"`
}

type NewMessageEventResp struct {
	Message datalayer.DetailedMessage `json:"message"`
	Room    datalayer.DetailedRoom    `json:"room"`
}

type JoinRoomEventReq struct {
	Username string `json:"username"`
	RoomCode string `json:"roomCode"`
}

type UpdateUnreadReq struct {
	Username string `json:"username"`
	RoomCode string `json:"roomCode"`
	Unread   int    `json:"unread"`
}

type NewMessageEventReq struct {
	Content  string `json:"content"`
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

func (e *Event) UnmarshalJSON(data []byte) error {
	var inner struct {
		EventType EventType `json:"eventType"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}

	switch inner.EventType {
	case EventTypeResumeRoomReq:
		e.Data = new(ResumeRoomEventReq)
	case EventTypeJoinRoomReq:
		e.Data = new(JoinRoomEventReq)
	case EventTypeUpdateUnreadReq:
		e.Data = new(UpdateUnreadReq)
	case EventTypeNewMessageReq:
		e.Data = new(NewMessageEventReq)
	}

	type eventAlias Event
	return json.Unmarshal(data, (*eventAlias)(e))
}
