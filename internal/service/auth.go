package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"kuu/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(ctx context.Context, payload models.InputRegisterPayload) (*models.User, error) {
	// 1. Hash the incoming plaintext password safely
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to secure user credentials")
	}

	// 2. Pass the record onto the persistence layer
	user, err := s.Repo.CreateUser(ctx, payload, string(hashedBytes))
	if err != nil {
		// Detect SQLite unique constraint failures for emails/usernames
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("username or email is already taken")
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, payload models.InputLoginPayload) (*models.Session, *models.User, error) {
	// 1. Fetch user records and verify bcrypt password inside the repo
	user, err := s.Repo.AuthenticationUser(ctx, payload)
	if err != nil {
		return nil, nil, errors.New("invalid identifier or password")
	}

	// 2. Provision secure UUID strings for session tracking
	token := uuid.NewString()
	duration := 24 * time.Hour

	// 3. Commit session record into database
	sessionInfo, err := s.Repo.CreateSession(ctx, user.ID, token, duration)
	if err != nil {
		return nil, nil, errors.New("failed to establish session")
	}

	return sessionInfo, user, nil
}
