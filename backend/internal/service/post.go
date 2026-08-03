package service

import (
	"context"
	"errors"

	"kuu/internal/models"
)

var ErrAccessDenied = errors.New("you do not have permission to access or post to this context")

func (s *Service) CreatePost(ctx context.Context, userID int64, payload models.CreatePostPayload) (*models.Post, error) {
	// If posting to a group, ensure user is an accepted member
	if payload.GroupID != nil && *payload.GroupID > 0 {
		status, err := s.Repo.GetMemberStatus(ctx, *payload.GroupID, userID)
		if err != nil || status != "accepted" {
			return nil, ErrAccessDenied
		}
		payload.Privacy = "group"
	}

	return s.Repo.CreatePost(ctx, userID, payload)
}

func (s *Service) GetFeed(ctx context.Context, userID int64) ([]models.Post, error) {
	return s.Repo.GetFeedPosts(ctx, userID)
}

func (s *Service) AddComment(ctx context.Context, userID int64, payload models.CreateCommentPayload) (*models.Comment, error) {
	// Check post existence and permission
	post, err := s.Repo.GetPostByID(ctx, payload.PostID)
	if err != nil {
		return nil, err
	}

	if err := s.checkPostVisibility(ctx, userID, post); err != nil {
		return nil, err
	}

	return s.Repo.CreateComment(ctx, userID, payload)
}

func (s *Service) GetComments(ctx context.Context, userID int64, postID int64) ([]models.Comment, error) {
	post, err := s.Repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if err := s.checkPostVisibility(ctx, userID, post); err != nil {
		return nil, err
	}

	return s.Repo.GetPostComments(ctx, postID, userID)
}

func (s *Service) checkPostVisibility(ctx context.Context, userID int64, post *models.Post) error {
	if post.UserID == userID || post.Privacy == "public" {
		return nil
	}
	if post.Privacy == "followers" {
		status, err := s.Repo.GetFollowRelation(ctx, userID, post.UserID)
		if err == nil && status == "accepted" {
			return nil
		}
	}
	if post.Privacy == "group" && post.GroupID != nil {
		status, err := s.Repo.GetMemberStatus(ctx, *post.GroupID, userID)
		if err == nil && status == "accepted" {
			return nil
		}
	}
	return ErrAccessDenied
}
