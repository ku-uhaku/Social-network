package requests

import (
	"strings"

	"kuu/internal/models"
)

func ValidateCreateGroup(payload models.CreateGroupPayload) []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(payload.Title) == "" {
		errs = append(errs, ValidationError{
			Field:   "title",
			Message: "group title is required",
		})
	} else if len(payload.Title) > 100 {
		errs = append(errs, ValidationError{
			Field:   "title",
			Message: "group title cannot exceed 100 characters",
		})
	}

	if strings.TrimSpace(payload.Description) == "" {
		errs = append(errs, ValidationError{
			Field:   "description",
			Message: "group description is required",
		})
	}

	return errs
}
