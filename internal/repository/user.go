package repository

import (
	"context"

	"kuu/internal/models"
)

func (r *Repository) GetUserByID(id int64) (*models.User, error) {
	var user models.User

	// Explicitly call out every table column column name you want to fetch
	query := `
		SELECT id, email, username, password, created_at 
		FROM users 
		WHERE id = $1 
		LIMIT 1
	`

	// Match the order of your struct scan references EXACTLY to the fields in the SELECT string above
	err := r.DB.Database.QueryRowContext(context.Background(), query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err // Automatically bubbles up sql.ErrNoRows if user is missing
	}

	return &user, nil
}
