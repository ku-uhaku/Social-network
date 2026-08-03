package repository

import (
	"context"

	"kuu/internal/models"
)

// CreatePost creates a new post record
func (r *Repository) CreatePost(ctx context.Context, userID int64, payload models.CreatePostPayload) (*models.Post, error) {
	query := `
		INSERT INTO posts (user_id, group_id, title, content, privacy, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, group_id, title, content, privacy, image_url, created_at
	`
	var p models.Post
	err := r.DB.Database.QueryRowContext(
		ctx, query,
		userID, payload.GroupID, payload.Title, payload.Content, payload.Privacy, payload.ImageURL,
	).Scan(&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Content, &p.Privacy, &p.ImageURL, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPostByID fetches a single post
func (r *Repository) GetPostByID(ctx context.Context, postID int64) (*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.group_id, p.title, p.content, p.privacy, p.image_url, p.created_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = $1
		LIMIT 1
	`
	var p models.Post
	err := r.DB.Database.QueryRowContext(ctx, query, postID).Scan(
		&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Content, &p.Privacy, &p.ImageURL, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetFeedPosts retrieves posts filtered by user visibility (Public, Follower posts, Group posts)
func (r *Repository) GetFeedPosts(ctx context.Context, currentUserID int64) ([]models.Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, p.group_id, p.title, p.content, p.privacy, p.image_url, p.created_at,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN follows f ON f.following_id = p.user_id AND f.follower_id = $1 AND f.status = 'accepted'
		LEFT JOIN group_members gm ON gm.group_id = p.group_id AND gm.user_id = $1 AND gm.status = 'accepted'
		WHERE 
			p.user_id = $1 OR
			(p.group_id IS NOT NULL AND gm.user_id IS NOT NULL) OR
			(p.group_id IS NULL AND (p.privacy = 'public' OR (p.privacy = 'almost private' AND f.follower_id IS NOT NULL)))
		ORDER BY p.created_at DESC
	`

	rows, err := r.DB.Database.QueryContext(ctx, query, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Content, &p.Privacy, &p.ImageURL, &p.CreatedAt,
			&p.CommentsCount,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// CreateComment inserts a new post comment
func (r *Repository) CreateComment(ctx context.Context, userID int64, payload models.CreateCommentPayload) (*models.Comment, error) {
	query := `
		INSERT INTO comments (post_id, user_id, content, image_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, post_id, user_id, content, image_url, created_at
	`
	var c models.Comment
	err := r.DB.Database.QueryRowContext(ctx, query, payload.PostID, userID, payload.Content, payload.ImageURL).Scan(
		&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ImageURL, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetPostComments retrieves comments for a post
func (r *Repository) GetPostComments(ctx context.Context, postID int64, currentUserID int64) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.image_url, c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, postID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ImageURL, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
