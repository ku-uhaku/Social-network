package service

import (
	"database/sql"
	"errors"
	"time"

	"kuu/internal/models"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrUnauthorized       = errors.New("you are not authorized to perform this action on this group")
	ErrAlreadyMember      = errors.New("user is already a member or has a pending request")
	ErrNotGroupCreator    = errors.New("only group creator can perform this action")
	ErrNotMember          = errors.New("you must be an accepted member of the group to perform this action")
	ErrNoPendingInvite    = errors.New("no pending invitation or join request found")
	ErrNothingToInvite    = errors.New("no users to invite: everyone is already a member or pending")
	ErrCreatorCannotLeave = errors.New("group creator cannot leave; delete the group instead")
	ErrEventNotFound      = errors.New("event not found")
	ErrEventExpired       = errors.New("event is no longer upcoming")
)

func (s *Service) CreateGroup(creatorID int64, payload models.CreateGroupPayload) (*models.Group, error) {
	return s.Repo.CreateGroup(creatorID, payload)
}

func (s *Service) GetGroupByID(groupID int64) (*models.Group, error) {
	group, err := s.Repo.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	return group, nil
}

func (s *Service) GetAllGroups() ([]models.Group, error) {
	return s.Repo.GetAllGroups()
}

func (s *Service) UpdateGroup(userID int64, groupID int64, payload models.UpdateGroupPayload) (*models.Group, error) {
	// 1. Check if group exists and verify creator ownership
	existingGroup, err := s.Repo.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	if existingGroup.CreatorID != userID {
		return nil, ErrUnauthorized
	}

	return s.Repo.UpdateGroup(groupID, payload)
}

func (s *Service) DeleteGroup(userID int64, groupID int64) error {
	// 1. Check if group exists and verify creator ownership
	existingGroup, err := s.Repo.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGroupNotFound
		}
		return err
	}

	if existingGroup.CreatorID != userID {
		return ErrUnauthorized
	}

	return s.Repo.DeleteGroup(groupID)
}

// Invites users; requester must be an accepted member. Self and users already
// in the group (or pending) are skipped, and only the actually-invited ids are
// returned.
func (s *Service) InviteUsers(requesterID int64, payload models.InviteMembersPayload) ([]int64, error) {
	status, err := s.Repo.GetMemberStatus(payload.GroupID, requesterID)
	if err != nil || status != "accepted" {
		return nil, ErrNotMember
	}

	invited := make([]int64, 0, len(payload.TargetUserIDs))
	for _, uid := range payload.TargetUserIDs {
		if uid == requesterID {
			continue
		}
		if memberStatus, err := s.Repo.GetMemberStatus(payload.GroupID, uid); err == nil && memberStatus != "" {
			continue // already accepted, pending, etc.
		}
		invited = append(invited, uid)
	}
	if len(invited) == 0 {
		return nil, ErrNothingToInvite
	}
	if err := s.Repo.InviteUsersBatch(payload.GroupID, invited); err != nil {
		return nil, err
	}
	return invited, nil
}

// Accepts or declines a group invite
func (s *Service) RespondToInvitation(userID int64, groupID int64, accept bool) error {
	status, err := s.Repo.GetMemberStatus(groupID, userID)
	if err != nil || status != "pending" {
		return ErrNoPendingInvite
	}

	if accept {
		return s.Repo.UpdateMemberStatus(groupID, userID, "accepted")
	}
	return s.Repo.RemoveMember(groupID, userID)
}

// Public: join; private: request to join
func (s *Service) JoinGroup(userID int64, groupID int64) error {
	
	status, err := s.Repo.GetMemberStatus(groupID, userID)
	if err == nil && (status == "accepted" || status == "pending") {
		return ErrAlreadyMember
	}

	// For private group: create a pending join request
	return s.Repo.AddMember(groupID, userID, "pending")
}

// Creator approves or rejects join requests (only actual pending requests)
func (s *Service) HandleJoinRequest(creatorID int64, groupID int64, targetUserID int64, approve bool) error {
	group, err := s.Repo.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if group.CreatorID != creatorID {
		return ErrNotGroupCreator
	}
	targetStatus, err := s.Repo.GetMemberStatus(groupID, targetUserID)
	if err != nil || targetStatus != "pending" {
		return ErrNoPendingInvite
	}
	if approve {
		return s.Repo.UpdateMemberStatus(groupID, targetUserID, "accepted")
	}
	return s.Repo.RemoveMember(groupID, targetUserID)
}

// Creator cannot leave (that would delete the whole group); use DeleteGroup.
func (s *Service) LeaveGroup(userID int64, groupID int64) error {
	group, err := s.Repo.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if group.CreatorID == userID {
		return ErrCreatorCannotLeave
	}
	return s.Repo.RemoveMember(groupID, userID)
}

func (s *Service) GetPendingInvitations(userID int64) ([]models.GroupInvitationView, error) {
	return s.Repo.GetUserPendingInvitations(userID)
}

// Membership status: accepted, pending, or none
func (s *Service) GetMembershipStatus(groupID int64, userID int64) (string, error) {
	status, err := s.Repo.GetMemberStatus(groupID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "none", nil
		}
		return "", err
	}
	return status, nil
}

func (s *Service) GetGroupMembers(groupID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetGroupMembers(groupID)
}

// Group posts; non-members get ErrAccessDenied
func (s *Service) GetGroupFeed(userID int64, groupID int64, limit int, cursor *int64) ([]models.Post, bool, error) {
	status, err := s.Repo.GetMemberStatus(groupID, userID)
	if err != nil || status != "accepted" {
		return nil, false, ErrAccessDenied
	}

	return s.Repo.GetGroupFeedPosts(groupID, limit, cursor)
}

// Accepted members can create events
func (s *Service) CreateGroupEvent(userID int64, payload models.CreateGroupEventPayload) (*models.GroupEvent, error) {
	status, err := s.Repo.GetMemberStatus(payload.GroupID, userID)
	if err != nil || status != "accepted" {
		return nil, ErrNotMember
	}
	return s.Repo.CreateGroupEvent(payload.GroupID, userID, payload)
}

// Group events (members only)
func (s *Service) GetGroupEvents(userID, groupID int64) ([]models.GroupEventWithCounts, error) {
	status, err := s.Repo.GetMemberStatus(groupID, userID)
	if err != nil || status != "accepted" {
		return nil, ErrAccessDenied
	}
	return s.Repo.GetGroupEvents(groupID, userID)
}

// Creator cancels an upcoming event
func (s *Service) CancelGroupEvent(userID, eventID int64) error {
	event, err := s.Repo.GetEventByID(eventID)
	if err != nil {
		return ErrEventNotFound
	}
	if event.CreatorID != userID {
		return ErrNotGroupCreator
	}
	if event.Status != "upcoming" || event.EventTime.Before(time.Now()) {
		return ErrEventExpired
	}
	return s.Repo.CancelGroupEvent(eventID)
}

// Member picks going/not_going for an upcoming event
func (s *Service) SetEventResponse(userID, eventID int64, status string) (*models.GroupEventWithCounts, error) {
	event, err := s.Repo.GetEventByID(eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}
	memberStatus, err := s.Repo.GetMemberStatus(event.GroupID, userID)
	if err != nil || memberStatus != "accepted" {
		return nil, ErrNotMember
	}
	if event.Status != "upcoming" || event.EventTime.Before(time.Now()) {
		return nil, ErrEventExpired
	}
	if err := s.Repo.SetEventResponse(eventID, userID, status); err != nil {
		return nil, err
	}
	return s.Repo.GetEventWithCounts(eventID, userID)
}
