package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Event defines our flexible real-time message payload distribution rules
type Event struct {
	Type    string      `json:"type"`    // e.g., "new_post", "group_invite", "global_alert"
	Payload interface{} `json:"payload"` // The actual JSON schema object data

	BroadcastToAll bool    `json:"broadcast_to_all"` // Set true to send to EVERYONE online across all tabs
	TargetUserID   int64   `json:"target_user_id"`   // Sent to a single explicit user ID (all tabs)
	TargetUsersIDs []int64 `json:"target_users_ids"` // Sent to a list of specific user IDs
	TargetGroupID  int64   `json:"target_group_id"`  // Tagged for group context distribution
}

type Client struct {
	UserID int64 // Track user ID so multiple browser tabs share the same identity marker
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Event
}

func New() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Event),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				client.Conn.Close()
			}

		case event := <-h.Broadcast:
			msgBytes, err := json.Marshal(event)
			if err != nil {
				log.Printf("[WS] Failed to compile JSON framing: %v", err)
				continue
			}

			// Distribute message downstream using your structural routing keys
			for client := range h.Clients {
				shouldSend := false

				// Pattern 1: Global Broadcast to Everyone
				if event.BroadcastToAll {
					shouldSend = true
				}

				// Pattern 2: Single Target User ID (Matches all open tabs for this user)
				if !shouldSend && event.TargetUserID != 0 && client.UserID == event.TargetUserID {
					shouldSend = true
				}

				// Pattern 3: List of Targeted User IDs
				if !shouldSend && len(event.TargetUsersIDs) > 0 {
					for _, id := range event.TargetUsersIDs {
						if client.UserID == id {
							shouldSend = true
							break
						}
					}
				}

				// Pattern 4: Group ID Context
				// Note: The frontend can filter this directly or read it if looking at a specific group feed view.
				if !shouldSend && event.TargetGroupID != 0 {
					// Option A: If your frontend handles the check dynamically, pass it down to everyone:
					shouldSend = true

					// Option B: If you prefer to fetch group membership out of the DB here to filter,
					// you can pass a repository instance to the hub and run a check instead.
				}

				// If it passed any of our routing keys, push it down the channel pipe
				if shouldSend {
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
	}
}
