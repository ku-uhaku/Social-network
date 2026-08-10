package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"kuu/internal/helper"
	"kuu/internal/models"
)

// ParseCreatePostPayload parses the request body into a CreatePostPayload struct
func ParseCreatePostPayload(r *http.Request) (models.CreatePostPayload, error) {
	var payload models.CreatePostPayload

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return payload, fmt.Errorf("malformed multipart form data: %w", err)
		}

		payload.Title = strings.TrimSpace(r.FormValue("title"))
		payload.Content = strings.TrimSpace(r.FormValue("content"))
		payload.Privacy = strings.TrimSpace(r.FormValue("privacy"))

		// Parse group_id when posting inside a group
		if groupIDStr := r.FormValue("group_id"); groupIDStr != "" {
			groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
			if err != nil {
				return payload, fmt.Errorf("invalid group_id")
			}
			payload.GroupID = &groupID
		}

		// Parse visible_to from multipart form
		if visibleToValues := r.Form["visible_to"]; len(visibleToValues) > 0 {
			visibleTo := make([]int64, 0, len(visibleToValues))
			for _, v := range visibleToValues {
				id, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return payload, fmt.Errorf("invalid user ID in visible_to: %w", err)
				}
				visibleTo = append(visibleTo, id)
			}
			payload.VisibleTo = visibleTo
		}

		if file, header, err := r.FormFile("image"); err == nil {
			defer file.Close()
			if header != nil && header.Size > 0 {
				imageName, err := helper.SaveUploadedImage(file, header)
				if err != nil {
					return payload, fmt.Errorf("failed to save post image: %w", err)
				}
				payload.ImageURL = &imageName
			}
		}
	} else if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, fmt.Errorf("invalid JSON payload")
	}

	return payload, nil
}

// ValidateCreatePost validates the CreatePostPayload struct
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

	// visible_to for private posts
	if p.Privacy == "private" && len(p.VisibleTo) == 0 {
		errs = append(errs, errors.New("private posts must specify visible_to users"))
	}

	return errs
}

// ParseCreateCommentPayload parses the request body into a CreateCommentPayload struct
func ParseCreateCommentPayload(r *http.Request) (models.CreateCommentPayload, error) {
	var payload models.CreateCommentPayload

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return payload, fmt.Errorf("malformed multipart form data: %w", err)
		}

		postID, err := strconv.ParseInt(r.FormValue("post_id"), 10, 64)
		if err != nil {
			return payload, fmt.Errorf("invalid post_id")
		}
		payload.PostID = postID
		payload.Title = strings.TrimSpace(r.FormValue("title"))
		payload.Content = strings.TrimSpace(r.FormValue("content"))

		if file, header, err := r.FormFile("image"); err == nil {
			defer file.Close()
			if header != nil && header.Size > 0 {
				imageName, err := helper.SaveUploadedImage(file, header)
				if err != nil {
					return payload, fmt.Errorf("failed to save comment image: %w", err)
				}
				payload.ImageURL = &imageName
			}
		}
	} else if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, fmt.Errorf("invalid JSON payload")
	}

	return payload, nil
}

// ValidateCreateComment validates the CreateCommentPayload struct
func ValidateCreateComment(p models.CreateCommentPayload) []error {
	var errs []error
	if p.PostID <= 0 {
		errs = append(errs, errors.New("post_id is required"))
	}
	if strings.TrimSpace(p.Title) == "" {
		errs = append(errs, errors.New("comment title is required"))
	}
	if strings.TrimSpace(p.Content) == "" {
		errs = append(errs, errors.New("comment content cannot be empty"))
	}

	return errs
}