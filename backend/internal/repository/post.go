package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kuu/internal/models"
)

// CreatePost creates a new post in the database
func (r *Repository) CreatePost(ctx context.Context, userID int64, payload models.CreatePostPayload) (*models.Post, error) {
	tx, err := r.DB.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("failed to commit transaction: %w", commitErr)
		}
	}()

	// Insert post into posts table
	post := models.Post{
		UserID:    userID,
		GroupID:   payload.GroupID,
		Title:     payload.Title,
		Content:   payload.Content,
		Privacy:   payload.Privacy,
		ImageURL:  payload.ImageURL,
		CreatedAt: time.Now(),
	}

	stmt, err := tx.Prepare("INSERT INTO posts (user_id, group_id, title, content, privacy, image_url, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(userID, payload.GroupID, payload.Title, payload.Content, payload.Privacy, payload.ImageURL, post.CreatedAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert post: %w", err)
	}
	post.ID = id

	// If private post, insert viewers into post_viewers table
	if payload.Privacy == "private" && len(payload.VisibleTo) > 0 {
		viewerStmt, err := tx.Prepare("INSERT INTO post_viewers (post_id, user_id) VALUES ($1, $2)")
		if err != nil {
			return nil, fmt.Errorf("failed to prepare viewer insert statement: %w", err)
		}
		defer viewerStmt.Close()

		for _, viewerID := range payload.VisibleTo {
			_, err = viewerStmt.Exec(post.ID, viewerID)
			if err != nil {
				return nil, fmt.Errorf("failed to insert viewer: %w", err)
			}
		}
	}

	// Get author metadata
	author, err := r.getUserMetadata(ctx, userID)
	if err != nil {
		return nil, err
	}
	post.User = *author

	return &post, nil
}

