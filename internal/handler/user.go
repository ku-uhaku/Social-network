package handler

import (
	"net/http"

	"kuu/internal/helper"
	"kuu/internal/middleware"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Pull the user out of the context cleanly using the helper function
	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusInternalServerError, "Failed to resolve identity context")
		return
	}

	helper.Success(w, http.StatusOK, "Profile retrieved", currentUser)
}
