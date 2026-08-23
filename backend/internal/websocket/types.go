package websocket

import "github.com/gorilla/websocket"

// Event types the server pushes to the browser.
const (
	EventNewDirectMessage    = "new_direct_message"
	EventNewGroupMessage     = "new_group_message"
	EventNewNotification     = "new_notification"
	EventNotificationRead    = "notification_read"
	EventNotificationExpired = "notification_expired"
)

// Event types the browser sends to the server.
const (
	ClientSendDirectMessage = "send_direct_message"
	ClientSendGroupMessage  = "send_group_message"
)

// Event is one real-time message. UserIDs lists who should receive it and is
// kept out of the JSON sent to the browser.
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	UserIDs []int64     `json:"-"`
}

func (e Event) isFor(userID int64) bool {
	for _, id := range e.UserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// Client is a single open connection (one browser tab).
type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}