// GetPostByID fetches a single post by ID
func (r *Repository) GetPostByID(ctx context.Context, postID int64) (*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.group_id, p.title, p.content, p.privacy, p.image_url, p.created_at,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_count,
			u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = $1
		LIMIT 1
	`
	var p models.Post
	err := r.DB.Database.QueryRowContext(ctx, query, postID).Scan(
		&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Content, &p.Privacy, &p.ImageURL, &p.CreatedAt,
		&p.CommentsCount,
		&p.User.ID, &p.User.Username, &p.User.FirstName, &p.User.LastName, &p.User.Avatar,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetFeedPosts retrieves posts from the current user and users they follow (accepted status),
// plus group posts for groups they belong to. Paginated by cursor (last post id) and limit.
func (r *Repository) GetFeedPosts(ctx context.Context, currentUserID int64, limit int, cursor *int64) ([]models.Post, bool, error) {
	query := `
		SELECT p.id, p.user_id, p.group_id, p.title, p.content, p.privacy, p.image_url, p.created_at,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_count,
			u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE
		p.group_id IS NULL AND 
			(p.user_id = $1   OR
			p.privacy = 'public' OR
			( p.privacy ='almost private' AND EXISTS (
				SELECT 1 FROM follows f
				WHERE f.following_id = p.user_id AND f.follower_id = $1 AND f.status = 'accepted'
			)) OR
			( p.privacy = 'private' AND EXISTS (
				SELECT 1 FROM post_viewers pv
				WHERE pv.post_id = p.id AND pv.user_id = $1
			)))
		ORDER BY p.created_at DESC, p.id DESC
	`
	args := []interface{}{currentUserID}
	if cursor != nil {
		query += ` AND p.id < $2`
		args = append(args, *cursor)
	}
	query += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
	args = append(args, limit+1)

	rows, err := r.DB.Database.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Content, &p.Privacy, &p.ImageURL, &p.CreatedAt,
			&p.CommentsCount,
			&p.User.ID, &p.User.Username, &p.User.FirstName, &p.User.LastName, &p.User.Avatar,
		); err != nil {
			return nil, false, err
		}
		posts = append(posts, p)
	}

	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	fmt.Println("POSTS:::", posts)
	return posts, hasMore, nil
}

// GetGroupFeedPosts retrieves posts belonging to a single group.
// Caller must verify the requesting user is an accepted group member.
// Paginated by cursor (last post id) and limit.
func (r *Repository) GetGroupFeedPosts(ctx context.Context, groupID int64, limit int, cursor *int64) ([]models.Post, bool, error) {
	query := `
		SELECT p.id, p.user_id, p.group_id, p.title, p.content, p.privacy, p.image_url, p.created_at,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_count,
			u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.group_id = $1
	`
	args := []interface{}{groupID}
	if cursor != nil {
		query += ` AND p.id < $2`
		args = append(args, *cursor)
	}
	query += ` ORDER BY p.created_at DESC, p.id DESC`
	query += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
	args = append(args, limit+1)

	rows, err := r.DB.Database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.GroupID, &p.Title, &p.Content, &p.Privacy, &p.ImageURL, &p.CreatedAt,
			&p.CommentsCount,
			&p.User.ID, &p.User.Username, &p.User.FirstName, &p.User.LastName, &p.User.Avatar,
		); err != nil {
			return nil, false, err
		}
		posts = append(posts, p)
	}

	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	return posts, hasMore, nil
}

// CreateComment inserts a new post comment
func (r *Repository) CreateComment(ctx context.Context, userID int64, payload models.CreateCommentPayload) (*models.Comment, error) {
	query := `
		INSERT INTO comments (post_id, user_id, title, content, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, post_id, user_id, title, content, image_url, created_at
	`
	var c models.Comment
	err := r.DB.Database.QueryRowContext(ctx, query, payload.PostID, userID, payload.Title, payload.Content, payload.ImageURL).Scan(
		&c.ID, &c.PostID, &c.UserID, &c.Title, &c.Content, &c.ImageURL, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	author, err := r.getUserMetadata(ctx, c.UserID)
	if err != nil {
		return nil, err
	}
	c.User = *author

	return &c, nil
}

// GetPostComments retrieves comments for a post
func (r *Repository) GetPostComments(ctx context.Context, postID int64) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.title, c.content, c.image_url, c.created_at,
			u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Title, &c.Content, &c.ImageURL, &c.CreatedAt,
			&c.User.ID, &c.User.Username, &c.User.FirstName, &c.User.LastName, &c.User.Avatar,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// getUserMetadata fetches author metadata for a post or comment
func (r *Repository) getUserMetadata(ctx context.Context, userID int64) (*models.UserMetadata, error) {
	query := `
		SELECT id, username, first_name, last_name, avatar
		FROM users
		WHERE id = $1
	`
	var u models.UserMetadata
	err := r.DB.Database.QueryRowContext(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Avatar,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// IsPostViewer checks if a user is allowed to view a post
func (r *Repository) IsPostViewer(ctx context.Context, postID int64, userID int64) (bool, error) {
	row := r.DB.Database.QueryRowContext(ctx, "SELECT 1 FROM post_viewers WHERE post_id = $1 AND user_id = $2", postID, userID)
	var found int
	err := row.Scan(&found)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check viewer status: %w", err)
	}
	return found > 0, nil
}

// GetPostViewers fetches the users a private post was explicitly shared with
func (r *Repository) GetPostViewers(ctx context.Context, postID int64) ([]models.UserMetadata, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM post_viewers pv
		JOIN users u ON u.id = pv.user_id
		WHERE pv.post_id = $1
		ORDER BY u.username ASC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viewers []models.UserMetadata
	for rows.Next() {
		var v models.UserMetadata
		if err := rows.Scan(
			&v.ID, &v.Username, &v.FirstName, &v.LastName, &v.Avatar,
		); err != nil {
			return nil, err
		}
		viewers = append(viewers, v)
	}
	return viewers, nil
}
