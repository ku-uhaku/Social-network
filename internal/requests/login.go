package requests

import (
	"errors"
	"regexp"
	"strings"

	"kuu/internal/models"
)

func ValidateLogin(payload models.InputLoginPayload) []error {
	var errs []error
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

	trimmedIdentifier := strings.TrimSpace(payload.Identifier)

	if trimmedIdentifier == "" {
		errs = append(errs, errors.New("identifier is required"))
	} else {
		isEmail := emailRegex.MatchString(trimmedIdentifier)
		isUsername := usernameRegex.MatchString(trimmedIdentifier)

		if !isEmail && !isUsername {
			errs = append(errs, errors.New("invalid identifier format (must be a valid username or email)"))
		}
	}

	if payload.Password == "" {
		errs = append(errs, errors.New("password is required"))
	} else if len(payload.Password) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters long"))
	}

	return errs
}
