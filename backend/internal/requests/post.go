package requests

import (
	"errors"
	"strings"

	"kuu/internal/models"
)

func ValidateCreatePost(p models.CreatePostPayload) []error {
	var errs []error
	if strings.TrimSpace(p.Title) == "" {
		errs = append(errs, errors.New("title cannot be empty"))
	}
	if strings.TrimSpace(p.Content) == "" {
		errs = append(errs, errors.New("content cannot be empty"))
	}

	p.Privacy = strings.ToLower(strings.TrimSpace(p.Privacy))
	if p.GroupID != nil && *p.GroupID > 0 {
		// If group_id is present, enforce group scope
		if p.Privacy != "group" {
			errs = append(errs, errors.New("privacy must be 'group' when posting to a group"))
		}
	} else {
		if p.Privacy != "public" && p.Privacy != "followers" {
			errs = append(errs, errors.New("privacy must be 'public' or 'followers' for user posts"))
		}
	}

	return errs
}

func ValidateCreateComment(p models.CreateCommentPayload) []error {
	var errs []error
	if p.PostID <= 0 {
		errs = append(errs, errors.New("post_id is required"))
	}
	if strings.TrimSpace(p.Content) == "" {
		errs = append(errs, errors.New("comment content cannot be empty"))
	}
	return errs
}

func ValidateReaction(p models.ReactionPayload) []error {
	var errs []error
	if (p.PostID == nil || *p.PostID <= 0) && (p.CommentID == nil || *p.CommentID <= 0) {
		errs = append(errs, errors.New("must specify either post_id or comment_id"))
	}
	if p.Type != 1 && p.Type != -1 && p.Type != 0 {
		errs = append(errs, errors.New("reaction type must be 1 (like), -1 (dislike), or 0 (remove)"))
	}
	return errs
}
