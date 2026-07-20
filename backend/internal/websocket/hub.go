package websocket

import (
	"encoding/json"
	"log"
)

type Hub struct {
	Clients     map[*Client]bool
	OnlineUsers map[int64]int // Tracks UserID -> Number of open tabs
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan Event
}

func New() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		OnlineUsers: make(map[int64]int),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan Event),
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
				close(client.Send)
				delete(h.Clients, client)
				client.Conn.Close()
			}
		}
	}
}
