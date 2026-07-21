package service

import (
	"context"
	"errors"

	"kuu/internal/models"
)

func (s *Service) SaveDirectMessage(ctx context.Context, senderID, receiverID int64, content string) (*models.DirectMessage, error) {
	return s.Repo.SaveDirectMessage(ctx, senderID, receiverID, content)
}

func (s *Service) SaveGroupMessage(ctx context.Context, senderID, groupID int64, content string) (*models.GroupMessage, error) {
	// Verify user is a group member
	status, err := s.Repo.GetMemberStatus(ctx, groupID, senderID)
	if err != nil || status != "accepted" {
		return nil, errors.New("unauthorized group access")
	}
	return s.Repo.SaveGroupMessage(ctx, senderID, groupID, content)
}

func (s *Service) GetDirectHistory(ctx context.Context, userA, userB int64, limit, offset int) ([]models.DirectMessage, error) {
	return s.Repo.GetDirectHistory(ctx, userA, userB, limit, offset)
}

func (s *Service) GetGroupHistory(ctx context.Context, userID, groupID int64, limit, offset int) ([]models.GroupMessage, error) {
	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err != nil || status != "accepted" {
		return nil, errors.New("unauthorized group access")
	}
	return s.Repo.GetGroupHistory(ctx, groupID, limit, offset)
}

func (s *Service) GetGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return s.Repo.GetGroupMemberIDs(ctx, groupID)
}
