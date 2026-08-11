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

func (r *Repository) InviteUsersBatch(ctx context.Context, groupID int64, userIDs []int64) error {
	tx, err := r.DB.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO group_members (group_id, user_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (group_id, user_id) DO NOTHING
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, uid := range userIDs {
		if _, err := stmt.ExecContext(ctx, groupID, uid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateMemberStatus updates membership status ('accepted', 'declined', etc.)
func (r *Repository) UpdateMemberStatus(ctx context.Context, groupID int64, userID int64, status string) error {
	query := `
		UPDATE group_members 
		SET status = $1 
		WHERE group_id = $2 AND user_id = $3
	`
	res, err := r.DB.Database.ExecContext(ctx, query, status, groupID, userID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// AddMember directly inserts an accepted member (for public group join)
func (r *Repository) AddMember(ctx context.Context, groupID int64, userID int64, status string) error {
	query := `
		INSERT INTO group_members (group_id, user_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (group_id, user_id) DO UPDATE SET status = EXCLUDED.status
	`
	_, err := r.DB.Database.ExecContext(ctx, query, groupID, userID, status)
	return err
}

// RemoveMember deletes a user membership record (Leave group or Kick user)
func (r *Repository) RemoveMember(ctx context.Context, groupID int64, userID int64) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.DB.Database.ExecContext(ctx, query, groupID, userID)
	return err
}

// IsMember returns current member status if row exists
func (r *Repository) GetMemberStatus(ctx context.Context, groupID int64, userID int64) (string, error) {
	var status string
	query := `SELECT status FROM group_members WHERE group_id = $1 AND user_id = $2 LIMIT 1`
	err := r.DB.Database.QueryRowContext(ctx, query, groupID, userID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

// GetGroupMembers returns all accepted members of a group with lightweight metadata
func (r *Repository) GetGroupMembers(ctx context.Context, groupID int64) ([]models.UserFollowView, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM group_members gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.group_id = $1 AND gm.status = 'accepted'
		ORDER BY u.username ASC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []models.UserFollowView{}
	for rows.Next() {
		var u models.UserFollowView
		if err := rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Avatar); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	return members, nil
}

// GetUserPendingInvitations returns all group invitations sent to a user
func (r *Repository) GetUserPendingInvitations(ctx context.Context, userID int64) ([]models.GroupInvitationView, error) {
	query := `
		SELECT g.id, g.title, g.creator_id, gm.status, gm.joined_at
		FROM group_members gm
		JOIN groups g ON g.id = gm.group_id
		WHERE gm.user_id = $1 AND gm.status = 'pending'
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []models.GroupInvitationView
	for rows.Next() {
		var inv models.GroupInvitationView
		if err := rows.Scan(&inv.GroupID, &inv.GroupTitle, &inv.InvitedBy, &inv.Status, &inv.RequestedAt); err != nil {
			return nil, err
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

// CreateGroupEvent inserts a new group event
func (r *Repository) CreateGroupEvent(ctx context.Context, groupID, creatorID int64, payload models.CreateGroupEventPayload) (*models.GroupEvent, error) {
	var e models.GroupEvent
	query := `
		INSERT INTO group_events (group_id, creator_id, title, description, event_time, status)
		VALUES ($1, $2, $3, $4, $5, 'upcoming')
		RETURNING id, group_id, creator_id, title, description, event_time, status, created_at
	`
	err := r.DB.Database.QueryRowContext(ctx, query,
		groupID, creatorID,
		strings.TrimSpace(payload.Title),
		strings.TrimSpace(payload.Description),
		payload.EventTime.UTC(),
	).Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.Status, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetEventByID fetches a single group event
func (r *Repository) GetEventByID(ctx context.Context, eventID int64) (*models.GroupEvent, error) {
	var e models.GroupEvent
	query := `
		SELECT id, group_id, creator_id, title, description, event_time, status, created_at
		FROM group_events
		WHERE id = $1
		LIMIT 1
	`
	err := r.DB.Database.QueryRowContext(ctx, query, eventID).Scan(
		&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.Status, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetGroupEvents marks past upcoming events as expired, then lists all events for a
// group with response tallies and the requesting user's choice.
func (r *Repository) GetGroupEvents(ctx context.Context, groupID, userID int64) ([]models.GroupEventWithCounts, error) {
	if _, err := r.DB.Database.ExecContext(ctx, `
		UPDATE group_events SET status = 'expired'
		WHERE group_id = $1 AND status = 'upcoming' AND event_time <= datetime('now')
	`, groupID); err != nil {
		return nil, err
	}

	query := `
		SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_time, e.status, e.created_at,
			(SELECT COUNT(*) FROM event_responses r WHERE r.event_id = e.id AND r.status = 'going'),
			(SELECT COUNT(*) FROM event_responses r WHERE r.event_id = e.id AND r.status = 'not_going'),
			COALESCE((SELECT myr.status FROM event_responses myr WHERE myr.event_id = e.id AND myr.user_id = $1), ''),
			COALESCE(cu.username, ''), cu.avatar
		FROM group_events e
		LEFT JOIN users cu ON cu.id = e.creator_id
		WHERE e.group_id = $2
		ORDER BY e.event_time ASC, e.id ASC
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, userID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.GroupEventWithCounts
	for rows.Next() {
		var ev models.GroupEventWithCounts
		if err := rows.Scan(
			&ev.ID, &ev.GroupID, &ev.CreatorID, &ev.Title, &ev.Description, &ev.EventTime, &ev.Status, &ev.CreatedAt,
			&ev.GoingCount, &ev.NotGoingCount, &ev.MyStatus,
			&ev.CreatorUsername, &ev.CreatorAvatar,
		); err != nil {
			return nil, err
		}
		if ev.CreatorUsername != nil && *ev.CreatorUsername == "" {
			ev.CreatorUsername = nil
		}
		events = append(events, ev)
	}
	return events, nil
}

// CancelGroupEvent marks an upcoming event as cancelled
func (r *Repository) CancelGroupEvent(ctx context.Context, eventID int64) error {
	res, err := r.DB.Database.ExecContext(ctx, `
		UPDATE group_events SET status = 'cancelled'
		WHERE id = $1 AND status = 'upcoming'
	`, eventID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// SetEventResponse upserts a user's going/not_going choice for an event
func (r *Repository) SetEventResponse(ctx context.Context, eventID, userID int64, status string) error {
	query := `
		INSERT INTO event_responses (event_id, user_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, user_id) DO UPDATE SET status = EXCLUDED.status
	`
	_, err := r.DB.Database.ExecContext(ctx, query, eventID, userID, status)
	return err
}
