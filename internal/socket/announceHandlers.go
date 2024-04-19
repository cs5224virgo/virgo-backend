package socket

import "github.com/cs5224virgo/virgo/internal/datalayer"

func (h *WebSocketHub) AnnounceAddUserToRoom(user datalayer.DetailedUser, room datalayer.DetailedRoom, isNewMember bool) {
	roomCode := room.Code
	h.roomsMutex.Lock()

	// create room if not exist
	if h.rooms[roomCode] == nil {
		h.rooms[roomCode] = make(map[*Client]bool)
	}

	// add this user to the room if not already in
	if h.rooms[roomCode] != nil {
		for client := range h.activeClients {
			if client.username == user.Username {
				if _, ok := h.rooms[roomCode][client]; !ok {
					client.roomCodes = append(client.roomCodes, roomCode)
					h.rooms[roomCode][client] = true
				}
			}
		}
	}

	// announce to everyone in the room
	if h.rooms[roomCode] != nil {
		for client := range h.rooms[roomCode] {
			// if isNewMember, send a NewRoomEventResp
			if client.username == user.Username && isNewMember {
				data := NewRoomEventResp{
					Room: room,
				}
				resp := Event{
					EventType: EventTypeNewRoomResp,
					Data:      data,
				}
				client.send <- resp
			} else {
				data := JoinRoomEventResp{
					RoomCode: roomCode,
					User:     user,
				}
				resp := Event{
					EventType: EventTypeJoinRoomResp,
					Data:      data,
				}
				client.send <- resp
			}
		}
	}

	// h.debugRooms()
	h.roomsMutex.Unlock()
}

func (h *WebSocketHub) AnnounceLeaveRoom(user datalayer.DetailedUser, roomCode string) {
	h.roomsMutex.Lock()

	// remove this user from the room
	if h.rooms[roomCode] != nil {
		for client := range h.activeClients {
			if client.username == user.Username {
				delete(h.rooms[roomCode], client)
			}
		}
	}

	// announce to everyone else in the room
	if h.rooms[roomCode] != nil {
		for client := range h.rooms[roomCode] {
			data := LeaveRoomEventResp{
				User:     user,
				RoomCode: roomCode,
			}
			resp := Event{
				EventType: EventTypeLeaveRoomResp,
				Data:      data,
			}
			client.send <- resp
		}
	}

	// h.debugRooms()
	h.roomsMutex.Unlock()
}
