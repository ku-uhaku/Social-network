package websocket

import "github.com/gorilla/websocket"

const (
	EventNewDirectMessage    = "new_direct_message"
	EventNewGroupMessage     = "new_group_message"
	EventNewNotification     = "new_notification"
	EventNotificationRead    = "notification_read"
	EventNotificationExpired = "notification_expired"
)

const (
	ClientSendDirectMessage = "send_direct_message"
	ClientSendGroupMessage  = "send_group_message"
)

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

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}
