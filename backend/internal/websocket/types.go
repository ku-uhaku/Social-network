package websocket

import "github.com/gorilla/websocket"

// Event type constants — call sites use these instead of magic strings
const (
	EventNewDirectMessage   = "new_direct_message"
	EventNewGroupMessage    = "new_group_message"
	EventUserOnline         = "user_online"
	EventUserOffline        = "user_offline"
	EventOnlineUsersList    = "online_users_list"
	EventNewNotification    = "new_notification"
	EventNotificationExpired = "notification_expired"
	EventNotificationRead   = "notification_read"
)

// Inbound event types sent from client to server
const (
	ClientSendDirectMessage = "send_direct_message"
	ClientSendGroupMessage  = "send_group_message"
)

// Event defines our flexible real-time message payload distribution rules
type Event struct {
	Type           string      `json:"type"`             // e.g. "new_notification", "user_online"
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