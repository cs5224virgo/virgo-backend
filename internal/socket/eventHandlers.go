package socket

import "github.com/cs5224virgo/virgo/logger"

func (h *WebSocketHub) handleResumeRoom(event *ResumeRoomEventReq, client *Client) {
	h.roomsMutex.Lock()
	if h.rooms[event.RoomCode] == nil {
		h.rooms[event.RoomCode] = make(map[*Client]bool)
	}
	h.rooms[event.RoomCode][client] = true
	client.roomCodes = append(client.roomCodes, event.RoomCode)
	h.roomsMutex.Unlock()
}

func (h *WebSocketHub) handleUpdateUnread(event *UpdateUnreadReq, client *Client) {
	err := h.DataLayer.SetUnreadCount(client.username, event.RoomCode, event.Unread)
	if err != nil {
		logger.Errorf("Error updating unread: %v", err)
	}
}

func (h *WebSocketHub) handleNewMessage(event *NewMessageEventReq, client *Client) {
	message, room, err := h.DataLayer.NewMessage(event.Content, event.RoomCode, event.Username)
	if err != nil {
		logger.Error(err)
	}
	// announce to everyone in the room
	h.roomsMutex.Lock()
	if h.rooms[event.RoomCode] != nil {
		for client := range h.rooms[event.RoomCode] {
			data := NewMessageEventResp{
				Message: message,
				Room:    room,
			}
			resp := Event{
				EventType: EventTypeNewMessageResp,
				Data:      data,
			}
			client.send <- resp
		}
	}
	h.roomsMutex.Unlock()
}

func (h *WebSocketHub) handleEvent(event Event, client *Client) {
	switch event.EventType {
	case EventTypeResumeRoomReq:
		data := event.Data.(*ResumeRoomEventReq)
		h.handleResumeRoom(data, client)
	case EventTypeUpdateUnreadReq:
		data := event.Data.(*UpdateUnreadReq)
		h.handleUpdateUnread(data, client)
	case EventTypeNewMessageReq:
		data := event.Data.(*NewMessageEventReq)
		h.handleNewMessage(data, client)
	default:
		logger.Error("unrecognized event: " + event.EventType)
	}
}
