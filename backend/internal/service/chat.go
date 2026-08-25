package service

import (
	"database/sql"
	"errors"
	"strings"

	"kuu/internal/models"
)

const maxMessageLength = 2000

var (
	ErrChatSelf         = errors.New("schizophrenia detected")
	ErrChatNoConnection = errors.New("can only chat with users you follow or who follow you")
	ErrChatEmpty        = errors.New("empty message")
	ErrChatTooLong      = errors.New("message too long")
)

func (s *Service) SaveDirectMessage(senderID, receiverID int64, content string) (*models.DirectMessage, error) {
	if senderID == receiverID {
		return nil, ErrChatSelf
	}

	if _, err := s.Repo.GetUserByID(receiverID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("receiver does not exist")
		}
		return nil, err
	}

	connected, err := s.Repo.CanChat(senderID, receiverID)
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
	if len(content) > maxMessageLength {
		return nil, ErrChatTooLong
	}

	return s.Repo.SaveDirectMessage(senderID, receiverID, content)
}

func (s *Service) SaveGroupMessage(senderID, groupID int64, content string) (*models.GroupMessage, error) {
	status, err := s.Repo.GetMemberStatus(groupID, senderID)
	if err != nil || status != "accepted" {
		return nil, errors.New("unauthorized group access")
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrChatEmpty
	}
	if len(content) > maxMessageLength {
		return nil, ErrChatTooLong
	}

	return s.Repo.SaveGroupMessage(senderID, groupID, content)
}

// Authorizes chat, returns a page, marks read on latest page
func (s *Service) GetDirectHistory(userA, userB int64, limit, offset int) ([]models.DirectMessage, error) {
	if userA == userB {
		return nil, ErrChatSelf
	}

	connected, err := s.Repo.CanChat(userA, userB)
	if err != nil {
		return nil, err
	}
	if !connected {
		return nil, ErrChatNoConnection
	}

	if offset == 0 {
		if err := s.Repo.MarkChatRead(userA, userB); err != nil {
			return nil, err
		}
	}

	return s.Repo.GetDirectHistory(userA, userB, limit, offset)
}

// Marks a conversation as read
func (s *Service) MarkChatRead(userA, userB int64) error {
	if userA == userB {
		return ErrChatSelf
	}
	return s.Repo.MarkChatRead(userA, userB)
}

// Chat-able users with their latest DM
func (s *Service) GetConversations(userID int64) ([]models.ConversationMetadata, error) {
	return s.Repo.ListConversations(userID)
}

func (s *Service) GetGroupHistory(userID, groupID int64, limit, offset int) ([]models.GroupMessage, error) {
	status, err := s.Repo.GetMemberStatus(groupID, userID)
	if err != nil || status != "accepted" {
		return nil, errors.New("unauthorized group access")
	}
	return s.Repo.GetGroupHistory(groupID, limit, offset)
}

func (s *Service) GetGroupMemberIDs(groupID int64) ([]int64, error) {
	return s.Repo.GetGroupMemberIDs(groupID)
}
