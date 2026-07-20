package repository

import (
	"context"
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
			date_of_birth, password, is_public, avatar, nick_name, about_me
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 1, $8, $9, $10)
		RETURNING id, username, email, first_name, last_name, gender, date_of_birth, is_public, avatar, nick_name, about_me, created_at
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

func (r *Repository) UpdateUserProfile(ctx context.Context, userID int64, payload models.UpdateProfilePayload) (*models.User, error) {
	var user models.User

	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, gender = $3, date_of_birth = $4, 
		    is_public = $5, avatar = $6, nick_name = $7, about_me = $8
		WHERE id = $9
		RETURNING id, username, email, first_name, last_name, gender, date_of_birth, 
		          is_public, avatar, nick_name, about_me, created_at
	`

	err := r.DB.Database.QueryRowContext(
		ctx, query,
		payload.FirstName,
		payload.LastName,
		payload.Gender,
		payload.DateOfBirth,
		payload.IsPublic,
		payload.Avatar,
		payload.NickName,
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
		&user.NickName,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
