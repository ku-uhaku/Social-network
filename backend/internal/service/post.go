package service

import (
	"errors"

	"kuu/internal/models"
)

var ErrAccessDenied = errors.New("you do not have permission to access or post to this context")

func (s *Service) CreatePost(userID int64, payload models.CreatePostPayload) (*models.Post, error) {
	// ensure user is an accepted
	if payload.GroupID != nil && *payload.GroupID > 0 {
		status, err := s.Repo.GetMemberStatus(*payload.GroupID, userID)
		if err != nil || status != "accepted" {
			return nil, ErrAccessDenied
		}
	}
	if payload.Privacy == "private" {
		for _, viewerID := range payload.VisibleTo {
			status, err := s.Repo.GetFollowRelation(viewerID, userID)
			if err != nil || status != "accepted" {
				return nil, ErrAccessDenied
			}
		}
	}
	return s.Repo.CreatePost(userID, payload)
}

func (s *Service) GetFeed(userID int64, limit int, cursor *int64) ([]models.Post, bool, error) {
	return s.Repo.GetFeedPosts(userID, limit, cursor)
}

func (s *Service) GetPost(userID int64, postID int64) (*models.Post, error) {
	post, err := s.Repo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if err := s.checkPostVisibility(userID, post); err != nil {
		return nil, err
	}
	if post.Privacy == "private" {
		viewers, err := s.Repo.GetPostViewers(post.ID)
		if err != nil {
			return nil, err
		}
		post.Viewers = viewers
	}

	return post, nil
}

func (s *Service) AddComment(userID int64, payload models.CreateCommentPayload) (*models.Comment, error) {
	// Check post existence and permission
	post, err := s.Repo.GetPostByID(payload.PostID)
	if err != nil {
		return nil, err
	}

	if err := s.checkPostVisibility(userID, post); err != nil {
		return nil, err
	}

	return s.Repo.CreateComment(userID, payload)
}

func (s *Service) GetComments(userID int64, postID int64) ([]models.Comment, error) {
	post, err := s.Repo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if err := s.checkPostVisibility(userID, post); err != nil {
		return nil, err
	}

	return s.Repo.GetPostComments(postID)
}

func (s *Service) checkPostVisibility(userID int64, post *models.Post) error {
	if post.UserID == userID {
		return nil
	}
	is_private, err := s.Repo.IsPrivate(post.UserID)
	if err != nil {
		return ErrAccessDenied
	}
	// Group posts: members only, regardless of privacy
	if post.GroupID != nil && *post.GroupID > 0 {
		status, err := s.Repo.GetMemberStatus(*post.GroupID, userID)
		if err == nil && status == "accepted" {
			return nil
		}
		return ErrAccessDenied
	}
	if post.Privacy == "private" {
		isViewer, err := s.Repo.IsPostViewer(post.ID, userID)
		if err != nil || !isViewer {
			return ErrAccessDenied
		}
		return nil
	}
	if post.Privacy == "almost private" || is_private {
		status, err := s.Repo.GetFollowRelation(userID, post.UserID)
		if err == nil && status == "accepted" {
			return nil
		}
		return ErrAccessDenied
	}
	if post.Privacy == "public" {
		return nil
	}
	return ErrAccessDenied
}
