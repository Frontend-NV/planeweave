package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*Client]struct{})}
}

func (h *Hub) Register(projectID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[projectID] == nil {
		h.rooms[projectID] = make(map[*Client]struct{})
	}
	h.rooms[projectID][c] = struct{}{}
	h.broadcastPresenceLocked(projectID)
}

func (h *Hub) Unregister(projectID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.rooms[projectID]; ok {
		delete(room, c)
		if len(room) == 0 {
			delete(h.rooms, projectID)
		} else {
			h.broadcastPresenceLocked(projectID)
		}
	}
}

func (h *Hub) Broadcast(projectID string, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	room := h.rooms[projectID]
	for client := range room {
		client.Send(event)
	}
}

func (h *Hub) broadcastPresenceLocked(projectID string) {
	count := len(h.rooms[projectID])
	event := Event{
		Type: "presence.updated",
		Payload: map[string]any{
			"projectId":   projectID,
			"onlineCount": count,
		},
	}
	for client := range h.rooms[projectID] {
		client.Send(event)
	}
}

type Client struct {
	conn *websocket.Conn
	send chan Event
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn, send: make(chan Event, 16)}
}

func (c *Client) Send(event Event) {
	select {
	case c.send <- event:
	default:
		log.Println("ws: drop slow client")
	}
}

func (c *Client) WritePump() {
	for event := range c.send {
		if err := c.conn.WriteJSON(event); err != nil {
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

func WriteConnected(conn *websocket.Conn, projectID string) error {
	return conn.WriteJSON(Event{
		Type:    "connected",
		Payload: map[string]string{"projectId": projectID},
	})
}

func DecodeEvent(data []byte) (Event, error) {
	var e Event
	err := json.Unmarshal(data, &e)
	return e, err
}
