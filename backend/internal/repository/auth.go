package repository

import (
	"database/sql"
	"errors"
	"strings"

	"kuu/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticationUser finds the profile and verifies the bcrypt password hash
func (r *Repository) AuthenticationUser(payload models.InputLoginPayload) (*models.User, error) {
	var user models.User
	var hashedPassword string
	query := `
		SELECT id, username, email, first_name, last_name, gender, date_of_birth, 
		       is_public, password, avatar, about_me, created_at 
		FROM users 
		WHERE LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1) 
		LIMIT 1
	`

	trimmedIdentifier := strings.TrimSpace(payload.Login)

	err := r.DB.Database.QueryRow(query, trimmedIdentifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.DateOfBirth,
		&user.IsPublic,
		&hashedPassword, // Scan password to check against plaintext input
		&user.Avatar,
		&user.AboutMe,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("credentials match failure")
		}
		return nil, err
	}

	// Verify plaintext matches the encoded hash safely
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
	if err != nil {
		return nil, errors.New("credentials match failure")
	}

	return &user, nil
}
