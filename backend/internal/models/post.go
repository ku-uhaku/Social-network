package models

import "time"

// Post represents a post entry
type Post struct {
	ID            int64        `json:"id"`
	UserID        int64        `json:"user_id"`
	User          UserMetadata `json:"user"`
	GroupID       *int64       `json:"group_id,omitempty"`
	Title         string       `json:"title"`
	Content       string       `json:"content"`
	Privacy       string       `json:"privacy"` // 'public', 'almost private', 'private'
	ImageURL      *string      `json:"image_url,omitempty"`
	CommentsCount int            `json:"comments_count"`
	CreatedAt     time.Time      `json:"created_at"`
	Viewers       []UserMetadata `json:"viewers,omitempty"` // Full viewer metadata for private posts
}

// CreatePostPayload represents data for creating a post
type CreatePostPayload struct {
	GroupID  *int64  `json:"group_id,omitempty"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Privacy  string  `json:"privacy"` // 'public', 'almost private', 'private'
	ImageURL *string `json:"image_url,omitempty"`
	VisibleTo []int64 `json:"visible_to,omitempty"` // New field to track viewers for private posts
}

// UserMetadata represents the author metadata displayed alongside a post or comment
type UserMetadata struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Avatar    *string `json:"avatar"`
}

// Comment represents a post comment
type Comment struct {
	ID        int64        `json:"id"`
	PostID    int64        `json:"post_id"`
	UserID    int64        `json:"user_id"`
	User      UserMetadata `json:"user"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	ImageURL  *string      `json:"image_url,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

// CreateCommentPayload represents data for creating a comment
type CreateCommentPayload struct {
	PostID   int64   `json:"post_id"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url,omitempty"`
}