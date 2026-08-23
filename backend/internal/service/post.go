package service

import (
	"context"
	"errors"
	"fmt"

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
	}
	fmt.Println("payload ::::::", payload)
	return s.Repo.CreatePost(ctx, userID, payload)
}

func (s *Service) GetFeed(ctx context.Context, userID int64, limit int, cursor *int64) ([]models.Post, bool, error) {
	return s.Repo.GetFeedPosts(ctx, userID, limit, cursor)
}

func (s *Service) GetPost(ctx context.Context, userID int64, postID int64) (*models.Post, error) {
	post, err := s.Repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if err := s.checkPostVisibility(ctx, userID, post); err != nil {
		return nil, err
	}
	fmt.Println("innnnnnnnn")
	if post.Privacy == "private" {
		viewers, err := s.Repo.GetPostViewers(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		post.Viewers = viewers
	}

	return post, nil
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

	return s.Repo.GetPostComments(ctx, postID)
}

func (s *Service) checkPostVisibility(ctx context.Context, userID int64, post *models.Post) error {
	fmt.Println("prayvacy::::::",post.Privacy)
	fmt.Println("post.UserID,userID",post.UserID,userID)
	if post.UserID == userID {
		return nil
	}
	is_private,err:=s.Repo.IsPrivate(ctx, post.UserID)
	if (err!=nil){
		return ErrAccessDenied
	}
	fmt.Println("is_private::::::",is_private)
	// Group posts are visible only to accepted group members, regardless of privacy
	if post.GroupID != nil && *post.GroupID > 0 {
		status, err := s.Repo.GetMemberStatus(ctx, *post.GroupID, userID)
		if err == nil && status == "accepted" {
			return nil
		}
		return ErrAccessDenied
	}
	if post.Privacy == "almost private" || is_private {
		status, err := s.Repo.GetFollowRelation(ctx, userID, post.UserID)
		if err == nil && status == "accepted" {
			return nil
		}
		return ErrAccessDenied
	}
	if post.Privacy == "private"  {
		isViewer, err := s.Repo.IsPostViewer(ctx, post.ID, userID)
		if err != nil || !isViewer {
			return ErrAccessDenied
		}
		return nil
	}
	if post.Privacy == "public" {
		return nil
	}
	return ErrAccessDenied
}
