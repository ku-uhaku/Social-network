package service

import (
	"errors"
	"strings"
	"time"

	"kuu/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(payload models.InputRegisterPayload) (*models.User, error) {
	// 1. Hash the incoming plaintext password safely
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to secure user credentials")
	}
	// 2. Pass the record onto the persistence layer
	user, err := s.Repo.CreateUser(payload, string(hashedBytes))
	if err != nil {
		// Detect SQLite unique constraint failures for emails/usernames
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("username or email is already taken")
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(payload models.InputLoginPayload) (*models.Session, *models.User, error) {
	// 1. Fetch user records and verify bcrypt password inside the repo
	user, err := s.Repo.AuthenticationUser(payload)
	if err != nil {
		return nil, nil, errors.New("invalid identifier or password")
	}

	// 2. Provision secure UUID strings for session tracking
	token := uuid.NewString()
	duration := 10 * 365 * 24 * time.Hour // persistent until explicit logout

	// 3. Commit session record into database
	sessionInfo, err := s.Repo.CreateSession(user.ID, token, duration)
	if err != nil {
		return nil, nil, errors.New("failed to establish session")
	}

	return sessionInfo, user, nil
}

func (s *Service) ValidateSession(token string) (*models.User, error) {
	user, err := s.Repo.GetUserBySessionToken(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Deletes a session
func (s *Service) DeleteSession(token string) error {
	return s.Repo.DeleteSession(token)
}
