package requests

import (
	"errors"
	"regexp"
	"strings"

	"kuu/internal/models"
)

func ValidateRegister(payload models.InputRegisterPayload) []error {
	var errs []error
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

	if strings.TrimSpace(payload.Username) == "" {
		errs = append(errs, errors.New("username is required"))
	} else if !usernameRegex.MatchString(payload.Username) {
		errs = append(errs, errors.New("username must be between 3 and 20 alphanumeric characters or underscores"))
	}

	if strings.TrimSpace(payload.Email) == "" {
		errs = append(errs, errors.New("email is required"))
	} else if !emailRegex.MatchString(payload.Email) {
		errs = append(errs, errors.New("invalid email address format"))
	}

	if strings.TrimSpace(payload.FirstName) == "" {
		errs = append(errs, errors.New("first name is required"))
	}
	if strings.TrimSpace(payload.LastName) == "" {
		errs = append(errs, errors.New("last name is required"))
	}

	if len(payload.Password) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters long"))
	}

	return errs
}

func ValidateLogin(payload models.InputLoginPayload) []error {
	var errs []error
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

	trimmedIdentifier := strings.TrimSpace(payload.Login)

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
