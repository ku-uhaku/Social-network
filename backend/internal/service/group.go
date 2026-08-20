package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"kuu/internal/models"
)

var (
	ErrGroupNotFound   = errors.New("group not found")
	ErrUnauthorized    = errors.New("you are not authorized to perform this action on this group")
	ErrAlreadyMember   = errors.New("user is already a member or has a pending request")
	ErrNotGroupCreator = errors.New("only group creator can perform this action")
	ErrNotMember       = errors.New("you must be an accepted member of the group to perform this action")
	ErrGroupIsPrivate  = errors.New("cannot directly join a private group; request required")
	ErrNoPendingInvite = errors.New("no pending invitation found")
	ErrEventNotFound   = errors.New("event not found")
	ErrEventExpired    = errors.New("event is no longer upcoming")
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

// InviteUsers ensures requester is an accepted member then issues batch invitations.
// Any member of the group may invite other users.
func (s *Service) InviteUsers(ctx context.Context, requesterID int64, payload models.InviteMembersPayload) error {
	status, err := s.Repo.GetMemberStatus(ctx, payload.GroupID, requesterID)
	if err != nil || status != "accepted" {
		return ErrNotMember
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
	// group, err := s.Repo.GetGroupByID(ctx, groupID)
	// if err != nil {
	// 	return err
	// }

	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err == nil && (status == "accepted" || status == "pending") {
		return ErrAlreadyMember
	}

	// if group.IsPublic == 1 {
	// 	return s.Repo.AddMember(ctx, groupID, userID, "accepted")
	// }

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

// GetMembershipStatus returns the current user's membership status in a group.
// Values: 'accepted', 'pending', or 'none' when no row exists.
func (s *Service) GetMembershipStatus(ctx context.Context, groupID int64, userID int64) (string, error) {
	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "none", nil
		}
		return "", err
	}
	return status, nil
}

// GetGroupMembers returns all accepted members of a group
func (s *Service) GetGroupMembers(ctx context.Context, groupID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetGroupMembers(ctx, groupID)
}

// GetGroupFeed returns the posts of a single group. Non-members get ErrAccessDenied.
func (s *Service) GetGroupFeed(ctx context.Context, userID int64, groupID int64, limit int, cursor *int64) ([]models.Post, bool, error) {
	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err != nil || status != "accepted" {
		return nil, false, ErrAccessDenied
	}

	return s.Repo.GetGroupFeedPosts(ctx, groupID, limit, cursor)
}

// CreateGroupEvent allows any accepted member to create an event
func (s *Service) CreateGroupEvent(ctx context.Context, userID int64, payload models.CreateGroupEventPayload) (*models.GroupEvent, error) {
	status, err := s.Repo.GetMemberStatus(ctx, payload.GroupID, userID)
	if err != nil || status != "accepted" {
		return nil, ErrNotMember
	}
	return s.Repo.CreateGroupEvent(ctx, payload.GroupID, userID, payload)
}

// GetGroupEvents lists events for a group (members only)
func (s *Service) GetGroupEvents(ctx context.Context, userID, groupID int64) ([]models.GroupEventWithCounts, error) {
	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err != nil || status != "accepted" {
		return nil, ErrAccessDenied
	}
	return s.Repo.GetGroupEvents(ctx, groupID, userID)
}

// CancelGroupEvent allows only the creator to cancel an upcoming event
func (s *Service) CancelGroupEvent(ctx context.Context, userID, eventID int64) error {
	event, err := s.Repo.GetEventByID(ctx, eventID)
	if err != nil {
		return ErrEventNotFound
	}
	if event.CreatorID != userID {
		return ErrNotGroupCreator
	}
	if event.Status != "upcoming" || event.EventTime.Before(time.Now()) {
		return ErrEventExpired
	}
	return s.Repo.CancelGroupEvent(ctx, eventID)
}

// SetEventResponse lets a member choose going/not_going for an upcoming event
func (s *Service) SetEventResponse(ctx context.Context, userID, eventID int64, status string) (*models.GroupEventWithCounts, error) {
	event, err := s.Repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}
	memberStatus, err := s.Repo.GetMemberStatus(ctx, event.GroupID, userID)
	if err != nil || memberStatus != "accepted" {
		return nil, ErrNotMember
	}
	if event.Status != "upcoming" || event.EventTime.Before(time.Now()) {
		return nil, ErrEventExpired
	}
	if err := s.Repo.SetEventResponse(ctx, eventID, userID, status); err != nil {
		return nil, err
	}
	return s.Repo.GetEventWithCounts(ctx, eventID, userID)
}
