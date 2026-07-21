package models

import "time"

// Post represents a post entry
type Post struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Username      string    `json:"username"`
	Avatar        *string   `json:"avatar,omitempty"`
	GroupID       *int64    `json:"group_id,omitempty"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Privacy       string    `json:"privacy"` // 'public', 'followers', 'group'
	ImageURL      *string   `json:"image_url,omitempty"`
	LikesCount    int       `json:"likes_count"`
	DislikesCount int       `json:"dislikes_count"`
	CommentsCount int       `json:"comments_count"`
	UserReaction  int       `json:"user_reaction"` // 1: liked, -1: disliked, 0: none
	CreatedAt     time.Time `json:"created_at"`
}

type CreatePostPayload struct {
	GroupID  *int64  `json:"group_id,omitempty"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Privacy  string  `json:"privacy"` // 'public', 'followers', 'group'
	ImageURL *string `json:"image_url,omitempty"`
}

// Comment represents a post comment
type Comment struct {
	ID            int64     `json:"id"`
	PostID        int64     `json:"post_id"`
	UserID        int64     `json:"user_id"`
	Username      string    `json:"username"`
	Avatar        *string   `json:"avatar,omitempty"`
	Content       string    `json:"content"`
	ImageURL      *string   `json:"image_url,omitempty"`
	LikesCount    int       `json:"likes_count"`
	DislikesCount int       `json:"dislikes_count"`
	UserReaction  int       `json:"user_reaction"` // 1: liked, -1: disliked, 0: none
	CreatedAt     time.Time `json:"created_at"`
}

type CreateCommentPayload struct {
	PostID   int64   `json:"post_id"`
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url,omitempty"`
}

// ReactionPayload handles like (1) and dislike (-1) toggles
type ReactionPayload struct {
	PostID    *int64 `json:"post_id,omitempty"`
	CommentID *int64 `json:"comment_id,omitempty"`
	Type      int    `json:"type"` // 1 for Like, -1 for Dislike, 0 to remove
}
