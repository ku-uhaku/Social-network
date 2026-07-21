package models

import "time"

// Group represents the core database entity and JSON contract for a group
type Group struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
	IsPublic    int       `json:"is_public"` // 1 for true, 0 for false (SQLite compliance)
	CreatedAt   time.Time `json:"created_at"`
}

type GroupMember struct {
	UserID   int64     `json:"user_id"`
	GroupID  int64     `json:"group_id"`
	Status   string    `json:"status"` // 'pending', 'accepted', 'declined'
	JoinedAt time.Time `json:"joined_at"`
}

// InviteMembersPayload allows the group creator to invite one or multiple users
type InviteMembersPayload struct {
	GroupID       int64   `json:"group_id"`
	TargetUserIDs []int64 `json:"target_user_ids"`
}

// GroupInvitationResponse represents incoming invitation/request items for a user UI
type GroupInvitationView struct {
	GroupID     int64     `json:"group_id"`
	GroupTitle  string    `json:"group_title"`
	InvitedBy   int64     `json:"invited_by"`
	Status      string    `json:"status"`
	RequestedAt time.Time `json:"requested_at"`
}
// GroupActionPayload handles single-user membership actions (accept, decline, join, leave)
type GroupActionPayload struct {
	GroupID      int64 `json:"group_id"`
	TargetUserID int64 `json:"target_user_id,omitempty"` // Used by creator when accepting join requests
}

// CreateGroupPayload defines incoming request payload for group creation
type CreateGroupPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *int   `json:"is_public"` // Pointer allows 0 to be valid JSON input
}

// UpdateGroupPayload defines incoming request payload for updates
type UpdateGroupPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *int   `json:"is_public"`
}
