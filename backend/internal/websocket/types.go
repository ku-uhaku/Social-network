package websocket

import "github.com/gorilla/websocket"

// Event defines our flexible real-time message payload distribution rules
type Event struct {
	Type           string      `json:"type"`             // e.g., "new_post", "user_online", "chat_message"
	Payload        interface{} `json:"payload"`          // The actual JSON schema object data
	BroadcastToAll bool        `json:"broadcast_to_all"` // Set true to send to EVERYONE online across all tabs
	TargetUserID   int64       `json:"target_user_id"`   // Sent to a single explicit user ID (all tabs)
	TargetUsersIDs []int64     `json:"target_users_ids"` // Sent to a list of specific user IDs
	TargetGroupID  int64       `json:"target_group_id"`  // Tagged for group context distribution
}

// Client represents a single open connection (like a single browser tab)
type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}
