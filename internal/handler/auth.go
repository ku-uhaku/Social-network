package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"kuu/internal/helper"
	"kuu/internal/models"
	"kuu/internal/requests"

	"github.com/google/uuid"
)

// Login handles verifying credentials and starting a session
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.InputLoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	r.Body.Close()

	validationErrs := requests.ValidateLogin(payload)
	if len(validationErrs) > 0 {
		errMessages := make([]string, len(validationErrs))
		for i, err := range validationErrs {
			errMessages[i] = err.Error()
		}
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", errMessages)
		return
	}

	user, err := h.Repo.AuthficationUser(payload)
	if err != nil {
		helper.Error(w, http.StatusUnauthorized, "Invalid identifier or password")
		return
	}

	token := uuid.NewString()
	sessionDuration := 24 * time.Hour

	sessionInfo, err := h.Repo.CreateSession(user.ID, token, sessionDuration)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, "Failed to create secure user session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "kuu_token",
		Value:    sessionInfo.ID,
		Expires:  sessionInfo.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	helper.Success(w, http.StatusOK, "Login successful", user)
}

// Logout deletes the session from the database and clears the browser cookie
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		helper.Error(w, http.StatusBadRequest, "No active session found")
		return
	}

	err = h.Repo.DeleteSession(cookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.Error(w, http.StatusNotFound, "Session already invalid or expired")
			return
		}
		helper.Error(w, http.StatusInternalServerError, "Failed to terminate session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "kuu_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	helper.Success(w, http.StatusOK, "Successfully logged out", nil)
}

// CheckAuth can be requested by frontend frameworks (React/Vue) on startup to see if the user is logged in
func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("kuu_token")
	if err != nil {
		helper.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	session, err := h.Repo.GetSession(cookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.Error(w, http.StatusUnauthorized, "Session has expired or is invalid")
			return
		}
		helper.Error(w, http.StatusInternalServerError, "Database verification failure")
		return
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		helper.Error(w, http.StatusUnauthorized, "Session expired")
		return
	}

	helper.Success(w, http.StatusOK, "Session is authenticated", session)
}
