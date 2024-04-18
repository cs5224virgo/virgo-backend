package socket

import (
	"time"

	"github.com/cs5224virgo/virgo/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxEventSize = 1024

	// maximum amount of message held in memory
	maxEventBuffer = 32
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	clientID  string
	username  string
	roomCodes []string
	hub       *WebSocketHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Event
}

// NewClient creates a new client
func NewClient(username string, conn *websocket.Conn, hub *WebSocketHub) *Client {
	return &Client{clientID: uuid.NewString(), username: username, conn: conn, send: make(chan Event, maxEventBuffer), hub: hub}
}

// Client goroutine to read messages from client
func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxEventSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		logger.Info("received a pong from " + c.username)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var event Event
		err := c.conn.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("error: %v", err)
			}
			logger.Errorf("error: %v", err)
			break
		}
		logger.Info("received an event from client " + c.username)
		c.hub.handleEvent(event, c)
	}
}

// Client goroutine to write messages to client
func (c *Client) write() {
	// send the first ping
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return
	}
	// set up write
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			logger.Info("sending an event to client " + c.username)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.conn.WriteJSON(message)
				if err != nil {
					logger.Error(err)
					break
				}
			}
		case <-ticker.C:
			logger.Info("sending a ping to " + c.username)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
