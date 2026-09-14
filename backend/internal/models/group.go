package models

import "time"

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

type InviteMembersPayload struct {
	GroupID       int64   `json:"group_id"`
	TargetUserIDs []int64 `json:"target_user_ids"`
}

type GroupInvitationView struct {
	GroupID     int64     `json:"group_id"`
	GroupTitle  string    `json:"group_title"`
	InvitedBy   int64     `json:"invited_by"`
	Status      string    `json:"status"`
	RequestedAt time.Time `json:"requested_at"`
}
type GroupActionPayload struct {
	GroupID      int64 `json:"group_id"`
	TargetUserID int64 `json:"target_user_id,omitempty"` // Used by creator when accepting join requests
}

type CreateGroupPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *int   `json:"is_public"` // Pointer allows 0 to be valid JSON input
}

type UpdateGroupPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *int   `json:"is_public"`
}

type GroupEvent struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	Status      string    `json:"status"` // 'upcoming', 'cancelled', 'expired'
	CreatedAt   time.Time `json:"created_at"`
}

type GroupEventWithCounts struct {
	GroupEvent
	GoingCount    int64  `json:"going_count"`
	NotGoingCount int64  `json:"not_going_count"`
	MyStatus      string `json:"my_status"` // 'going', 'not_going', or '' when unset

	// Joined creator info for display
	CreatorUsername *string `json:"creator_username,omitempty"`
	CreatorAvatar   *string `json:"creator_avatar,omitempty"`
}

type CreateGroupEventPayload struct {
	GroupID     int64     `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
}

type EventResponsePayload struct {
	EventID int64  `json:"event_id"`
	Status  string `json:"status"` // 'going' | 'not_going'
}
