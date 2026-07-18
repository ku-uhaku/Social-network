package service

import (
	"context"

	"kuu/internal/models"
)

func (s *Service) CreateGroup(ctx context.Context, creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	// Call repository to handle group instantiation and creator enrollment transactions
	return s.Repo.CreateGroupWithCreator(ctx, creatorID, payload)
}

func (s *Service) GetPublicGroups(ctx context.Context) ([]models.Group, error) {
	return s.Repo.FetchAllPublicGroups(ctx)
}
