package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/service"
)

// Messages per history page
const pageSize = 30

// GetDirectHistory GET /api/v1/chat/direct?user_id=123&page=1
func (h *Handler) GetDirectHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	targetUserID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil || targetUserID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid user_id parameter")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	messages, err := h.Service.GetDirectHistory(user.ID, targetUserID, pageSize, (page-1)*pageSize)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrChatSelf) || errors.Is(err, service.ErrChatNoConnection) {
			status = http.StatusForbidden
		}
		helper.Error(w, status, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Direct message history retrieved", messages)
}

// MarkChatRead POST /api/v1/chat/read { user_id }
func (h *Handler) MarkChatRead(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	if err := h.Service.MarkChatRead(user.ID, req.UserID); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Conversation marked as read", nil)
}

// GetConversations GET /api/v1/chat/conversations
func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversations, err := h.Service.GetConversations(user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Conversations retrieved", conversations)
}

// GetGroupHistory GET /api/v1/chat/group?group_id=456&page=1
func (h *Handler) GetGroupHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groupID, err := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid group_id parameter")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	messages, err := h.Service.GetGroupHistory(user.ID, groupID, pageSize, (page-1)*pageSize)
	if err != nil {
		helper.Error(w, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group message history retrieved", messages)
}
