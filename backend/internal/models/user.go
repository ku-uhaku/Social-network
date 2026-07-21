package models

import "time"

// User represents the database entity and JSON contract for a user profile
type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"` // Handled as ISO 8601 string format
	IsPublic    int    `json:"is_public"`     // 1 for true, 0 for false (SQLite compliance)
	Password    string `json:"-"`             // Password is never exposed in JSON responses

	// The 3 optional database fields (using pointers to allow for nil/null)
	Avatar   *string `json:"avatar"`
	NickName *string `json:"nick_name"`
	AboutMe  *string `json:"about_me"`

	CreatedAt time.Time `json:"created_at"`
}

type InputLoginPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InputRegisterPayload struct {
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Gender      string  `json:"gender"`
	DateOfBirth string  `json:"date_of_birth"`
	Password    string  `json:"password"`
	Avatar      *string `json:"avatar,omitempty"`
	NickName    *string `json:"nick_name,omitempty"`
	AboutMe     *string `json:"about_me,omitempty"`
}

// UpdateProfilePayload represents the incoming JSON contract for changing user details
type UpdateProfilePayload struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Gender      string  `json:"gender"`
	DateOfBirth string  `json:"date_of_birth"`
	IsPublic    int     `json:"is_public"` // 1 for true, 0 for false
	Avatar      *string `json:"avatar"`
	NickName    *string `json:"nick_name"`
	AboutMe     *string `json:"about_me"`
}

// Follow represents the database record for follower relationships
type Follow struct {
	FollowerID  int64     `json:"follower_id"`
	FollowingID int64     `json:"following_id"`
	Status      string    `json:"status"` // 'pending' or 'accepted'
	CreatedAt   time.Time `json:"created_at"`
}

// FollowActionPayload handles targets for follow/unfollow and request actions
type FollowActionPayload struct {
	TargetUserID int64 `json:"target_user_id"`
}

// UserFollowView represents user metadata displayed in follower/following lists
type UserFollowView struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Avatar    *string `json:"avatar"`
	Status    string  `json:"status,omitempty"` // For pending requests
}

// FollowStats summarizes follow counts
type FollowStats struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}
