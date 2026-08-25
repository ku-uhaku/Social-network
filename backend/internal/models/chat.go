package models

import "time"

// DirectMessage model
type DirectMessage struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupMessage model
type GroupMessage struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	SenderID  int64     `json:"sender_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    *string   `json:"avatar,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Incoming DTO from WebSocket payload
type SendMessagePayload struct {
	ReceiverID *int64 `json:"receiver_id,omitempty"` // For DM
	GroupID    *int64 `json:"group_id,omitempty"`    // For Group
	Content    string `json:"content"`
}

type ConversationMetadata struct {
	UserID        int64      `json:"user_id"`
	Username      string     `json:"username"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Avatar        *string    `json:"avatar"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty"`
	UnreadCount   int        `json:"unread_count"`
}

// DirectHistoryPage is a paginated slice of direct-message history.
type DirectHistoryPage struct {
	Messages []DirectMessage `json:"messages"`
	HasMore  bool            `json:"has_more"`
}
