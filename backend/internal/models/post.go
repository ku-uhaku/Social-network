package models

import "time"

// Post represents a post entry
type Post struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	GroupID       *int64    `json:"group_id,omitempty"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Privacy       string    `json:"privacy"` // 'public', 'almost private', 'private'
	ImageURL      *string   `json:"image_url,omitempty"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreatePostPayload struct {
	GroupID  *int64  `json:"group_id,omitempty"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Privacy  string  `json:"privacy"` // 'public', 'almost private', 'private'
	ImageURL *string `json:"image_url,omitempty"`
}

// Comment represents a post comment
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCommentPayload struct {
	PostID   int64   `json:"post_id"`
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url,omitempty"`
}
