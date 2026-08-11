package websocket

import (
	"encoding/json"
	"log"
)

type Hub struct {
	Clients               map[*Client]bool
	OnlineUsers           map[int64]int // Tracks UserID -> Number of open tabs
	Register              chan *Client
	Unregister            chan *Client
	Broadcast             chan Event
	onlineSnapshotRequest chan chan []int64
}

func New() *Hub {
	return &Hub{
		Clients:               make(map[*Client]bool),
		OnlineUsers:           make(map[int64]int),
		Register:              make(chan *Client),
		Unregister:            make(chan *Client),
		Broadcast:             make(chan Event),
		onlineSnapshotRequest: make(chan chan []int64),
	}
}

// Run manages channel coordination safely across concurrent routines
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.handleClientRegister(client)

		case client := <-h.Unregister:
			h.handleClientUnregister(client)

		case event := <-h.Broadcast:
			h.routeEvent(event)

		case respCh := <-h.onlineSnapshotRequest:
			users := make([]int64, 0, len(h.OnlineUsers))
			for userID := range h.OnlineUsers {
				users = append(users, userID)
			}
			respCh <- users
		}
	}
}

// routeEvent marshals and sends an event to clients matching the target parameters
func (h *Hub) routeEvent(event Event) {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("[WS] Failed to marshal event: %v", err)
		return
	}

	for client := range h.Clients {
		if h.shouldSendToClient(client, event) {
			select {
			case client.Send <- msgBytes:
			default:
				// Slow client defer cleanup to handleClientUnregister serially
				h.Unregister <- client
			}
		}
	}
}

// BroadcastToUser pushes an event to all tabs of a single user, non-blocking
func (h *Hub) BroadcastToUser(userID int64, eventType string, payload interface{}) {
	h.broadcast(Event{
		Type:         eventType,
		TargetUserID: userID,
		Payload:      payload,
	})
}

// BroadcastToUsers pushes an event to all tabs of multiple users, non-blocking
func (h *Hub) BroadcastToUsers(userIDs []int64, eventType string, payload interface{}) {
	h.broadcast(Event{
		Type:           eventType,
		TargetUsersIDs: userIDs,
		Payload:        payload,
	})
}

func (h *Hub) broadcast(event Event) {
	go func() {
		h.Broadcast <- event
	}()
}