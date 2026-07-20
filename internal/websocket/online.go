package websocket

import "log"

// handleClientRegister executes actions when a new tab opens
func (h *Hub) handleClientRegister(client *Client) {
	h.Clients[client] = true
	h.OnlineUsers[client.UserID]++

	// First tab open means the user is officially online
	if h.OnlineUsers[client.UserID] == 1 {
		h.BroadcastUserStatus("user_online", client.UserID)
	}
	log.Printf("[WS] User %d opened a tab. Total connections: %d", client.UserID, len(h.Clients))
}

// handleClientUnregister executes actions when a tab closes
func (h *Hub) handleClientUnregister(client *Client) {
	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		close(client.Send)
		client.Conn.Close()

		h.OnlineUsers[client.UserID]--

		// If no tabs remain open, the user is officially offline
		if h.OnlineUsers[client.UserID] <= 0 {
			delete(h.OnlineUsers, client.UserID)
			h.BroadcastUserStatus("user_offline", client.UserID)
		}
		log.Printf("[WS] User %d closed a tab. Remaining connections: %d", client.UserID, len(h.Clients))
	}
}

// shouldSendToClient handles our routing matrix rules
func (h *Hub) shouldSendToClient(client *Client, event Event) bool {
	if event.BroadcastToAll {
		return true
	}
	if event.TargetUserID != 0 && client.UserID == event.TargetUserID {
		return true
	}
	if len(event.TargetUsersIDs) > 0 {
		for _, id := range event.TargetUsersIDs {
			if client.UserID == id {
				return true
			}
		}
	}
	if event.TargetGroupID != 0 {
		// Frontend handles view-filtering natively based on group context markers
		return true
	}
	return false
}

// BroadcastUserStatus safely tells everyone when a user profile shifts visibility state
func (h *Hub) BroadcastUserStatus(eventType string, userID int64) {
	event := Event{
		Type:           eventType,
		BroadcastToAll: true,
		Payload: map[string]interface{}{
			"user_id": userID,
		},
	}
	// Run non-blocking to prevent locking the execution loop
	go func() {
		h.Broadcast <- event
	}()
}

// GetOnlineUsers snapshot utility helper
func (h *Hub) GetOnlineUsers() []int64 {
	users := make([]int64, 0, len(h.OnlineUsers))
	for userID := range h.OnlineUsers {
		users = append(users, userID)
	}
	return users
}
