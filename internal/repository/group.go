package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

	err = tx.QueryRowContext(ctx, groupQuery, payload.Title, payload.Description, creatorID, payload.IsPublic).Scan(
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

// IsGroupMember checks if a user is an accepted member/creator of the group
func (r *Repository) IsGroupMember(ctx context.Context, userID int64, groupID int64) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 FROM group_members 
            WHERE user_id = $1 AND group_id = $2 AND status = 'accepted'
        )
    `
	var exists bool
	err := r.DB.Database.QueryRowContext(ctx, query, userID, groupID).Scan(&exists)
	return exists, err
}

// AddGroupMembersPending safely bulk inserts invitations with a 'pending' state
func (r *Repository) AddGroupMembersPending(ctx context.Context, groupID int64, userIDs []int64) error {
	tx, err := r.DB.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// SQLite compliant UPSERT/IGNORE string so we don't crash on duplicate invitations
	query := `
        INSERT INTO group_members (user_id, group_id, status)
        VALUES ($1, $2, 'pending')
        ON CONFLICT(user_id, group_id) DO NOTHING
    `

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, userID := range userIDs {
		_, err := stmt.ExecContext(ctx, userID, groupID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) AcceptGroupInvitation(ctx context.Context, userID int64, groupID int64) error {
	// First, check if an invitation or request actually exists for this user
	var currentStatus string
	checkQuery := `SELECT status FROM group_members WHERE user_id = $1 AND group_id = $2`

	err := r.DB.Database.QueryRowContext(ctx, checkQuery, userID, groupID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no invitation or join request found for this group")
		}
		return err
	}

	if currentStatus == "accepted" {
		return errors.New("you are already a member of this group")
	}

	// Update the status to accepted
	updateQuery := `
		UPDATE group_members 
		SET status = 'accepted', joined_at = $1 
		WHERE user_id = $2 AND group_id = $3
	`
	_, err = r.DB.Database.ExecContext(ctx, updateQuery, time.Now(), userID, groupID)
	return err
}
