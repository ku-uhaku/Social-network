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
	if p.Privacy != "public" && p.Privacy != "almost private" && p.Privacy != "private" {
		errs = append(errs, errors.New("privacy must be 'public', 'almost private', or 'private'"))
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
