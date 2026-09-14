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
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to secure user credentials")
	}
	user, err := s.Repo.CreateUser(payload, string(hashedBytes))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("username or email is already taken")
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(payload models.InputLoginPayload) (*models.Session, *models.User, error) {
	user, err := s.Repo.AuthenticationUser(payload)
	if err != nil {
		return nil, nil, errors.New("invalid identifier or password")
	}

	token := uuid.NewString()
	duration := 10 * 365 * 24 * time.Hour // persistent until explicit logout

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

func (s *Service) DeleteSession(token string) error {
	return s.Repo.DeleteSession(token)
}
