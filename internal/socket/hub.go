package socket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// WebSocketHub maintains the set of active clients and broadcasts messages to the
// clients.
type WebSocketHub struct {
	DataLayer DataLayer

	// Registered clients.
	activeClients map[*Client]bool

	// Registered clients, sorted by the rooms they're in
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	// broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == viper.GetString("frontend_url")
	},
}

func NewWebSocketHub(datalayer DataLayer) *WebSocketHub {
	return &WebSocketHub{
		DataLayer:     datalayer,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		rooms:         make(map[string]map[*Client]bool),
		activeClients: make(map[*Client]bool),
		// broadcast:  make(chan Message),
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
	h.activeClients[client] = true
}

func (h *WebSocketHub) removeClient(client *Client) {
	delete(h.activeClients, client)
	for _, roomCode := range client.roomCodes {
		if _, ok := h.rooms[roomCode]; ok {
			delete(h.rooms[roomCode], client)
		}
	}
	close(client.send)
}

// ServeWs handles websocket requests from the peer.
func (h *WebSocketHub) ServeWs(c *gin.Context, username string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(username, conn, h)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
