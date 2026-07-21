package service

import (
	"context"
	"database/sql"
	"errors"

	"kuu/internal/models"
)

var (
	ErrCannotFollowSelf  = errors.New("you cannot follow yourself")
	ErrAlreadyFollowing  = errors.New("you are already following or requested to follow this user")
	ErrFollowReqNotFound = errors.New("no pending follow request found")
)

// UpdateProfile orchestrates incoming structural profile changes
func (s *Service) UpdateProfile(ctx context.Context, userID int64, payload models.UpdateProfilePayload) (*models.User, error) {
	return s.Repo.UpdateUserProfile(ctx, userID, payload)
}

func (s *Service) FollowUser(ctx context.Context, followerID, targetID int64) (string, error) {
	if followerID == targetID {
		return "", ErrCannotFollowSelf
	}

	// 1. Verify target user exists
	targetUser, err := s.Repo.GetUserByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("target user does not exist")
		}
		return "", err
	}

	// 2. Check current status
	status, err := s.Repo.GetFollowRelation(ctx, followerID, targetID)
	if err == nil && (status == "accepted" || status == "pending") {
		return "", ErrAlreadyFollowing
	}

	// 3. Determine status based on privacy
	initialStatus := "accepted"
	if targetUser.IsPublic == 0 {
		initialStatus = "pending"
	}

	if err := s.Repo.InsertFollowRelation(ctx, followerID, targetID, initialStatus); err != nil {
		return "", err
	}

	return initialStatus, nil
}

// UnfollowUser handles unfollowing or cancelling a pending request
func (s *Service) UnfollowUser(ctx context.Context, followerID, targetID int64) error {
	return s.Repo.RemoveFollowRelation(ctx, followerID, targetID)
}

// HandleFollowRequest accepts or declines an incoming request (called by target user)
func (s *Service) HandleFollowRequest(ctx context.Context, targetUserID, requesterID int64, accept bool) error {
	status, err := s.Repo.GetFollowRelation(ctx, requesterID, targetUserID)
	if err != nil || status != "pending" {
		return ErrFollowReqNotFound
	}

	if accept {
		return s.Repo.UpdateFollowStatus(ctx, requesterID, targetUserID, "accepted")
	}
	return s.Repo.RemoveFollowRelation(ctx, requesterID, targetUserID)
}

func (s *Service) GetFollowers(ctx context.Context, userID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetFollowers(ctx, userID)
}

func (s *Service) GetFollowing(ctx context.Context, userID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetFollowing(ctx, userID)
}

func (s *Service) GetPendingRequests(ctx context.Context, userID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetPendingFollowRequests(ctx, userID)
}
