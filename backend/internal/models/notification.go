package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Notification types
const (
	NotificationFollowRequest    = "follow_request"
	NotificationGroupInvitation = "group_invitation"
	NotificationGroupJoinRequest = "group_join_request"
	NotificationGroupEvent      = "group_event_created"
)

// JSONText stores arbitrary JSON in a TEXT column, marshaling NULL/empty to null
type JSONText map[string]interface{}

func (j JSONText) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONText) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}
	return json.Unmarshal(data, j)
}

func (j JSONText) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal(map[string]interface{}(j))
}

// Notification represents a single notification in a user's stack
type Notification struct {
	ID          int64     `json:"id"`
	RecipientID int64     `json:"recipient_id"`
	ActorID     *int64    `json:"actor_id,omitempty"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Message    string    `json:"message"`
	Payload    JSONText  `json:"payload"`
	Actions    JSONText  `json:"actions"`
	IsRead     int       `json:"is_read"`
	IsExpired  int       `json:"is_expired"`
	CreatedAt  time.Time `json:"created_at"`

	// Joined actor info for display
	ActorUsername *string `json:"actor_username,omitempty"`
	ActorAvatar  *string `json:"actor_avatar,omitempty"`
}

// NotificationListResponse wraps a page of notifications with the unread count
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount  int64          `json:"unread_count"`
	HasMore      bool           `json:"has_more"`
}

// MarkNotificationReadPayload marks a single notification or all as read
type MarkNotificationReadPayload struct {
	NotificationID *int64 `json:"notification_id,omitempty"`
	All           bool   `json:"all"`
}

// ExpireNotificationPayload expires a single notification
type ExpireNotificationPayload struct {
	NotificationID int64 `json:"notification_id"`
}