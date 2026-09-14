package websocket

import (
	"encoding/json"
)

// Hub holds every open connection and routes events to them.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	disconnect chan string
	events     chan Event
}

func New() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		disconnect: make(chan string),
		events:     make(chan Event, 64),
	}
}

// Run owns the clients map: every change to it happens in this one goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			h.remove(client)

		case sessionID := <-h.disconnect:
			h.disconnectSession(sessionID)

		case event := <-h.events:
			h.route(event)
		}
	}
}

// Register adds a connection, Unregister drops and closes it.
func (h *Hub) Register(client *Client)   { h.register <- client }
func (h *Hub) Unregister(client *Client) { h.unregister <- client }

// DisconnectSession closes every connection bound to the given session token,
// used when that session is destroyed (e.g. logout).
func (h *Hub) DisconnectSession(sessionID string) {
	h.disconnect <- sessionID
}

func (h *Hub) disconnectSession(sessionID string) {
	for client := range h.clients {
		if client.SessionID == sessionID {
			h.remove(client)
		}
	}
}

// BroadcastToUser sends an event to every open tab of one user.
func (h *Hub) BroadcastToUser(userID int64, eventType string, payload interface{}) {
	h.BroadcastToUsers([]int64{userID}, eventType, payload)
}

// BroadcastToUsers sends an event to every open tab of the given users.
func (h *Hub) BroadcastToUsers(userIDs []int64, eventType string, payload interface{}) {
	h.events <- Event{Type: eventType, Payload: payload, UserIDs: userIDs}
}

func (h *Hub) route(event Event) {
	msg, err := json.Marshal(event)
	if err != nil {
		return
	}

	for client := range h.clients {
		if !event.isFor(client.UserID) {
			continue
		}
		select {
		case client.Send <- msg:
		default:
			h.remove(client) // client is not reading, drop it
		}
	}
}

func (h *Hub) remove(client *Client) {
	if !h.clients[client] {
		return
	}
	delete(h.clients, client)
	close(client.Send)
	client.Conn.Close()
}
