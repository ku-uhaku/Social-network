package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/requests"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	payload, err := requests.ParseRegisterPayload(r)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if validationErrs := requests.ValidateRegister(payload); len(validationErrs) > 0 {
		helper.ValidationErrorResponse(w, http.StatusUnprocessableEntity, validationErrs)
		return
	}

	user, err := h.Service.RegisterUser(payload)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(w, http.StatusCreated, "Account created successfully", user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.InputLoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	if validationErrs := requests.ValidateLogin(payload); len(validationErrs) > 0 {
		helper.ValidationErrorResponse(w, http.StatusUnprocessableEntity, validationErrs)
		return
	}

	// 1. Service returns BOTH sessionInfo and user
	sessionInfo, user, err := h.Service.LoginUser(payload)
	if err != nil {
		helper.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	const sessionCookieDuration = 10 * 365 * 24 * time.Hour // never expires

	// 2. Set the secure HTTP-Only cookie so the middleware can authenticate upcoming requests
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionInfo.ID,
		Expires:  time.Now().Add(sessionCookieDuration),
		MaxAge:   int(sessionCookieDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// 3. Respond with OK status and the user profile schema payload
	helper.Success(w, http.StatusOK, "Login successful", user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete session from the database asynchronously or synchronously
		_ = h.Service.DeleteSession(cookie.Value)
	}

	// Explicitly instruct the browser to wipe the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Tells browser to delete immediately
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	helper.Success(w, http.StatusOK, "Logged out successfully", nil)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusInternalServerError, "Failed to retrieve profile context")
		return
	}
	helper.Success(w, http.StatusOK, "Profile retrieved successfully", user)
}
