package service

import (
	"errors"

	"kuu/internal/models"
)

// CreateNotification persists a notification for a recipient
func (s *Service) CreateNotification(n *models.Notification) (*models.Notification, error) {
	return s.Repo.CreateNotification(n)
}

// GetUserNotifications returns a paginated stack with the unread count.
// lastID acts as a cursor: 0 loads the newest page, otherwise only
// notifications with a lower id are returned.
func (s *Service) GetUserNotifications(recipientID int64, limit int, lastID int64) (*models.NotificationListResponse, error) {
	notifications, hasMore, err := s.Repo.GetUserNotifications(recipientID, limit, lastID)
	if err != nil {
		return nil, err
	}
	unreadCount, err := s.Repo.GetUnreadCount(recipientID)
	if err != nil {
		return nil, err
	}
	return &models.NotificationListResponse{
		Notifications: notifications,
		UnreadCount:   unreadCount,
		HasMore:       hasMore,
	}, nil
}

// MarkNotificationRead marks a single notification or all as read for a recipient
func (s *Service) MarkNotificationRead(recipientID, notificationID int64) (int64, error) {
	if notificationID > 0 {
		n, err := s.Repo.GetNotification(recipientID, notificationID)
		if err != nil || n == nil {
			return 0, errors.New("notification not found")
		}
	}
	return s.Repo.MarkNotificationRead(recipientID, notificationID)
}

// ExpireNotification marks a single notification as expired for a recipient
func (s *Service) ExpireNotification(recipientID, notificationID int64) error {
	return s.Repo.ExpireNotification(recipientID, notificationID)
}

// GetNotification returns a single notification owned by the recipient
func (s *Service) GetNotification(recipientID, notificationID int64) (*models.Notification, error) {
	return s.Repo.GetNotification(recipientID, notificationID)
}

// GetNotificationByActorType finds a notification by recipient + actor + type
func (s *Service) GetNotificationByActorType(recipientID, actorID int64, notifType string) (*models.Notification, error) {
	return s.Repo.GetNotificationByActorType(recipientID, actorID, notifType)
}

// ExpireNotificationsByType marks all unread, non-expired notifications of a type as expired
func (s *Service) ExpireNotificationsByType(recipientID int64, notifType string) error {
	return s.Repo.ExpireNotificationsByType(recipientID, notifType)
}
