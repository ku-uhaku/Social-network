package service

import (
	"context"
	"database/sql"
	"errors"

	"kuu/internal/models"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrUnauthorized  = errors.New("you are not authorized to perform this action on this group")
)

func (s *Service) CreateGroup(ctx context.Context, creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	return s.Repo.CreateGroup(ctx, creatorID, payload)
}

func (s *Service) GetGroupByID(ctx context.Context, groupID int64) (*models.Group, error) {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	return group, nil
}

func (s *Service) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return s.Repo.GetAllGroups(ctx)
}

func (s *Service) UpdateGroup(ctx context.Context, userID int64, groupID int64, payload models.UpdateGroupPayload) (*models.Group, error) {
	// 1. Check if group exists and verify creator ownership
	existingGroup, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	if existingGroup.CreatorID != userID {
		return nil, ErrUnauthorized
	}

	return s.Repo.UpdateGroup(ctx, groupID, payload)
}

func (s *Service) DeleteGroup(ctx context.Context, userID int64, groupID int64) error {
	// 1. Check if group exists and verify creator ownership
	existingGroup, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGroupNotFound
		}
		return err
	}

	if existingGroup.CreatorID != userID {
		return ErrUnauthorized
	}

	return s.Repo.DeleteGroup(ctx, groupID)
}
