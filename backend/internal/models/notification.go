package models

import "time"

// Notification types
const (
	NotificationFollowRequest    = "follow_request"
	NotificationGroupInvitation  = "group_invitation"
	NotificationGroupJoinRequest = "group_join_request"
	NotificationGroupEvent       = "group_event_created"
)

// NotificationPayload is the extra data a notification carries. It is stored as
// plain columns rather than JSON, and shown to the browser as one object.
type NotificationPayload struct {
	GroupID *int64 `json:"group_id,omitempty"`
}

// NotificationActions are the buttons a notification offers. They are derived
// from the type rather than stored, so button labels never live in the database.
type NotificationActions struct {
	Buttons []NotificationButton `json:"buttons"`
}

type NotificationButton struct {
	Action string `json:"action"`
	Label  string `json:"label"`
}

var (
	acceptDeclineActions = &NotificationActions{Buttons: []NotificationButton{
		{Action: "accept", Label: "Accept"},
		{Action: "decline", Label: "Decline"},
	}}
	viewEventActions = &NotificationActions{Buttons: []NotificationButton{
		{Action: "view", Label: "View event"},
	}}
)

// NotificationActionsFor returns the buttons a notification type offers, or nil
// when the type is informational.
func NotificationActionsFor(notifType string) *NotificationActions {
	switch notifType {
	case NotificationFollowRequest, NotificationGroupInvitation, NotificationGroupJoinRequest:
		return acceptDeclineActions
	case NotificationGroupEvent:
		return viewEventActions
	default:
		return nil
	}
}

// Notification represents a single notification in a user's stack
type Notification struct {
	ID          int64                `json:"id"`
	RecipientID int64                `json:"recipient_id"`
	ActorID     *int64               `json:"actor_id,omitempty"`
	Type        string               `json:"type"`
	Title       string               `json:"title"`
	Message     string               `json:"message"`
	Payload     NotificationPayload  `json:"payload"`
	Actions     *NotificationActions `json:"actions"`
	IsRead      int                  `json:"is_read"`
	IsExpired   int                  `json:"is_expired"`
	CreatedAt   time.Time            `json:"created_at"`

	// Joined actor info for display
	ActorUsername *string `json:"actor_username,omitempty"`
	ActorAvatar   *string `json:"actor_avatar,omitempty"`
}

// NotificationListResponse wraps a page of notifications with the unread count
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int64          `json:"unread_count"`
	HasMore       bool           `json:"has_more"`
}

// MarkNotificationReadPayload marks a single notification or all as read
type MarkNotificationReadPayload struct {
	NotificationID *int64 `json:"notification_id,omitempty"`
	All            bool   `json:"all"`
}

// ExpireNotificationPayload expires a single notification
type ExpireNotificationPayload struct {
	NotificationID int64 `json:"notification_id"`
}
