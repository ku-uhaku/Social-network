package repository

import (
	"database/sql"
	"errors"

	"kuu/internal/models"
)

const notificationColumns = `
	n.id, n.recipient_id, n.actor_id, n.type, n.title, n.message,
	n.payload, n.actions, n.is_read, n.is_expired, n.created_at,
	u.username, u.avatar
`

const notificationFrom = `
	FROM notifications n
	LEFT JOIN users u ON u.id = n.actor_id
`

// CreateNotification persists a new notification and returns it with the actor joined in
func (r *Repository) CreateNotification(n *models.Notification) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (recipient_id, actor_id, type, title, message, payload, actions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	res, err := r.DB.Database.Exec(query,
		n.RecipientID, n.ActorID, n.Type, n.Title, n.Message, n.Payload, n.Actions)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	created, err := r.GetNotification(n.RecipientID, id)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, sql.ErrNoRows
	}
	return created, nil
}

// GetUserNotifications returns a page of notifications for a recipient, newest
// first, using lastID as a cursor (0 = newest page). Fetches limit+1 rows to
// report whether a further page exists.
func (r *Repository) GetUserNotifications(recipientID int64, limit int, lastID int64) ([]models.Notification, bool, error) {
	query := `
		SELECT ` + notificationColumns + notificationFrom + `
		WHERE n.recipient_id = $1 AND ($2 = 0 OR n.id < $2)
		ORDER BY n.id DESC
		LIMIT $3
	`
	rows, err := r.DB.Database.Query(query, recipientID, lastID, limit+1)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	notifications := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(notificationFields(&n)...); err != nil {
			return nil, false, err
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(notifications) > limit
	if hasMore {
		notifications = notifications[:limit]
	}
	return notifications, hasMore, nil
}

// GetNotification returns a single notification owned by the recipient, or nil if there is none
func (r *Repository) GetNotification(recipientID, notificationID int64) (*models.Notification, error) {
	query := `
		SELECT ` + notificationColumns + notificationFrom + `
		WHERE n.id = $1 AND n.recipient_id = $2
	`
	return r.queryOneNotification(query, notificationID, recipientID)
}

// GetNotificationByActorType returns the latest notification for a recipient +
// actor + type, or nil if there is none
func (r *Repository) GetNotificationByActorType(recipientID, actorID int64, notifType string) (*models.Notification, error) {
	query := `
		SELECT ` + notificationColumns + notificationFrom + `
		WHERE n.recipient_id = $1 AND n.actor_id = $2 AND n.type = $3
		ORDER BY n.id DESC
		LIMIT 1
	`
	return r.queryOneNotification(query, recipientID, actorID, notifType)
}

// GetUnreadCount returns the number of unread, non-expired notifications
func (r *Repository) GetUnreadCount(recipientID int64) (int64, error) {
	query := `
		SELECT COUNT(*) FROM notifications
		WHERE recipient_id = $1 AND is_read = 0 AND is_expired = 0
	`
	var count int64
	err := r.DB.Database.QueryRow(query, recipientID).Scan(&count)
	return count, err
}

// MarkNotificationRead marks one notification as read; reports whether the
// notification exists for that recipient
func (r *Repository) MarkNotificationRead(recipientID, notificationID int64) (bool, error) {
	query := `
		UPDATE notifications SET is_read = 1
		WHERE id = $1 AND recipient_id = $2
	`
	return r.execAffected(query, notificationID, recipientID)
}

// MarkAllNotificationsRead marks every unread notification of a recipient as read
func (r *Repository) MarkAllNotificationsRead(recipientID int64) error {
	query := `
		UPDATE notifications SET is_read = 1
		WHERE recipient_id = $1 AND is_read = 0 AND is_expired = 0
	`
	_, err := r.DB.Database.Exec(query, recipientID)
	return err
}

// ExpireNotification marks a notification as expired. Expiring an already
// expired notification is a no-op, so callers can retry safely.
func (r *Repository) ExpireNotification(recipientID, notificationID int64) error {
	query := `
		UPDATE notifications SET is_expired = 1
		WHERE id = $1 AND recipient_id = $2
	`
	_, err := r.DB.Database.Exec(query, notificationID, recipientID)
	return err
}

// ExpireNotificationsByType marks all unread, non-expired notifications of a type as expired
func (r *Repository) ExpireNotificationsByType(recipientID int64, notifType string) error {
	query := `
		UPDATE notifications SET is_expired = 1
		WHERE recipient_id = $1 AND type = $2 AND is_read = 0 AND is_expired = 0
	`
	_, err := r.DB.Database.Exec(query, recipientID, notifType)
	return err
}

func (r *Repository) execAffected(query string, args ...interface{}) (bool, error) {
	res, err := r.DB.Database.Exec(query, args...)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

// notificationFields lists the scan targets for notificationColumns, in the
// same order. Both the single-row and the paged query scan through it, so the
// column list and the destinations can only ever change together.
func notificationFields(n *models.Notification) []interface{} {
	return []interface{}{
		&n.ID, &n.RecipientID, &n.ActorID, &n.Type, &n.Title, &n.Message,
		&n.Payload, &n.Actions, &n.IsRead, &n.IsExpired, &n.CreatedAt,
		&n.ActorUsername, &n.ActorAvatar,
	}
}

// queryOneNotification runs a single-row notification query; no match is
// reported as (nil, nil) rather than sql.ErrNoRows.
func (r *Repository) queryOneNotification(query string, args ...interface{}) (*models.Notification, error) {
	var n models.Notification
	err := r.DB.Database.QueryRow(query, args...).Scan(notificationFields(&n)...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}
