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
        SELECT id, username, email, first_name, last_name, gender, date_of_birth, is_public, created_at
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
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if the ID doesn't exist
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

func (r *Repository) UpdateUserProfile(ctx context.Context, userID int64, payload models.UpdateProfilePayload) (*models.User, error) {
	var user models.User

	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, gender = $3, date_of_birth = $4, 
		    is_public = $5, avatar = $6, about_me = $7
		WHERE id = $9
		RETURNING id, username, email, first_name, last_name, gender, date_of_birth, 
		          is_public, avatar, about_me, created_at
	`

	err := r.DB.Database.QueryRowContext(
		ctx, query,
		payload.FirstName,
		payload.LastName,
		payload.Gender,
		payload.DateOfBirth,
		payload.IsPublic,
		payload.Avatar,
		payload.AboutMe,
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
