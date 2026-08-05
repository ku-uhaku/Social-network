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
	DateOfBirth string `json:"date_of_birth"`
	IsPublic    int    `json:"is_public"`
	Password    string `json:"-"`

	Avatar  *string `json:"avatar"`
	AboutMe *string `json:"about_me"`

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
	AboutMe     *string `json:"about_me,omitempty"`
}

// UpdateProfilePayload represents the incoming JSON contract for changing user details
type UpdateProfilePayload struct {
	IsPublic int `json:"is_public"` // 1 for true, 0 for false
}

// UserProfileView represents a user profile as displayed on the profile page,
// including follow statistics and the viewer's relationship to the user.
type UserProfileView struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Gender      string  `json:"gender"`
	DateOfBirth string  `json:"date_of_birth"`
	IsPublic    int     `json:"is_public"`
	Avatar      *string `json:"avatar"`
	AboutMe     *string `json:"about_me"`
	CreatedAt   time.Time `json:"created_at"`

	FollowersCount int64  `json:"followers_count"`
	FollowingCount int64  `json:"following_count"`
	FollowStatus   string `json:"follow_status"` // 'self', 'none', 'pending', 'accepted'
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
}

// FollowStats summarizes follow counts
type FollowStats struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}