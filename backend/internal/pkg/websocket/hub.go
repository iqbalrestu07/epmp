package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"
)

// Client is one connected websocket consumer.
type Client struct {
	UserID string
	OrgID  string
	conn   *websocket.Conn
	send   chan []byte
}

// Hub manages websocket clients and broadcasts messages to them,
// grouped by organization for multi-tenant isolation.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
	log     zerolog.Logger
}

// NewHub creates a new Hub.
func NewHub(log zerolog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		log:     log,
	}
}

// Register attaches a new websocket connection and starts its writer loop.
// The caller should run the read pump to detect disconnects.
func (h *Hub) Register(userID, orgID string, conn *websocket.Conn) *Client {
	c := &Client{
		UserID: userID,
		OrgID:  orgID,
		conn:   conn,
		send:   make(chan []byte, 32),
	}

	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	go h.writeLoop(c)
	return c
}

// Unregister detaches a client and closes its channel.
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
	h.mu.Unlock()
}

func (h *Hub) writeLoop(c *Client) {
	for msg := range c.send {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.conn.Write(ctx, websocket.MessageText, msg)
		cancel()
		if err != nil {
			h.Unregister(c)
			_ = c.conn.Close(websocket.StatusInternalError, "write failed")
			return
		}
	}
}

// BroadcastToOrg sends a payload to every connected client of an organization.
func (h *Hub) BroadcastToOrg(orgID string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.OrgID != orgID {
			continue
		}
		select {
		case c.send <- payload:
		default:
			// Drop message for slow consumers instead of blocking the hub.
			h.log.Warn().Str("user_id", c.UserID).Msg("websocket send buffer full, dropping notification")
		}
	}
}

// PushToUser sends a payload to every connection owned by a user.
func (h *Hub) PushToUser(userID string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.UserID != userID {
			continue
		}
		select {
		case c.send <- payload:
		default:
			h.log.Warn().Str("user_id", c.UserID).Msg("websocket send buffer full, dropping notification")
		}
	}
}
