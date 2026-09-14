package repository

import (
	"database/sql"
	"errors"

	"kuu/internal/models"
)

const notificationColumns = `
	n.id, n.recipient_id, n.actor_id, n.type, n.title, n.message,
	n.group_id, n.is_read, n.is_expired, n.created_at,
	u.username, u.avatar
`

const notificationFrom = `
	FROM notifications n
	LEFT JOIN users u ON u.id = n.actor_id
`

func (r *Repository) CreateNotification(n *models.Notification) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (recipient_id, actor_id, type, title, message, group_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	res, err := r.DB.Database.Exec(query,
		n.RecipientID, n.ActorID, n.Type, n.Title, n.Message, n.Payload.GroupID)
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
		n.Actions = models.NotificationActionsFor(n.Type)
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

func (r *Repository) GetNotification(recipientID, notificationID int64) (*models.Notification, error) {
	query := `
		SELECT ` + notificationColumns + notificationFrom + `
		WHERE n.id = $1 AND n.recipient_id = $2
	`
	return r.queryOneNotification(query, notificationID, recipientID)
}

func (r *Repository) GetNotificationByActorType(recipientID, actorID int64, notifType string) (*models.Notification, error) {
	query := `
		SELECT ` + notificationColumns + notificationFrom + `
		WHERE n.recipient_id = $1 AND n.actor_id = $2 AND n.type = $3
		ORDER BY n.id DESC
		LIMIT 1
	`
	return r.queryOneNotification(query, recipientID, actorID, notifType)
}

func (r *Repository) GetUnreadCount(recipientID int64) (int64, error) {
	query := `
		SELECT COUNT(*) FROM notifications
		WHERE recipient_id = $1 AND is_read = 0 AND is_expired = 0
	`
	var count int64
	err := r.DB.Database.QueryRow(query, recipientID).Scan(&count)
	return count, err
}

func (r *Repository) MarkNotificationRead(recipientID, notificationID int64) (bool, error) {
	query := `
		UPDATE notifications SET is_read = 1
		WHERE id = $1 AND recipient_id = $2
	`
	return r.execAffected(query, notificationID, recipientID)
}

func (r *Repository) MarkAllNotificationsRead(recipientID int64) error {
	query := `
		UPDATE notifications SET is_read = 1
		WHERE recipient_id = $1 AND is_read = 0 AND is_expired = 0
	`
	_, err := r.DB.Database.Exec(query, recipientID)
	return err
}

func (r *Repository) ExpireNotification(recipientID, notificationID int64) error {
	query := `
		UPDATE notifications SET is_expired = 1
		WHERE id = $1 AND recipient_id = $2
	`
	_, err := r.DB.Database.Exec(query, notificationID, recipientID)
	return err
}

func (r *Repository) ExpireNotificationsByType(recipientID int64, notifType string) ([]int64, error) {
	return r.expireNotifications(
		`recipient_id = $1 AND type = $2 AND is_read = 0`,
		recipientID, notifType,
	)
}

func (r *Repository) ExpireNotificationsByActorType(recipientID, actorID int64, notifType string) ([]int64, error) {
	return r.expireNotifications(
		`recipient_id = $1 AND actor_id = $2 AND type = $3`,
		recipientID, actorID, notifType,
	)
}

func (r *Repository) ExpireGroupNotifications(recipientID, groupID int64, notifType string) ([]int64, error) {
	return r.expireNotifications(
		`recipient_id = $1 AND type = $2 AND group_id = $3`,
		recipientID, notifType, groupID,
	)
}

func (r *Repository) expireNotifications(where string, args ...interface{}) ([]int64, error) {
	query := `
		UPDATE notifications SET is_expired = 1
		WHERE is_expired = 0 AND ` + where + `
		RETURNING id
	`
	rows, err := r.DB.Database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
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

func notificationFields(n *models.Notification) []interface{} {
	return []interface{}{
		&n.ID, &n.RecipientID, &n.ActorID, &n.Type, &n.Title, &n.Message,
		&n.Payload.GroupID, &n.IsRead, &n.IsExpired, &n.CreatedAt,
		&n.ActorUsername, &n.ActorAvatar,
	}
}

func (r *Repository) queryOneNotification(query string, args ...interface{}) (*models.Notification, error) {
	var n models.Notification
	err := r.DB.Database.QueryRow(query, args...).Scan(notificationFields(&n)...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	n.Actions = models.NotificationActionsFor(n.Type)
	return &n, nil
}
