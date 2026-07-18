package requests

import (
	"errors"
	"strings"

	"kuu/internal/models"
)

// ValidateCreateGroup runs layout filters on structural metadata inputs
func ValidateCreateGroup(payload models.CreateGroupPayload) []error {
	var errs []error

	trimmedTitle := strings.TrimSpace(payload.Title)
	if trimmedTitle == "" {
		errs = append(errs, errors.New("group title is required"))
	} else if len(trimmedTitle) < 3 || len(trimmedTitle) > 50 {
		errs = append(errs, errors.New("group title must be between 3 and 50 characters long"))
	}

	if payload.IsPublic == nil {
		errs = append(errs, errors.New("is_public visibility state is required"))
	} else if *payload.IsPublic != 0 && *payload.IsPublic != 1 {
		errs = append(errs, errors.New("is_public must be 1 (true) or 0 (false)"))
	}

	return errs
}
