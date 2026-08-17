package repository

import (
	"context"
	"database/sql"
	"strings"

	"kuu/internal/models"
)

// GetUserByID fetches a complete user profile by its primary key integer
func (r *Repository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User

	query := `
        SELECT id, username, email, first_name, last_name, gender, date_of_birth, is_public, avatar, about_me, created_at
        FROM users 
        WHERE id = $1 
        LIMIT 1
    `

	err := r.DB.Database.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&user.Avatar,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if the ID doesn't exist
	}

	return &user, nil
}

// GetUserByUsername fetches a complete user profile by its unique username
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := `
        SELECT id, username, email, first_name, last_name, gender, date_of_birth, is_public, avatar, about_me, created_at
        FROM users 
        WHERE username = $1 
        LIMIT 1
    `

	err := r.DB.Database.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&user.Avatar,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if the username doesn't exist
	}

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, payload models.InputRegisterPayload, hashedPassword string) (*models.User, error) {
	var user models.User

	query := `
		INSERT INTO users (
			username, email, first_name, last_name, gender, 
			date_of_birth, password, is_public, avatar, about_me
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 1, $8, $9)
		RETURNING id, username, email, first_name, last_name, gender, date_of_birth, is_public, avatar, about_me, created_at
	`

	err := r.DB.Database.QueryRowContext(
		ctx, query,
		strings.TrimSpace(payload.Username),
		strings.ToLower(strings.TrimSpace(payload.Email)),
		payload.FirstName,
		payload.LastName,
		payload.Gender,
		payload.DateOfBirth,
		hashedPassword,
		payload.Avatar,
		payload.AboutMe,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&user.Avatar,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserProfile updates the public/private state of a user profile
func (r *Repository) UpdateUserProfile(ctx context.Context, userID int64, payload models.UpdateProfilePayload) (*models.User, error) {
	var user models.User

	query := `
		UPDATE users 
		SET is_public = $1
		WHERE id = $2
		RETURNING id, username, email, first_name, last_name, gender, date_of_birth, 
		          is_public, avatar, about_me, created_at
	`

	err := r.DB.Database.QueryRowContext(
		ctx, query,
		payload.IsPublic,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&user.Avatar,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserPosts retrieves a user's posts (excluding group posts) with author metadata
// viewerID is the ID of the user viewing the posts (for visibility checks)
func (r *Repository) GetUserPosts(ctx context.Context, targetUserID int64, viewerID int64) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.group_id, p.title, p.content, p.privacy, p.image_url, p.created_at,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comments_count,
			u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = $1 AND p.group_id IS NULL
		AND (
			$2 = $1 OR
			p.privacy = 'public' OR
			(p.privacy = 'almost private' AND EXISTS (
				SELECT 1 FROM follows f
				WHERE f.following_id = p.user_id AND f.follower_id = $2 AND f.status = 'accepted'
			)) OR
			(p.privacy = 'private' AND EXISTS (
				SELECT 1 FROM post_viewers pv
				WHERE pv.post_id = p.id AND pv.user_id = $2
			))
		)
		ORDER BY p.created_at DESC, p.id DESC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, targetUserID, viewerID)
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
			&p.User.ID, &p.User.Username, &p.User.FirstName, &p.User.LastName, &p.User.Avatar,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// GetFollowStats returns accepted follower and following counts for a user
func (r *Repository) GetFollowStats(ctx context.Context, userID int64) (*models.FollowStats, error) {
	var stats models.FollowStats

	query := `
		SELECT
			(SELECT COUNT(*) FROM follows WHERE following_id = $1 AND status = 'accepted') AS followers_count,
			(SELECT COUNT(*) FROM follows WHERE follower_id = $1 AND status = 'accepted') AS following_count
	`
	err := r.DB.Database.QueryRowContext(ctx, query, userID).Scan(&stats.FollowersCount, &stats.FollowingCount)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// InsertFollowRelation inserts or updates a follow relation with appropriate status
func (r *Repository) InsertFollowRelation(ctx context.Context, followerID, targetID int64, status string) error {
	query := `
		INSERT INTO follows (follower_id, following_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (follower_id, following_id) 
		DO UPDATE SET status = EXCLUDED.status
	`
	_, err := r.DB.Database.ExecContext(ctx, query, followerID, targetID, status)
	return err
}

// RemoveFollowRelation removes the follower relationship
func (r *Repository) RemoveFollowRelation(ctx context.Context, followerID, targetID int64) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	_, err := r.DB.Database.ExecContext(ctx, query, followerID, targetID)
	return err
}

// AcceptAllPendingFollows accepts all pending follow requests for a target user,
// returning the list of follower IDs that were accepted.
func (r *Repository) AcceptAllPendingFollows(ctx context.Context, targetUserID int64) ([]int64, error) {
	query := `
		UPDATE follows SET status = 'accepted'
		WHERE following_id = $1 AND status = 'pending'
		RETURNING follower_id
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followerIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, id)
	}
	return followerIDs, rows.Err()
}

// UpdateFollowStatus changes request status ('pending' -> 'accepted')
func (r *Repository) UpdateFollowStatus(ctx context.Context, followerID, targetID int64, status string) error {
	query := `
		UPDATE follows 
		SET status = $1 
		WHERE follower_id = $2 AND following_id = $3
	`
	res, err := r.DB.Database.ExecContext(ctx, query, status, followerID, targetID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetFollowRelation fetches current relation status if any
func (r *Repository) GetFollowRelation(ctx context.Context, followerID, targetID int64) (string, error) {
	var status string
	query := `SELECT status FROM follows WHERE follower_id = $1 AND following_id = $2 LIMIT 1`
	err := r.DB.Database.QueryRowContext(ctx, query, followerID, targetID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

// CanChat reports whether two users share an accepted follow relation in
// either direction (a chat eligibility check, not a login state).
func (r *Repository) CanChat(ctx context.Context, userA, userB int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM follows
			WHERE status = 'accepted'
			  AND ((follower_id = $1 AND following_id = $2)
				OR (follower_id = $2 AND following_id = $1))
		)
	`
	var connected bool
	err := r.DB.Database.QueryRowContext(ctx, query, userA, userB).Scan(&connected)
	if err != nil {
		return false, err
	}
	return connected, nil
}

// GetFollowers retrieves all users following a target user (accepted status only)
func (r *Repository) GetFollowers(ctx context.Context, targetUserID int64) ([]models.UserFollowView, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM follows f
		JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = $1 AND f.status = 'accepted'
		ORDER BY f.created_at DESC
	`
	return r.scanUserFollowViews(ctx, query, targetUserID)
}

// GetFollowing retrieves all users followed by target user (accepted status only)
func (r *Repository) GetFollowing(ctx context.Context, userID int64) ([]models.UserFollowView, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM follows f
		JOIN users u ON u.id = f.following_id
		WHERE f.follower_id = $1 AND f.status = 'accepted'
		ORDER BY f.created_at DESC
	`
	return r.scanUserFollowViews(ctx, query, userID)
}

// GetPendingFollowRequests retrieves follow requests waiting for target user's approval
func (r *Repository) GetPendingFollowRequests(ctx context.Context, targetUserID int64) ([]models.UserFollowView, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM follows f
		JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = $1 AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`
	return r.scanUserFollowViews(ctx, query, targetUserID)
}

// GetAllUsers retrieves lightweight metadata for every user (used by group invites)
func (r *Repository) GetAllUsers(ctx context.Context) ([]models.UserFollowView, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM users u
		ORDER BY u.username ASC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.UserFollowView{}
	for rows.Next() {
		var u models.UserFollowView
		if err := rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Avatar); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetSuggestedUsers returns users the viewer does not follow (and has no
// pending request to), ranked by accepted follower count descending.
func (r *Repository) GetSuggestedUsers(ctx context.Context, viewerID int64, limit int) ([]models.UserFollowView, error) {
	query := `
	SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM users u
		LEFT JOIN follows f ON f.following_id = u.id AND f.follower_id = $1
		WHERE u.id != $1 AND f.follower_id IS NULL
		ORDER BY (
			SELECT COUNT(*) FROM follows ff WHERE ff.following_id = u.id AND ff.status = 'accepted'
		) DESC, u.id ASC
		LIMIT $2
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, viewerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.UserFollowView{}
	for rows.Next() {
		var user models.UserFollowView
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Avatar); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	// println("check if we reciave users",users)
	return users, rows.Err()
}

func (r *Repository) scanUserFollowViews(ctx context.Context, query string, arg int64) ([]models.UserFollowView, error) {
	rows, err := r.DB.Database.QueryContext(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	views := []models.UserFollowView{}
	for rows.Next() {
		var u models.UserFollowView
		if err := rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Avatar); err != nil {
			return nil, err
		}
		views = append(views, u)
	}
	return views, nil
}