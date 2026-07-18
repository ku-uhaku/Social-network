package repository

import (
	"context" // <-- Ensure this is imported
	"database/sql"
	"errors"
	"strings"

	"kuu/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) AuthficationUser(payload models.InputLoginPayload) (*models.User, error) {
	var user models.User
	var hashedPassword string

	query := `
		SELECT id, username, email, first_name, last_name, gender, date_of_birth, 
		       is_public, password, avatar, nick_name, about_me, created_at 
		FROM users 
		WHERE LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1) 
		LIMIT 1
	`

	trimmedIdentifier := strings.TrimSpace(payload.Identifier)

	// Here we use context.Background() as the baseline runtime scope context
	err := r.DB.Database.QueryRowContext(context.Background(), query, trimmedIdentifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&hashedPassword,
		&user.Avatar,
		&user.NickName,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
