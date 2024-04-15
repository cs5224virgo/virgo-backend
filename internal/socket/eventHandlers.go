package socket

import "github.com/cs5224virgo/virgo/logger"

func (h *WebSocketHub) handleJoinRoom(event JoinRoomEventReq, client *Client) {
	if h.rooms[event.RoomCode] == nil {
		h.rooms[event.RoomCode] = make(map[*Client]bool)
	}
	h.rooms[event.RoomCode][client] = true
	client.roomCodes = append(client.roomCodes, event.RoomCode)

	data := JoinRoomEventResp{
		RoomCode: "test",
	}
	resp := Event{
		EventType: EventTypeJoinRoomResp,
		Data:      data,
	}
	client.send <- resp
}

func (h *WebSocketHub) handleUpdateUnread(event UpdateUnreadReq, client *Client) {
	err := h.DataLayer.SetUnreadCount(client.username, event.RoomCode, event.Unread)
	if err != nil {
		logger.Errorf("Error updating unread: %v", err)
	}
}

func (h *WebSocketHub) handleEvent(event Event, client *Client) {
	switch event.EventType {
	case EventTypeJoinRoomReq:
		data := event.Data.(JoinRoomEventReq)
		h.handleJoinRoom(data, client)
	case EventTypeUpdateUnreadReq:
		data := event.Data.(UpdateUnreadReq)
		h.handleUpdateUnread(data, client)
	}
}
