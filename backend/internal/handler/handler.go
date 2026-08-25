package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/service"
	"kuu/internal/websocket"
)

type Handler struct {
	Service *service.Service
	Hub     *websocket.Hub
}

func New(svc *service.Service, hub *websocket.Hub) *Handler {
	return &Handler{
		Service: svc,
		Hub:     hub,
	}
}

// Helper functions

func (h *Handler) authUser(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, false
	}
	return user, true
}

func queryPositiveInt64(w http.ResponseWriter, r *http.Request, key string, message string) (int64, bool) {
	value, err := helper.GetParamInt64(r, key)
	if err != nil || value <= 0 {
		helper.Error(w, http.StatusBadRequest, message)
		return 0, false
	}
	return value, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return false
	}
	defer r.Body.Close()
	return true
}

// Writes 422 if errs is non-empty; returns whether it did.
func writeValidationErrors(w http.ResponseWriter, errs []error) bool {
	if len(errs) == 0 {
		return false
	}
	helper.ValidationErrorResponse(w, http.StatusUnprocessableEntity, errs)
	return true
}

func parseFeedParams(w http.ResponseWriter, r *http.Request) (int, *int64, bool) {
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil || parsed < 1 {
			helper.Error(w, http.StatusBadRequest, "Invalid limit parameter")
			return 0, nil, false
		}
		limit = parsed
		if limit > 50 {
			limit = 50
		}
	}

	var cursor *int64
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		parsed, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || parsed < 1 {
			helper.Error(w, http.StatusBadRequest, "Invalid cursor parameter")
			return 0, nil, false
		}
		cursor = &parsed
	}

	return limit, cursor, true
}
