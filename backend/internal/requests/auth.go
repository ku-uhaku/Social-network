package requests

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
			defer file.Close()
			if header != nil && header.Size > 0 {
				fileBytes, err := io.ReadAll(file)
				if err != nil {
					return payload, fmt.Errorf("failed to read avatar file: %w", err)
				}

				hash := sha256.Sum256(fileBytes)
				avatarName := hex.EncodeToString(hash[:])
				mediaDir := filepath.Join("backend", "media")
				if err := os.MkdirAll(mediaDir, 0o755); err != nil {
					return payload, fmt.Errorf("failed to prepare avatar directory: %w", err)
				}

				avatarPath := filepath.Join(mediaDir, avatarName)
				if err := os.WriteFile(avatarPath, fileBytes, 0o644); err != nil {
					return payload, fmt.Errorf("failed to save avatar file: %w", err)
				}

				payload.Avatar = &avatarName
			}
		}

		return payload, nil
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, fmt.Errorf("malformed json request body: %w", err)
	}
	defer r.Body.Close()

	return payload, nil
}

func ValidateRegister(payload models.InputRegisterPayload) []ValidationError {
	var errs []ValidationError

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

	if strings.TrimSpace(payload.Gender) == "" {
		errs = append(errs, ValidationError{
			Field:   "gender",
			Message: "gender is required",
		})
	}

	if strings.TrimSpace(payload.DateOfBirth) == "" {
		errs = append(errs, ValidationError{
			Field:   "date_of_birth",
			Message: "date of birth is required",
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
