package service

import (
	"context"
	"errors"

	"kuu/internal/models"
)

func (s *Service) CreateGroup(ctx context.Context, creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	// Call repository to handle group instantiation and creator enrollment transactions
	return s.Repo.CreateGroupWithCreator(ctx, creatorID, payload)
}

func (s *Service) GetPublicGroups(ctx context.Context) ([]models.Group, error) {
	return s.Repo.FetchAllPublicGroups(ctx)
}

func (s *Service) InviteMembersToGroup(ctx context.Context, requesterID int64, groupID int64, payload models.InviteMembersPayload) error {
	if len(payload.TargetUsersID) == 0 {
		return nil // Quick return if no users specified
	}

	// 1. Authorization check: Is the requester a member/creator of this group?
	isMember, err := s.Repo.IsGroupMember(ctx, requesterID, groupID)
	if err != nil {
		return err
	}
	if !isMember {
		// FIXED: Replaced models.ErrUnauthorizedAction with a standard text error
		return errors.New("unauthorized action: you are not a member of this group")
	}

	// 2. Delegate database bulk insertion to repository
	return s.Repo.AddGroupMembersPending(ctx, groupID, payload.TargetUsersID)
}

func (s *Service) AcceptInvitation(ctx context.Context, userID int64, groupID int64) error {
	// You can add validation logic here if needed
	return s.Repo.AcceptGroupInvitation(ctx, userID, groupID)
}

func (s *Service) JoinPublicGroup(ctx context.Context, userID int64, groupID int64) error {
	return s.Repo.JoinPublicGroup(ctx, userID, groupID)
}