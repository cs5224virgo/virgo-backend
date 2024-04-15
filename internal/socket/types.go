package socket

import (
	"encoding/json"
)

type DataLayer interface {
	SetUnreadCount(username string, roomCode string, unread int) error
}

type EventType string

const (
	EventTypeJoinRoomReq     EventType = "req-join-room"
	EventTypeJoinRoomResp    EventType = "resp-join-room"
	EventTypeUpdateUnreadReq EventType = "req-update-unread"
)

type Event struct {
	EventType EventType `json:"eventType"`
	Data      any       `json:"data"`
}

type JoinRoomEventResp struct {
	RoomCode string `json:"roomCode"`
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
	case EventTypeJoinRoomReq:
		e.Data = new(JoinRoomEventReq)
	case EventTypeJoinRoomResp:
		e.Data = new(JoinRoomEventResp)
	case EventTypeUpdateUnreadReq:
		e.Data = new(UpdateUnreadReq)
	}

	type eventAlias Event
	return json.Unmarshal(data, (*eventAlias)(e))
}
