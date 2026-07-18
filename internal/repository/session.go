package repository

import (
	"context"
	"database/sql"
	"time"

	"kuu/internal/models"
)

// CreateSession deletes all existing sessions for a user, then creates and returns a fresh one.
func (r *Repository) CreateSession(userID int64, sessionID string, duration time.Duration) (*models.Session, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(duration)

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	query := `
		DELETE FROM sessions WHERE user_id = $1;
		INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES ($2, $1, $3, $4);
	`

	_, err := r.DB.Database.ExecContext(context.Background(), query, userID, session.ID, session.ExpiresAt, session.CreatedAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by its token ID to verify if a user is logged in.
func (r *Repository) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = $1 LIMIT 1`

	err := r.DB.Database.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if not found
	}

	return &session, nil
}

// DeleteSession destroys a specific session token (used for logging out).
func (r *Repository) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`

	result, err := r.DB.Database.ExecContext(context.Background(), query, sessionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Let the handler know nothing was deleted
	}

	return nil
}
