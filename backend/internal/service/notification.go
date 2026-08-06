package service

import (
	"context"
	"errors"

	"kuu/internal/models"
)

// CreateNotification persists a notification for a recipient
func (s *Service) CreateNotification(ctx context.Context, n *models.Notification) (*models.Notification, error) {
	return s.Repo.CreateNotification(ctx, n)
}

// GetUserNotifications returns a paginated stack with the unread count
func (s *Service) GetUserNotifications(ctx context.Context, recipientID int64, limit, offset int) (*models.NotificationListResponse, error) {
	notifications, err := s.Repo.GetUserNotifications(ctx, recipientID, limit, offset)
	if err != nil {
		return nil, err
	}
	unreadCount, err := s.Repo.GetUnreadCount(ctx, recipientID)
	if err != nil {
		return nil, err
	}
	return &models.NotificationListResponse{
		Notifications: notifications,
		UnreadCount:  unreadCount,
	}, nil
}

// MarkNotificationRead marks a single notification or all as read for a recipient
func (s *Service) MarkNotificationRead(ctx context.Context, recipientID, notificationID int64) (int64, error) {
	if notificationID > 0 {
		n, err := s.Repo.GetNotification(ctx, recipientID, notificationID)
		if err != nil || n == nil {
			return 0, errors.New("notification not found")
		}
	}
	return s.Repo.MarkNotificationRead(ctx, recipientID, notificationID)
}

// ExpireNotification marks a single notification as expired for a recipient
func (s *Service) ExpireNotification(ctx context.Context, recipientID, notificationID int64) error {
	return s.Repo.ExpireNotification(ctx, recipientID, notificationID)
}

// GetNotification returns a single notification owned by the recipient
func (s *Service) GetNotification(ctx context.Context, recipientID, notificationID int64) (*models.Notification, error) {
	return s.Repo.GetNotification(ctx, recipientID, notificationID)
}

// GetNotificationByActorType finds a notification by recipient + actor + type
func (s *Service) GetNotificationByActorType(ctx context.Context, recipientID, actorID int64, notifType string) (*models.Notification, error) {
	return s.Repo.GetNotificationByActorType(ctx, recipientID, actorID, notifType)
}

// ExpireNotificationsByType marks all unread, non-expired notifications of a type as expired
func (s *Service) ExpireNotificationsByType(ctx context.Context, recipientID int64, notifType string) error {
	return s.Repo.ExpireNotificationsByType(ctx, recipientID, notifType)
}
