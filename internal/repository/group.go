package repository

import (
	"context"

	"kuu/internal/models"
)

func (r *Repository) CreateGroupWithCreator(ctx context.Context, creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	tx, err := r.DB.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Safely roll back changes if execution crashes or errors out

	var group models.Group

	// SQLite conversion compliance mapping: 1 for public, 0 for private
	groupQuery := `
        INSERT INTO groups (title, description, creator_id, is_public) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id, title, description, creator_id, is_public, created_at
    `

	var isPublicInt int
	err = tx.QueryRowContext(ctx, groupQuery, payload.Title, payload.Description, creatorID, payload.IsPublic).Scan(
		&group.ID,
		&group.Title,
		&group.Description,
		&group.CreatorID,
		&isPublicInt,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Automatically enroll creator inside the membership management list as 'accepted'
	memberQuery := `
        INSERT INTO group_members (user_id, group_id, status) 
        VALUES ($1, $2, 'accepted')
    `
	_, err = tx.ExecContext(ctx, memberQuery, creatorID, group.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *Repository) FetchAllPublicGroups(ctx context.Context) ([]models.Group, error) {
	query := `
        SELECT id, title, description, creator_id, is_public, created_at 
        FROM groups 
        WHERE is_public = 1 
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]models.Group, 0) // Prevents returning null in JSON responses
	for rows.Next() {
		var group models.Group
		var isPublicInt int

		err := rows.Scan(
			&group.ID,
			&group.Title,
			&group.Description,
			&group.CreatorID,
			&isPublicInt,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}
