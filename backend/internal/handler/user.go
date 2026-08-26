package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/models"
	"kuu/internal/requests"
)

// UpdateProfile updates the current user's profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.UpdateProfilePayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if writeValidationErrors(w, requests.ValidateUpdateProfile(payload)) {
		return
	}

	if payload.IsPublic == 1 && user.IsPublic == 0 {
		if _, err := h.Service.AcceptAllPendingFollows(user.ID); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !h.expireNotificationsByType(w, user.ID, models.NotificationFollowRequest) {
			return
		}
	}

	updatedUser, err := h.Service.UpdateProfile(user.ID, payload)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Profile updated successfully", updatedUser)
}

// GetUserProfile GET /api/v1/user/profile?username=john
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		helper.Error(w, http.StatusBadRequest, "Invalid username parameter")
		return
	}
	targetUser, err := h.Service.GetUserProfile(user.ID, username)
	if err != nil {
		if (err==sql.ErrNoRows){
			helper.Error(w, http.StatusNotFound, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Only the owner, public profiles, or accepted followers can view posts
	if targetUser.FollowStatus != "self" && targetUser.IsPublic == 0 && targetUser.FollowStatus != "accepted" {
		// helper.Error(w, http.StatusForbidden, "This account is private")
		helper.ValidationErrorResponse(w, http.StatusOK, models.ErrorMsg{
			Code:    200,
			Message: "This account is private",
		})
		return
	}

	profile, err := h.Service.GetUserProfile(user.ID, username)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(w, http.StatusOK, "Profile retrieved successfully", profile)
}

// GetUserPosts GET /api/v1/user/posts?username=john
func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		helper.Error(w, http.StatusBadRequest, "Invalid username parameter")
		return
	}

	targetUser, err := h.Service.GetUserProfile(user.ID, username)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Private posts: owner or accepted followers only
	if targetUser.FollowStatus != "self" && targetUser.IsPublic == 0 && targetUser.FollowStatus != "accepted" {
		helper.Error(w, http.StatusForbidden, "This account is private")
		return
	}

	posts, err := h.Service.GetUserPosts(targetUser.ID, user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Posts retrieved successfully", posts)
}

// FollowUser POST /api/v1/user/follow
func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}
	var payload models.FollowActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if writeValidationErrors(w, requests.ValidateFollowAction(payload)) {
		return
	}

	status, err := h.Service.FollowUser(user.ID, payload.TargetUserID)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	msg := "Successfully followed user"
	if status == "pending" {
		msg = "Follow request sent successfully"
		// Avoid duplicate follow request notifications
		existing, _ := h.Service.GetNotificationByActorType(
			payload.TargetUserID, user.ID, models.NotificationFollowRequest,
		)
		if existing == nil || existing.IsExpired == 1 {
			actorID := user.ID
			_, err := h.DispatchNotification(payload.TargetUserID, &models.Notification{
				RecipientID: payload.TargetUserID,
				ActorID:     &actorID,
				Type:        models.NotificationFollowRequest,
				Title:       "New follow request",
				Message:     helper.DisplayName(user) + " wants to follow you",
				Actions:     acceptDeclineActions(),
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.FollowActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if err := h.Service.UnfollowUser(user.ID, payload.TargetUserID); err != nil {
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.FollowActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if err := h.Service.HandleFollowRequest(user.ID, payload.TargetUserID, accept); err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Expire follow request notifications
	if !h.expireNotificationsByType(w, user.ID, models.NotificationFollowRequest) {
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
	users, err := h.Service.GetAllUsers()
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Users retrieved successfully", users)
}

func (h *Handler) userList(w http.ResponseWriter, r *http.Request, fetch func(int64) ([]models.UserFollowView, error), message string) {
	userID, ok := queryPositiveInt64(w, r, "id", "Invalid user ID")
	if !ok {
		return
	}

	users, err := fetch(userID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, message, users)
}

// GetFollowers GET /api/v1/user/followers?id=123
func (h *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	h.userList(w, r, h.Service.GetFollowers, "Followers retrieved successfully")
}

// GetFollowing GET /api/v1/user/following?id=123
func (h *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	h.userList(w, r, h.Service.GetFollowing, "Following list retrieved successfully")
}

// GetFollowRequests GET /api/v1/user/follow/requests
func (h *Handler) GetFollowRequests(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	requestsList, err := h.Service.GetPendingRequests(user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Pending requests retrieved successfully", requestsList)
}

// GetSuggestedUsers GET /api/v1/user/suggestions?limit=5
func (h *Handler) GetSuggestedUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	limit := 5
	// take five for now then i will take more
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 20 {
		limit = 20
	}
	suggestions, err := h.Service.GetSuggestedUsers(user.ID, limit)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// pass it to succes func
	helper.Success(w, http.StatusOK, "Suggestions retrieved successfully", suggestions)
}
