package service

import (
	"context"
	"database/sql"
	"errors"

	"kuu/internal/models"
)

var (
	ErrGroupNotFound   = errors.New("group not found")
	ErrUnauthorized    = errors.New("you are not authorized to perform this action on this group")
	ErrAlreadyMember   = errors.New("user is already a member or has a pending request")
	ErrNotGroupCreator = errors.New("only group creator can perform this action")
	ErrGroupIsPrivate  = errors.New("cannot directly join a private group; request required")
	ErrNoPendingInvite = errors.New("no pending invitation found")
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

// InviteUsers ensures requester is group creator then issues batch invitations
func (s *Service) InviteUsers(ctx context.Context, requesterID int64, payload models.InviteMembersPayload) error {
	group, err := s.Repo.GetGroupByID(ctx, payload.GroupID)
	if err != nil {
		return err
	}
	if group.CreatorID != requesterID {
		return ErrNotGroupCreator
	}

	return s.Repo.InviteUsersBatch(ctx, payload.GroupID, payload.TargetUserIDs)
}

// RespondToInvitation handles accepting or declining group invites by the target user
func (s *Service) RespondToInvitation(ctx context.Context, userID int64, groupID int64, accept bool) error {
	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err != nil || status != "pending" {
		return ErrNoPendingInvite
	}

	if accept {
		return s.Repo.UpdateMemberStatus(ctx, groupID, userID, "accepted")
	}
	return s.Repo.RemoveMember(ctx, groupID, userID)
}

// JoinGroup handles public join or requesting join for private groups
func (s *Service) JoinGroup(ctx context.Context, userID int64, groupID int64) error {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return err
	}

	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err == nil && (status == "accepted" || status == "pending") {
		return ErrAlreadyMember
	}

	if group.IsPublic == 1 {
		return s.Repo.AddMember(ctx, groupID, userID, "accepted")
	}

	// For private group: create a pending join request
	return s.Repo.AddMember(ctx, groupID, userID, "pending")
}

// HandleJoinRequest lets creator approve or reject pending requests to private groups
func (s *Service) HandleJoinRequest(ctx context.Context, creatorID int64, groupID int64, targetUserID int64, approve bool) error {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.CreatorID != creatorID {
		return ErrNotGroupCreator
	}

	if approve {
		return s.Repo.UpdateMemberStatus(ctx, groupID, targetUserID, "accepted")
	}
	return s.Repo.RemoveMember(ctx, groupID, targetUserID)
}

// LeaveGroup allows a member to leave
func (s *Service) LeaveGroup(ctx context.Context, userID int64, groupID int64) error {
	return s.Repo.RemoveMember(ctx, groupID, userID)
}

func (s *Service) GetPendingInvitations(ctx context.Context, userID int64) ([]models.GroupInvitationView, error) {
	return s.Repo.GetUserPendingInvitations(ctx, userID)
}
