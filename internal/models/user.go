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
