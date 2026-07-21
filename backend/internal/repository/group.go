package repository

import (
	"context"
	"database/sql"
	"strings"

	"kuu/internal/models"
)

// CreateGroup inserts the group AND adds the creator as an accepted member atomically using a transaction
func (r *Repository) CreateGroup(ctx context.Context, creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	tx, err := r.DB.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Safe to call; no-op if Tx is already committed

	// 1. Insert Group
	var group models.Group
	groupQuery := `
		INSERT INTO groups (title, description, creator_id, is_public)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, creator_id, is_public, created_at
	`
	err = tx.QueryRowContext(
		ctx, groupQuery,
		strings.TrimSpace(payload.Title),
		strings.TrimSpace(payload.Description),
		creatorID,
		*payload.IsPublic,
	).Scan(
		&group.ID,
		&group.Title,
		&group.Description,
		&group.CreatorID,
		&group.IsPublic,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 2. Automatically add creator as an accepted member
	memberQuery := `
		INSERT INTO group_members (group_id, user_id, status)
		VALUES ($1, $2, 'accepted')
	`
	_, err = tx.ExecContext(ctx, memberQuery, group.ID, creatorID)
	if err != nil {
		return nil, err
	}

	// 3. Commit both actions
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &group, nil
}

// GetGroupByID fetches a single group by ID
func (r *Repository) GetGroupByID(ctx context.Context, groupID int64) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, title, description, creator_id, is_public, created_at
		FROM groups
		WHERE id = $1
		LIMIT 1
	`
	err := r.DB.Database.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID,
		&group.Title,
		&group.Description,
		&group.CreatorID,
		&group.IsPublic,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if not found
	}

	return &group, nil
}

// GetAllGroups fetches all groups (can be extended later for pagination/filtering)
func (r *Repository) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	query := `
		SELECT id, title, description, creator_id, is_public, created_at
		FROM groups
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []models.Group{}
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.CreatorID, &g.IsPublic, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

// UpdateGroup modifies group details if the caller is the creator
func (r *Repository) UpdateGroup(ctx context.Context, groupID int64, payload models.UpdateGroupPayload) (*models.Group, error) {
	var group models.Group
	query := `
		UPDATE groups
		SET title = $1, description = $2, is_public = $3
		WHERE id = $4
		RETURNING id, title, description, creator_id, is_public, created_at
	`
	err := r.DB.Database.QueryRowContext(
		ctx, query,
		strings.TrimSpace(payload.Title),
		strings.TrimSpace(payload.Description),
		*payload.IsPublic,
		groupID,
	).Scan(
		&group.ID,
		&group.Title,
		&group.Description,
		&group.CreatorID,
		&group.IsPublic,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// DeleteGroup removes a group and its cascading relations (members, posts, invites)
func (r *Repository) DeleteGroup(ctx context.Context, groupID int64) error {
	query := `DELETE FROM groups WHERE id = $1`
	res, err := r.DB.Database.ExecContext(ctx, query, groupID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
