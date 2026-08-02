package requests

import (
	"errors"
	"strings"

	"kuu/internal/models"
)

// ValidateUpdateProfile ensures required core fields are present and structurally sound
func ValidateUpdateProfile(payload models.UpdateProfilePayload) []error {
	var errs []error

	if strings.TrimSpace(payload.FirstName) == "" {
		errs = append(errs, errors.New("first_name cannot be blank"))
	}
	if strings.TrimSpace(payload.LastName) == "" {
		errs = append(errs, errors.New("last_name cannot be blank"))
	}

	gender := strings.ToLower(strings.TrimSpace(payload.Gender))
	if gender != "male" && gender != "female" {
		errs = append(errs, errors.New("gender must be either 'male', 'female'"))
	}

	if strings.TrimSpace(payload.DateOfBirth) == "" {
		errs = append(errs, errors.New("date_of_birth is required"))
	}

	if payload.IsPublic != 0 && payload.IsPublic != 1 {
		errs = append(errs, errors.New("is_public must be either 0 (private) or 1 (public)"))
	}

	return errs
}

func ValidateFollowAction(payload models.FollowActionPayload) []error {
	var errs []error
	if payload.TargetUserID <= 0 {
		errs = append(errs, errors.New("target_user_id is required and must be valid"))
	}
	return errs
}
