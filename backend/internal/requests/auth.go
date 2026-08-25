package requests

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"kuu/internal/helper"
	"kuu/internal/models"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

func ParseRegisterPayload(r *http.Request) (models.InputRegisterPayload, error) {
	var payload models.InputRegisterPayload

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return payload, fmt.Errorf("malformed multipart form data: %w", err)
		}

		payload.Username = strings.TrimSpace(r.FormValue("username"))
		payload.Email = strings.TrimSpace(r.FormValue("email"))
		payload.Password = r.FormValue("password")
		payload.FirstName = strings.TrimSpace(r.FormValue("first_name"))
		payload.LastName = strings.TrimSpace(r.FormValue("last_name"))
		payload.Gender = strings.TrimSpace(r.FormValue("gender"))
		payload.DateOfBirth = strings.TrimSpace(r.FormValue("date_of_birth"))

		if aboutMe := strings.TrimSpace(r.FormValue("about_me")); aboutMe != "" {
			payload.AboutMe = &aboutMe
		}
		if file, header, err := r.FormFile("avatar"); err == nil {
			// _,err=helper.IsValidImage([]byte(*payload.Avatar))
			// if err!=nil{
			// 	return payload, fmt.Errorf("the avatar not good: %w", err)
			// }
			defer file.Close()
			if header != nil && header.Size > 0 {
				avatarName, err := helper.SaveUploadedImage(file, header)
				if err != nil {
					return payload, fmt.Errorf("failed to save avatar file: %w", err)
				}

				payload.Avatar = &avatarName
			}
		}
	}
	return payload, nil
}

func ValidateRegister(payload models.InputRegisterPayload) []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(payload.Username) != "" && !usernameRegex.MatchString(payload.Username) {
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
	} else if _, err := mail.ParseAddress(payload.Email); err != nil {
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

	gender := strings.TrimSpace(payload.Gender)
	if gender == "" {
		errs = append(errs, ValidationError{
			Field:   "gender",
			Message: "gender is required",
		})
	} else if gender != "male" && gender != "female" {
		errs = append(errs, ValidationError{
			Field:   "gender",
			Message: "gender must be either male or female",
		})
	}
	// i would validate usee age 
	dateOfBirth := strings.TrimSpace(payload.DateOfBirth)
	if dateOfBirth == "" {
		errs = append(errs, ValidationError{
			Field:   "date_of_birth",
			Message: "date of birth is required",
		})
	} else if birth, parseErr := time.Parse("2006-01-02", dateOfBirth); parseErr != nil {
		errs = append(errs, ValidationError{
			Field:   "date_of_birth",
			Message: "date of birth must be a valid date (YYYY-MM-DD)",
		})
	} else if birth.AddDate(16, 0, 0).After(time.Now()) {
		errs = append(errs, ValidationError{
			Field:   "date_of_birth",
			Message: "you must be at least 16 years old to register",
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

	login := strings.TrimSpace(payload.Login)
	_, emailErr := mail.ParseAddress(login)
	isValidEmail := emailErr == nil

	if login == "" {
		errs = append(errs, ValidationError{
			Field:   "login",
			Message: "login is required",
		})
	} else if !isValidEmail && !usernameRegex.MatchString(login) {
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
