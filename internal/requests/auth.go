package requests

import (
	"regexp"
	"strings"

	"kuu/internal/models"
)

func ValidateRegister(payload models.InputRegisterPayload) []ValidationError {
	var errs []ValidationError

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

	if strings.TrimSpace(payload.Username) == "" {
		errs = append(errs, ValidationError{
			Field:   "username",
			Message: "username is required",
		})
	} else if !usernameRegex.MatchString(payload.Username) {
		errs = append(errs, ValidationError{
			Field:   "username",
			Message: "username must be between 3 and 20 alphanumeric characters or underscores",
		})
	}

	if strings.TrimSpace(payload.Email) == "" {
		errs = append(errs, ValidationError{
			Field:   "email",
			Message: "email is required",
		})
	} else if !emailRegex.MatchString(payload.Email) {
		errs = append(errs, ValidationError{
			Field:   "email",
			Message: "invalid email address format",
		})
	}

	if strings.TrimSpace(payload.FirstName) == "" {
		errs = append(errs, ValidationError{
			Field:   "first_name",
			Message: "first name is required",
		})
	}

	if strings.TrimSpace(payload.LastName) == "" {
		errs = append(errs, ValidationError{
			Field:   "last_name",
			Message: "last name is required",
		})
	}

	if len(payload.Password) < 8 {
		errs = append(errs, ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters long",
		})
	}

	return errs
}

func ValidateLogin(payload models.InputLoginPayload) []ValidationError {
	var errs []ValidationError

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

	login := strings.TrimSpace(payload.Login)

	if login == "" {
		errs = append(errs, ValidationError{
			Field:   "login",
			Message: "login is required",
		})
	} else if !emailRegex.MatchString(login) && !usernameRegex.MatchString(login) {
		errs = append(errs, ValidationError{
			Field:   "login",
			Message: "must be a valid username or email",
		})
	}

	if payload.Password == "" {
		errs = append(errs, ValidationError{
			Field:   "password",
			Message: "password is required",
		})
	} else if len(payload.Password) < 8 {
		errs = append(errs, ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters long",
		})
	}

	return errs
}
