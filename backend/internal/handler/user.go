package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/requests"
)

// UpdateProfile reads active session context data and applies new settings
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	var payload models.UpdateProfilePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	if validationErrs := requests.ValidateUpdateProfile(payload); len(validationErrs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, validationErrs)
		return
	}

	if payload.IsPublic == 1 && user.IsPublic == 0 {
		if _, err := h.Service.AcceptAllPendingFollows(r.Context(), user.ID); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := h.Service.ExpireNotificationsByType(r.Context(), user.ID, models.NotificationFollowRequest); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	updatedUser, err := h.Service.UpdateProfile(r.Context(), user.ID, payload)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Profile updated successfully", updatedUser)
}

// GetUserProfile GET /api/v1/user/profile?username=john
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		helper.Error(w, http.StatusBadRequest, "Invalid username parameter")
		return
	}

	profile, err := h.Service.GetUserProfile(r.Context(), user.ID, username)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Profile retrieved successfully", profile)
}

// GetUserPosts GET /api/v1/user/posts?username=john
func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		helper.Error(w, http.StatusBadRequest, "Invalid username parameter")
		return
	}

	targetUser, err := h.Service.GetUserProfile(r.Context(), user.ID, username)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Only the owner, public profiles, or accepted followers can view posts
	if targetUser.FollowStatus != "self" && targetUser.IsPublic == 0 && targetUser.FollowStatus != "accepted" {
		helper.Error(w, http.StatusForbidden, "This account is private")
		return
	}

	posts, err := h.Service.GetUserPosts(r.Context(), targetUser.ID, user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Posts retrieved successfully", posts)
}

// FollowUser POST /api/v1/user/follow
func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.FollowActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if errs := requests.ValidateFollowAction(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	status, err := h.Service.FollowUser(r.Context(), user.ID, payload.TargetUserID)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	msg := "Successfully followed user"
	if status == "pending" {
		msg = "Follow request sent successfully"
		// Avoid duplicate follow_request notifications from spam/retries.
		existing, _ := h.Service.GetNotificationByActorType(
			r.Context(), payload.TargetUserID, user.ID, models.NotificationFollowRequest,
		)
		if existing == nil || existing.IsExpired == 1 {
			actorID := user.ID
			_, err := h.DispatchNotification(r.Context(), payload.TargetUserID, &models.Notification{
				RecipientID: payload.TargetUserID,
				ActorID:     &actorID,
				Type:        models.NotificationFollowRequest,
				Title:       "New follow request",
				Message:     user.Username + " wants to follow you",
				Actions: models.JSONText{
					"buttons": []interface{}{
						map[string]interface{}{"action": "accept", "label": "Accept"},
						map[string]interface{}{"action": "decline", "label": "Decline"},
					},
				},
			})
			if err != nil {
				helper.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	helper.Success(w, http.StatusOK, msg, map[string]string{"status": status})
}

// UnfollowUser POST /api/v1/user/unfollow
func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.FollowActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.Service.UnfollowUser(r.Context(), user.ID, payload.TargetUserID); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Unfollowed user successfully", nil)
}

// AcceptFollowRequest POST /api/v1/user/follow/accept
func (h *Handler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	h.respondToFollowRequest(w, r, true)
}

// DeclineFollowRequest POST /api/v1/user/follow/decline
func (h *Handler) DeclineFollowRequest(w http.ResponseWriter, r *http.Request) {
	h.respondToFollowRequest(w, r, false)
}

func (h *Handler) respondToFollowRequest(w http.ResponseWriter, r *http.Request, accept bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.FollowActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.Service.HandleFollowRequest(r.Context(), user.ID, payload.TargetUserID, accept); err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Expire all follow_request notifications for this requester
	if err := h.Service.ExpireNotificationsByType(r.Context(), user.ID, models.NotificationFollowRequest); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := "Follow request declined"
	if accept {
		msg = "Follow request accepted"
	}
	helper.Success(w, http.StatusOK, msg, nil)
}


// GetAllUsers GET /api/v1/user/all
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers(r.Context())
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Users retrieved successfully", users)
}

// GetFollowers GET /api/v1/user/followers?id=123
func (h *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	followers, err := h.Service.GetFollowers(r.Context(), userID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Followers retrieved successfully", followers)
}

// GetFollowing GET /api/v1/user/following?id=123
func (h *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	following, err := h.Service.GetFollowing(r.Context(), userID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Following list retrieved successfully", following)
}

// GetFollowRequests GET /api/v1/user/follow/requests
func (h *Handler) GetFollowRequests(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestsList, err := h.Service.GetPendingRequests(r.Context(), user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Pending requests retrieved successfully", requestsList)
}


// GetSuggestedUsers GET /api/v1/user/suggestions?limit=5
func (h *Handler) GetSuggestedUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	println("there is the user for suggest ",user)
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limit := 5
	//take five for now then i will take more 
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 20 {
		limit = 20
	}
	suggestions, err := h.Service.GetSuggestedUsers(r.Context(), user.ID, limit)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	//pass it to succes func 
	// println("not every thing  is okay we bring exactlry what we need ",suggestions)
	helper.Success(w, http.StatusOK, "Suggestions retrieved successfully", suggestions)
}
