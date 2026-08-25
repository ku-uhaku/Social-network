package service

import (
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
func (s *Service) UpdateProfile(userID int64, payload models.UpdateProfilePayload) (*models.User, error) {
	return s.Repo.UpdateUserProfile(userID, payload)
}

// AcceptAllPendingFollows accepts all pending follow requests for a user
func (s *Service) AcceptAllPendingFollows(targetUserID int64) ([]int64, error) {
	return s.Repo.AcceptAllPendingFollows(targetUserID)
}

// GetUserProfile fetches a user profile by username, enriched with follow stats
// and the requesting viewer's relationship to the target user.
func (s *Service) GetUserProfile(viewerID int64, username string) (*models.UserProfileView, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	stats, err := s.Repo.GetFollowStats(user.ID)
	if err != nil {
		return nil, err
	}

	view := &models.UserProfileView{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Gender:         user.Gender,
		DateOfBirth:    user.DateOfBirth,
		IsPublic:       user.IsPublic,
		Avatar:         user.Avatar,
		AboutMe:        user.AboutMe,
		CreatedAt:      user.CreatedAt,
		FollowersCount: stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
		FollowStatus:   "none",
	}

	if viewerID == user.ID {
		view.FollowStatus = "self"
		return view, nil
	}

	status, err := s.Repo.GetFollowRelation(viewerID, user.ID)
	if err == nil && (status == "pending" || status == "accepted") {
		view.FollowStatus = status
	}

	// Private profiles: strip personal fields unless the viewer is a follower.
	if user.IsPublic == 0 && view.FollowStatus != "accepted" {
		view.Email = ""
		view.FirstName = ""
		view.LastName = ""
		view.Gender = ""
		view.DateOfBirth = ""
		view.AboutMe = nil
	}

	return view, nil
}

// GetUserPosts retrieves a user's posts (excluding group posts) with author metadata
func (s *Service) GetUserPosts(targetUserID int64, viewerID int64) ([]models.Post, error) {
	return s.Repo.GetUserPosts(targetUserID, viewerID)
}

func (s *Service) FollowUser(followerID, targetID int64) (string, error) {
	if followerID == targetID {
		return "", ErrCannotFollowSelf
	}

	// 1. Verify target user exists
	targetUser, err := s.Repo.GetUserByID(targetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("target user does not exist")
		}
		return "", err
	}

	// 2. Check current status
	status, err := s.Repo.GetFollowRelation(followerID, targetID)
	if err == nil && (status == "accepted" || status == "pending") {
		return "", ErrAlreadyFollowing
	}

	// 3. Determine status based on privacy
	initialStatus := "accepted"
	if targetUser.IsPublic == 0 {
		initialStatus = "pending"
	}

	if err := s.Repo.InsertFollowRelation(followerID, targetID, initialStatus); err != nil {
		return "", err
	}

	return initialStatus, nil
}

// UnfollowUser handles unfollowing or cancelling a pending request
func (s *Service) UnfollowUser(followerID, targetID int64) error {
	return s.Repo.RemoveFollowRelation(followerID, targetID)
}

// HandleFollowRequest accepts or declines an incoming request (called by target user).
// if already accepted, returns nil; if no row exists, returns nil.
func (s *Service) HandleFollowRequest(targetUserID, requesterID int64, accept bool) error {
	status, err := s.Repo.GetFollowRelation(requesterID, targetUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	if status == "accepted" {
		return nil
	}
	if status != "pending" {
		return ErrFollowReqNotFound
	}

	if accept {
		return s.Repo.UpdateFollowStatus(requesterID, targetUserID, "accepted")
	}
	return s.Repo.RemoveFollowRelation(requesterID, targetUserID)
}

func (s *Service) GetFollowers(userID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetFollowers(userID)
}

func (s *Service) GetFollowing(userID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetFollowing(userID)
}

func (s *Service) GetPendingRequests(userID int64) ([]models.UserFollowView, error) {
	return s.Repo.GetPendingFollowRequests(userID)
}

func (s *Service) GetAllUsers() ([]models.UserFollowView, error) {
	return s.Repo.GetAllUsers()
}

// GetSuggestedUsers returns a shortlist of users for the viewer to follow
func (s *Service) GetSuggestedUsers(viewerID int64, limit int) ([]models.UserFollowView, error) {
	return s.Repo.GetSuggestedUsers(viewerID, limit)
}
