package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"kuu/internal/models"
)

// GetUserByID fetches a complete user profile by its primary key integer
func (r *Repository) GetUserByID(userID int64) (*models.User, error) {
	var user models.User

	// 1. Select all properties matching your struct architecture
	query := `
		SELECT id, username, email, first_name, last_name, gender, 
		       date_of_birth, is_public, avatar, nick_name, about_me, created_at
		FROM users 
		WHERE id = $1 
		LIMIT 1
	`

	// 2. Scan row values using memory pointers directly into the target object properties
	err := r.DB.Database.QueryRowContext(context.Background(), query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&user.Avatar,
		&user.NickName,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, payload models.InputRegisterPayload, hashedPassword string) (*models.User, error) {
	var user models.User

	query := `
		INSERT INTO users (
			username, email, first_name, last_name, gender, 
			date_of_birth, password, is_public, avatar, nick_name, about_me
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 1, $8, $9, $10)
		RETURNING id, username, email, first_name, last_name, gender, date_of_birth, is_public, avatar, nick_name, about_me, created_at
	`

	err := r.DB.Database.QueryRowContext(ctx, query,
		strings.TrimSpace(payload.Username),
		strings.ToLower(strings.TrimSpace(payload.Email)),
		payload.FirstName,
		payload.LastName,
		payload.Gender,
		payload.DateOfBirth,
		hashedPassword,
		payload.Avatar,
		payload.NickName,
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
		&user.NickName,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
