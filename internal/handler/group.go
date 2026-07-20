package handler

import (
	"encoding/json"
	"net/http"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/requests"
)

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate that we have a valid context user
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	var payload models.CreateGroupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	// 2. Run validations
	if validationErrs := requests.ValidateCreateGroup(payload); len(validationErrs) > 0 {
		helper.ValidationErrorResponse(w, http.StatusUnprocessableEntity, validationErrs)
		return
	}

	// 3. Fallback default to public (1) if field wasn't explicitly provided
	// 4. Fire service processing
	group, err := h.Service.CreateGroup(r.Context(), user.ID, payload)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusCreated, "Group created successfully", group)
}

func (h *Handler) GetPublicGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.Service.GetPublicGroups(r.Context())
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, "Failed to retrieve groups")
		return
	}

	helper.Success(w, http.StatusOK, "Public groups retrieved successfully", groups)
}

func (h *Handler) InviteMembers(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate user context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	// 2. Extract group ID from URL path (adjust helper according to your router framework)
	groupID, err := helper.GetParamInt64(r, "id")
	if err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid group ID format")
		return
	}

	// 3. Decode payload containing targets
	var payload models.InviteMembersPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	// 4. Fire service processing
	err = h.Service.InviteMembersToGroup(r.Context(), user.ID, groupID, payload)
	if err != nil {
		// You can check for specific errors here (e.g., unauthorized to invite, group not found)
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Invitations sent successfully", nil)
}

func (h *Handler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate user context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	// 2. Decode the incoming group target payload
	var payload models.RespondToInvitePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	if payload.GroupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid group ID parameter")
		return
	}

	// 3. Execute change operation
	err := h.Service.AcceptInvitation(r.Context(), user.ID, payload.GroupID)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Successfully joined the group!", nil)
}
