package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"kuu/internal/models"
)

const maxDirectMessageLength = 2000

var (
	ErrChatSelf         = errors.New("schizophrenia detected")
	ErrChatNoConnection = errors.New("can only chat with users you follow or who follow you")
	ErrChatEmpty        = errors.New("empty message")
	ErrChatTooLong      = errors.New("message too long")
)

func (s *Service) SaveDirectMessage(ctx context.Context, senderID, receiverID int64, content string) (*models.DirectMessage, error) {
	if senderID == receiverID {
		return nil, ErrChatSelf
	}

	// Verify the receiver actually exists
	if _, err := s.Repo.GetUserByID(ctx, receiverID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("receiver does not exist")
		}
		return nil, err
	}

	connected, err := s.Repo.CanChat(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	if !connected {
		return nil, ErrChatNoConnection
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrChatEmpty
	}
	if len(content) > maxDirectMessageLength {
		return nil, ErrChatTooLong
	}

	return s.Repo.SaveDirectMessage(ctx, senderID, receiverID, content)
}

func (s *Service) SaveGroupMessage(ctx context.Context, senderID, groupID int64, content string) (*models.GroupMessage, error) {
	// Verify user is a group member
	status, err := s.Repo.GetMemberStatus(ctx, groupID, senderID)
	if err != nil || status != "accepted" {
		return nil, errors.New("unauthorized group access")
	}
	return s.Repo.SaveGroupMessage(ctx, senderID, groupID, content)
}

// GetDirectHistory authorizes the conversation before returning its full
// history: users may only view chats with someone they follow or who follows them.
func (s *Service) GetDirectHistory(ctx context.Context, userA, userB int64) ([]models.DirectMessage, error) {
	if userA == userB {
		return nil, ErrChatSelf
	}

	connected, err := s.Repo.CanChat(ctx, userA, userB)
	if err != nil {
		return nil, err
	}
	if !connected {
		return nil, ErrChatNoConnection
	}

	return s.Repo.GetDirectHistory(ctx, userA, userB)
}

// GetConversations lists all chat-able users for the viewer with their latest DM.
func (s *Service) GetConversations(ctx context.Context, userID int64) ([]models.ConversationMetadata, error) {
	return s.Repo.ListConversations(ctx, userID)
}

func (s *Service) GetGroupHistory(ctx context.Context, userID, groupID int64, limit, offset int) ([]models.GroupMessage, error) {
	status, err := s.Repo.GetMemberStatus(ctx, groupID, userID)
	if err != nil || status != "accepted" {
		return nil, errors.New("unauthorized group access")
	}
	return s.Repo.GetGroupHistory(ctx, groupID, limit, offset)
}

func (s *Service) GetGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return s.Repo.GetGroupMemberIDs(ctx, groupID)
}
