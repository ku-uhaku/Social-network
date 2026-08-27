package service

import (
	"errors"

	"kuu/internal/models"
)

// CreateNotification persists a notification for a recipient
func (s *Service) CreateNotification(n *models.Notification) (*models.Notification, error) {
	return s.Repo.CreateNotification(n)
}

// GetUserNotifications returns a page of notifications plus the unread count;
// lastID is the cursor (0 = newest page)
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

// MarkNotificationRead marks one notification owned by the recipient as read
func (s *Service) MarkNotificationRead(recipientID, notificationID int64) error {
	found, err := s.Repo.MarkNotificationRead(recipientID, notificationID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("notification not found")
	}
	return nil
}

// MarkAllNotificationsRead marks every unread notification of a recipient as read
func (s *Service) MarkAllNotificationsRead(recipientID int64) error {
	return s.Repo.MarkAllNotificationsRead(recipientID)
}

// ExpireNotification marks a single notification as expired for a recipient
func (s *Service) ExpireNotification(recipientID, notificationID int64) error {
	return s.Repo.ExpireNotification(recipientID, notificationID)
}

// GetNotificationByActorType returns the latest notification for a recipient +
// actor + type, or nil if there is none
func (s *Service) GetNotificationByActorType(recipientID, actorID int64, notifType string) (*models.Notification, error) {
	return s.Repo.GetNotificationByActorType(recipientID, actorID, notifType)
}

// ExpireNotificationsByType expires all unread notifications of a type
func (s *Service) ExpireNotificationsByType(recipientID int64, notifType string) error {
	return s.Repo.ExpireNotificationsByType(recipientID, notifType)
}
