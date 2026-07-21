package handler

import (
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
)

// GetDirectHistory GET /api/v1/chat/direct?user_id=123&page=1
func (h *Handler) GetDirectHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	targetUserIDStr := r.URL.Query().Get("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil || targetUserID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid user_id parameter")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit := 30
	offset := (page - 1) * limit

	messages, err := h.Service.GetDirectHistory(r.Context(), user.ID, targetUserID, limit, offset)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Direct message history retrieved", messages)
}

// GetGroupHistory GET /api/v1/chat/group?group_id=456&page=1
func (h *Handler) GetGroupHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid group_id parameter")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit := 30
	offset := (page - 1) * limit

	messages, err := h.Service.GetGroupHistory(r.Context(), user.ID, groupID, limit, offset)
	if err != nil {
		helper.Error(w, http.StatusForbidden, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group message history retrieved", messages)
}
