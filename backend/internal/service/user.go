package service

import (
	"context"

	"kuu/internal/models"
)

// UpdateProfile orchestrates incoming structural profile changes
func (s *Service) UpdateProfile(ctx context.Context, userID int64, payload models.UpdateProfilePayload) (*models.User, error) {
	return s.Repo.UpdateUserProfile(ctx, userID, payload)
}
