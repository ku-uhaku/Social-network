package repository

import (
	"context"
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

	err := r.DB.Database.QueryRowContext(context.Background(), query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Clean return of nil for the data payload when nothing matches
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Return the memory address of our populated session struct
	return &session, nil
}

func (r *Repository) CreateSession(ctx context.Context, userID int64, token string, duration time.Duration) (*models.Session, error) {
	var session models.Session
	expiresAt := time.Now().Add(duration)

	query := `
		INSERT INTO sessions (id, user_id, expires_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, expires_at, created_at
	`

	err := r.DB.Database.QueryRowContext(ctx, query, token, userID, expiresAt).Scan(
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
