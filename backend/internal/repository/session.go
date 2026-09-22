package repository

import (
	"database/sql"
	"errors"
	"time"

	"kuu/internal/models"
)

// GetSession retrieves a user session using clean pointer semantics
func (r *Repository) GetSession(token string) (*models.Session, error) {
	var session models.Session

	query := `
        SELECT id, user_id, expires_at, created_at 
        FROM sessions 
        WHERE id = $1 
        LIMIT 1
    `

	err := r.DB.Database.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &session, nil
}

// CreateSession inserts a new active session into SQLite
func (r *Repository) CreateSession(userID int64, token string, duration time.Duration) (*models.Session, error) {
	var session models.Session
	expiresAt := time.Now().Add(duration)

	query := `
        INSERT INTO sessions (id, user_id, expires_at) 
        VALUES ($1, $2, $3) 
        RETURNING id, user_id, expires_at, created_at
    `

	err := r.DB.Database.QueryRow(query, token, userID, expiresAt).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetUserBySessionToken joins sessions and users to authenticate active tokens
func (r *Repository) GetUserBySessionToken(token string) (*models.User, error) {
	var user models.User

	query := `
        SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.gender, u.date_of_birth, 
               u.is_public, u.avatar, u.about_me, u.created_at
        FROM sessions s
        JOIN users u ON s.user_id = u.id
        WHERE s.id = $1 AND s.expires_at > CURRENT_TIMESTAMP
        LIMIT 1
    `

	err := r.DB.Database.QueryRow(query, token).Scan(
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

// DeleteSession physically removes the session row on logout
func (r *Repository) DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.DB.Database.Exec(query, token)
	return err
}
