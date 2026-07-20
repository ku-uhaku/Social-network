package handler

import (
	"encoding/json"
	"net/http"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/requests"
)

// UpdateProfile reads active session context data and applies new settings
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate user context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	// 2. Decode incoming update values
	var payload models.UpdateProfilePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	// 3. Fire structural assertions
	if validationErrs := requests.ValidateUpdateProfile(payload); len(validationErrs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, validationErrs)
		return
	}
	// 4. Update profiles inside persistence layer
	updatedUser, err := h.Service.UpdateProfile(r.Context(), user.ID, payload)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 5. Return success data
	helper.Success(w, http.StatusOK, "Profile updated successfully", updatedUser)
}
