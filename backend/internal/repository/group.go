package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"kuu/internal/models"
)

// CreateGroup inserts the group AND adds the creator as an accepted member atomically using a transaction
func (r *Repository) CreateGroup(creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	tx, err := r.DB.Database.Begin()
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
	err = tx.QueryRow(
		groupQuery,
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
	_, err = tx.Exec(memberQuery, group.ID, creatorID)
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
func (r *Repository) GetGroupByID(groupID int64) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, title, description, creator_id, is_public, created_at
		FROM groups
		WHERE id = $1
		LIMIT 1
	`
	err := r.DB.Database.QueryRow(query, groupID).Scan(
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
func (r *Repository) GetAllGroups() ([]models.Group, error) {
	query := `
		SELECT id, title, description, creator_id, is_public, created_at
		FROM groups
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Database.Query(query)
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
func (r *Repository) UpdateGroup(groupID int64, payload models.UpdateGroupPayload) (*models.Group, error) {
	var group models.Group
	query := `
		UPDATE groups
		SET title = $1, description = $2, is_public = $3
		WHERE id = $4
		RETURNING id, title, description, creator_id, is_public, created_at
	`
	err := r.DB.Database.QueryRow(
		query,
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
func (r *Repository) DeleteGroup(groupID int64) error {
	query := `DELETE FROM groups WHERE id = $1`
	res, err := r.DB.Database.Exec(query, groupID)
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

func (r *Repository) InviteUsersBatch(groupID int64, userIDs []int64) error {
	tx, err := r.DB.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO group_members (group_id, user_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (group_id, user_id) DO NOTHING
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, uid := range userIDs {
		if _, err := stmt.Exec(groupID, uid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateMemberStatus updates membership status ('accepted', 'declined', etc.)
func (r *Repository) UpdateMemberStatus(groupID int64, userID int64, status string) error {
	query := `
		UPDATE group_members 
		SET status = $1 
		WHERE group_id = $2 AND user_id = $3
	`
	res, err := r.DB.Database.Exec(query, status, groupID, userID)
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
func (r *Repository) AddMember(groupID int64, userID int64, status string) error {
	query := `
		INSERT INTO group_members (group_id, user_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (group_id, user_id) DO UPDATE SET status = EXCLUDED.status
	`
	_, err := r.DB.Database.Exec(query, groupID, userID, status)
	return err
}

// RemoveMember deletes a user membership record (Leave group or Kick user)
func (r *Repository) RemoveMember(groupID int64, userID int64) error {
	fmt.Println("userID:",userID,"groupID:",groupID)
	tx, err := r.DB.Database.Begin()
	if err != nil {

		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
    DELETE FROM group_members
    WHERE group_id = $1 AND user_id = $2
`, groupID, userID)
	if err != nil {

		return err
	}

	_, err = tx.Exec(`
    DELETE FROM groups
    WHERE id = $1 AND creator_id = $2
`, groupID, userID)
	if err != nil {
		fmt.Println("error database:::",err)
		return err
	}

	return tx.Commit()
}

// IsMember returns current member status if row exists
func (r *Repository) GetMemberStatus(groupID int64, userID int64) (string, error) {
	var status string
	query := `SELECT status FROM group_members WHERE group_id = $1 AND user_id = $2 LIMIT 1`
	err := r.DB.Database.QueryRow(query, groupID, userID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

// GetGroupMembers returns all accepted members of a group with lightweight metadata
func (r *Repository) GetGroupMembers(groupID int64) ([]models.UserFollowView, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM group_members gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.group_id = $1 AND gm.status = 'accepted'
		ORDER BY u.username ASC
	`
	rows, err := r.DB.Database.Query(query, groupID)
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
func (r *Repository) GetUserPendingInvitations(userID int64) ([]models.GroupInvitationView, error) {
	query := `
		SELECT g.id, g.title, g.creator_id, gm.status, gm.joined_at
		FROM group_members gm
		JOIN groups g ON g.id = gm.group_id
		WHERE gm.user_id = $1 AND gm.status = 'pending'
	`
	rows, err := r.DB.Database.Query(query, userID)
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
func (r *Repository) CreateGroupEvent(groupID, creatorID int64, payload models.CreateGroupEventPayload) (*models.GroupEvent, error) {
	var e models.GroupEvent
	query := `
		INSERT INTO group_events (group_id, creator_id, title, description, event_time, status)
		VALUES ($1, $2, $3, $4, $5, 'upcoming')
		RETURNING id, group_id, creator_id, title, description, event_time, status, created_at
	`
	err := r.DB.Database.QueryRow(
		query,
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
func (r *Repository) GetEventByID(eventID int64) (*models.GroupEvent, error) {
	var e models.GroupEvent
	query := `
		SELECT id, group_id, creator_id, title, description, event_time, status, created_at
		FROM group_events
		WHERE id = $1
		LIMIT 1
	`
	err := r.DB.Database.QueryRow(query, eventID).Scan(
		&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.Status, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetGroupEvents lists all events for a group with response tallies and the
// requesting user's choice. Past upcoming events are reported as 'expired'
// without writing to the database, keeping this a read-only query.
func (r *Repository) GetGroupEvents(groupID, userID int64) ([]models.GroupEventWithCounts, error) {
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
	rows, err := r.DB.Database.Query(query, userID, groupID)
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
		if ev.Status == "upcoming" && ev.EventTime.Before(time.Now()) {
			ev.Status = "expired"
		}
		if ev.CreatorUsername != nil && *ev.CreatorUsername == "" {
			ev.CreatorUsername = nil
		}
		events = append(events, ev)
	}
	return events, nil
}

// GetEventWithCounts fetches a single event with response tallies and the
// caller's choice. Used to return fresh data after recording a response.
func (r *Repository) GetEventWithCounts(eventID, userID int64) (*models.GroupEventWithCounts, error) {
	query := `
		SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_time, e.status, e.created_at,
			(SELECT COUNT(*) FROM event_responses r WHERE r.event_id = e.id AND r.status = 'going'),
			(SELECT COUNT(*) FROM event_responses r WHERE r.event_id = e.id AND r.status = 'not_going'),
			COALESCE((SELECT myr.status FROM event_responses myr WHERE myr.event_id = e.id AND myr.user_id = $1), ''),
			COALESCE(cu.username, ''), cu.avatar
		FROM group_events e
		LEFT JOIN users cu ON cu.id = e.creator_id
		WHERE e.id = $2
		LIMIT 1
	`
	var ev models.GroupEventWithCounts
	if err := r.DB.Database.QueryRow(query, userID, eventID).Scan(
		&ev.ID, &ev.GroupID, &ev.CreatorID, &ev.Title, &ev.Description, &ev.EventTime, &ev.Status, &ev.CreatedAt,
		&ev.GoingCount, &ev.NotGoingCount, &ev.MyStatus,
		&ev.CreatorUsername, &ev.CreatorAvatar,
	); err != nil {
		return nil, err
	}
	if ev.Status == "upcoming" && ev.EventTime.Before(time.Now()) {
		ev.Status = "expired"
	}
	if ev.CreatorUsername != nil && *ev.CreatorUsername == "" {
		ev.CreatorUsername = nil
	}
	return &ev, nil
}

// CancelGroupEvent marks an upcoming event as cancelled
func (r *Repository) CancelGroupEvent(eventID int64) error {
	res, err := r.DB.Database.Exec(`
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
func (r *Repository) SetEventResponse(eventID, userID int64, status string) error {
	query := `
		INSERT INTO event_responses (event_id, user_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, user_id) DO UPDATE SET status = EXCLUDED.status
	`
	_, err := r.DB.Database.Exec(query, eventID, userID, status)
	return err
}
