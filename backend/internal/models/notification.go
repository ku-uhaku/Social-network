package models

import "time"

const (
	NotificationFollowRequest    = "follow_request"
	NotificationNewFollower      = "new_follower"
	NotificationGroupInvitation  = "group_invitation"
	NotificationGroupJoinRequest = "group_join_request"
	NotificationGroupEvent       = "group_event_created"
)

type NotificationPayload struct {
	GroupID *int64 `json:"group_id,omitempty"`
}

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

type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int64          `json:"unread_count"`
	HasMore       bool           `json:"has_more"`
}

type MarkNotificationReadPayload struct {
	NotificationID *int64 `json:"notification_id,omitempty"`
	All            bool   `json:"all"`
}

type ExpireNotificationPayload struct {
	NotificationID int64 `json:"notification_id"`
}
