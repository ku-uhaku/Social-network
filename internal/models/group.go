package models

import "time"

// Group represents the core database entity and JSON contract for a group
type Group struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
	IsPublic    bool      `json:"is_public"` // 1 for true, 0 for false (SQLite compliance)
	CreatedAt   time.Time `json:"created_at"`
}

// GroupMember represents a membership relation status entry
type GroupMember struct {
	UserID   int64     `json:"user_id"`
	GroupID  int64     `json:"group_id"`
	Status   string    `json:"status"` // 'pending' or 'accepted'
	JoinedAt time.Time `json:"joined_at"`
}

type CreateGroupPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *int   `json:"is_public"` // Pointer to distinguish between missing vs explicitly sending 0
}

// InviteMemberPayload represents an invitation structure
type InviteMemberPayload struct {
	TargetUserID int64 `json:"target_user_id"`
}

// HandleRequestPayload represents an approval action payload schema
type HandleRequestPayload struct {
	TargetUserID int64 `json:"target_user_id"`
}
