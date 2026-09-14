package service

import (
	"errors"

	"kuu/internal/models"
)

func (s *Service) CreateNotification(n *models.Notification) (*models.Notification, error) {
	return s.Repo.CreateNotification(n)
}

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

func (s *Service) MarkAllNotificationsRead(recipientID int64) error {
	return s.Repo.MarkAllNotificationsRead(recipientID)
}

func (s *Service) ExpireNotification(recipientID, notificationID int64) error {
	return s.Repo.ExpireNotification(recipientID, notificationID)
}

func (s *Service) GetNotificationByActorType(recipientID, actorID int64, notifType string) (*models.Notification, error) {
	return s.Repo.GetNotificationByActorType(recipientID, actorID, notifType)
}

func (s *Service) ExpireNotificationsByType(recipientID int64, notifType string) ([]int64, error) {
	return s.Repo.ExpireNotificationsByType(recipientID, notifType)
}

func (s *Service) ExpireNotificationsByActorType(recipientID, actorID int64, notifType string) ([]int64, error) {
	return s.Repo.ExpireNotificationsByActorType(recipientID, actorID, notifType)
}

func (s *Service) ExpireGroupNotifications(recipientID, groupID int64, notifType string) ([]int64, error) {
	return s.Repo.ExpireGroupNotifications(recipientID, groupID, notifType)
}
