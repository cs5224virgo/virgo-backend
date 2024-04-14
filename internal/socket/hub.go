package socket

import (
	"log"

	"github.com/gin-gonic/gin"
)

// WebSocketHub maintains the set of active clients and broadcasts messages to the
// clients.
type WebSocketHub struct {
	// Registered clients.
	// clients map[*Client]bool
	allClients map[*Client]bool

	// Registered clients, sorted by the rooms they're in
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		allClients: make(map[*Client]bool),
	}
}

func (h *WebSocketHub) Run() {
	for {
		select {
		// Register a client
		case client := <-h.register:
			h.registerNewClient(client)
		// Unregister a client
		case client := <-h.unregister:
			h.removeClient(client)
			// Send a message
			// case message := <-h.broadcast:
			// 	h.handleMessage(message)
		}
	}
}

func (h *WebSocketHub) registerNewClient(client *Client) {
	h.allClients[client] = true
	for _, roomCode := range client.roomCodes {
		if h.rooms[roomCode] == nil {
			h.rooms[roomCode] = make(map[*Client]bool)
		}
		h.rooms[roomCode][client] = true
	}
}

func (h *WebSocketHub) removeClient(client *Client) {
	delete(h.allClients, client)
	for _, roomCode := range client.roomCodes {
		if _, ok := h.rooms[roomCode]; ok {
			delete(h.rooms[roomCode], client)
		}
	}
	close(client.send)
}

func (h *WebSocketHub) handleMessage(message Message) {
	for client := range h.rooms[message.RoomCode] {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.rooms[message.RoomCode], client)
			delete(h.allClients, client)
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (h *WebSocketHub) ServeWs(c *gin.Context, username string, roomCodes []string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(username, roomCodes, conn, h)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
