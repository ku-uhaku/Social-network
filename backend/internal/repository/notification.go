package repository

import (
	"context"

	"kuu/internal/models"
)

const notificationSelectColumns = `
	n.id, n.recipient_id, n.actor_id, n.type, n.title, n.message,
	n.payload, n.actions, n.is_read, n.is_expired, n.created_at,
	COALESCE(u.username, ''), u.avatar
`

// CreateNotification persists a new notification for a recipient
func (r *Repository) CreateNotification(ctx context.Context, n *models.Notification) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (recipient_id, actor_id, type, title, message, payload, actions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	res, err := r.DB.Database.ExecContext(ctx, query,
		n.RecipientID, n.ActorID, n.Type, n.Title, n.Message, n.Payload, n.Actions)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetNotification(ctx, n.RecipientID, id)
}

// GetUserNotifications returns a paginated stack for a recipient, newest first
func (r *Repository) GetUserNotifications(ctx context.Context, recipientID int64, limit, offset int) ([]models.Notification, error) {
	query := `
		SELECT ` + notificationSelectColumns + `
		FROM notifications n
		LEFT JOIN users u ON u.id = n.actor_id
		WHERE n.recipient_id = $1
		ORDER BY n.created_at DESC, n.id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Database.QueryContext(ctx, query, recipientID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		n, err := r.scanNotification(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, *n)
	}
	return notifications, nil
}

// GetUnreadCount returns the number of unread, non-expired notifications
func (r *Repository) GetUnreadCount(ctx context.Context, recipientID int64) (int64, error) {
	query := `
		SELECT COUNT(*) FROM notifications
		WHERE recipient_id = $1 AND is_read = 0 AND is_expired = 0
	`
	var count int64
	err := r.DB.Database.QueryRowContext(ctx, query, recipientID).Scan(&count)
	return count, err
}

// MarkNotificationRead marks a single notification (by id) or all as read for a recipient.
// Returns the id of the single notification if it was marked, or 0 when marking all.
func (r *Repository) MarkNotificationRead(ctx context.Context, recipientID, notificationID int64) (int64, error) {
	if notificationID > 0 {
		query := `
			UPDATE notifications SET is_read = 1
			WHERE id = $1 AND recipient_id = $2 AND is_expired = 0
		`
		_, err := r.DB.Database.ExecContext(ctx, query, notificationID, recipientID)
		return notificationID, err
	}

	query := `
		UPDATE notifications SET is_read = 1
		WHERE recipient_id = $1 AND is_read = 0 AND is_expired = 0
	`
	_, err := r.DB.Database.ExecContext(ctx, query, recipientID)
	return 0, err
}

// ExpireNotification marks a notification as expired for a recipient
func (r *Repository) ExpireNotification(ctx context.Context, recipientID, notificationID int64) error {
	query := `
		UPDATE notifications SET is_expired = 1
		WHERE id = $1 AND recipient_id = $2 AND is_expired = 0
	`
	_, err := r.DB.Database.ExecContext(ctx, query, notificationID, recipientID)
	return err
}

// GetNotification returns a single notification if it belongs to the recipient
func (r *Repository) GetNotification(ctx context.Context, recipientID, notificationID int64) (*models.Notification, error) {
	query := `
		SELECT ` + notificationSelectColumns + `
		FROM notifications n
		LEFT JOIN users u ON u.id = n.actor_id
		WHERE n.id = $1 AND n.recipient_id = $2
	`
	return r.scanNotification(r.DB.Database.QueryRowContext(ctx, query, notificationID, recipientID))
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func (r *Repository) scanNotification(row rowScanner) (*models.Notification, error) {
	var n models.Notification
	err := row.Scan(
		&n.ID, &n.RecipientID, &n.ActorID, &n.Type, &n.Title, &n.Message,
		&n.Payload, &n.Actions, &n.IsRead, &n.IsExpired, &n.CreatedAt,
		&n.ActorUsername, &n.ActorAvatar,
	)
	if err != nil {
		return nil, err
	}
	if n.ActorUsername != nil && *n.ActorUsername == "" {
		n.ActorUsername = nil
	}
	return &n, nil
}

