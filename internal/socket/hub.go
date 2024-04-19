package socket

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cs5224virgo/virgo/logger"
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

	// mutex
	roomsMutex sync.Mutex

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
		roomsMutex:    sync.Mutex{},
		// broadcast:  make(chan Message),
	}
}

func (h *WebSocketHub) Run() {
	go h.cleanRooms()
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
	h.roomsMutex.Lock()
	// logger.Info("registering client " + client.username)
	h.activeClients[client] = true
	// h.debugRooms()
	h.roomsMutex.Unlock()
}

func (h *WebSocketHub) removeClient(client *Client) {
	// logger.Info("removing client " + client.username)
	h.roomsMutex.Lock()
	delete(h.activeClients, client)
	for _, roomCode := range client.roomCodes {
		if _, ok := h.rooms[roomCode]; ok {
			delete(h.rooms[roomCode], client)
		}
	}
	close(client.send)
	// h.debugRooms()
	h.roomsMutex.Unlock()
}

// ServeWs handles websocket requests from the peer.
func (h *WebSocketHub) ServeWs(c *gin.Context, username string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(username, conn, h)
	h.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.read()
	go client.write()
}

func (h *WebSocketHub) cleanRooms() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		h.roomsMutex.Lock()
		for roomCode := range h.rooms {
			hasClient := false
			for client := range h.activeClients {
				for _, clientRoomCode := range client.roomCodes {
					if clientRoomCode == roomCode {
						hasClient = true
					}
				}
			}
			if !hasClient && len(h.rooms[roomCode]) == 0 {
				delete(h.rooms, roomCode)
			}
		}
		// h.debugRooms()
		h.roomsMutex.Unlock()
	}
}

func (h *WebSocketHub) debugRooms() {
	for roomCode := range h.rooms {
		for client := range h.rooms[roomCode] {
			logger.Info("room " + roomCode + " has user " + client.username + " with ID " + client.clientID)
		}
	}
	for client := range h.activeClients {
		logger.Info("user " + client.username + " is connected with client " + client.clientID)
	}
}
