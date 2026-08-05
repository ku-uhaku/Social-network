package requests

import (
	"errors"

	"kuu/internal/models"
)

// ValidateMarkNotificationRead ensures the read payload is structurally sound
func ValidateMarkNotificationRead(payload models.MarkNotificationReadPayload) []error {
	var errs []error
	if !payload.All && payload.NotificationID == nil {
		errs = append(errs, errors.New("either notification_id or all must be provided"))
	}
	if payload.NotificationID != nil && *payload.NotificationID <= 0 {
		errs = append(errs, errors.New("notification_id must be a positive integer"))
	}
	return errs
}

// ValidateExpireNotification ensures the expire payload is structurally sound
func ValidateExpireNotification(payload models.ExpireNotificationPayload) []error {
	var errs []error
	if payload.NotificationID <= 0 {
		errs = append(errs, errors.New("notification_id is required and must be valid"))
	}
	return errs
}