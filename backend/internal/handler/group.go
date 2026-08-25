package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/requests"
	"kuu/internal/service"
)

// CreateGroup handles POST /api/v1/groups
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
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

	if validationErrs := requests.ValidateCreateGroup(payload); len(validationErrs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, validationErrs)
		return
	}

	group, err := h.Service.CreateGroup(user.ID, payload)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusCreated, "Group created successfully", group)
}

// GetGroup handles GET /api/v1/groups/detail?id=123
// Returns { group, membership } where membership is the requesting user's
// status in the group ('accepted', 'pending', or 'none').
func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid or missing group ID parameter")
		return
	}

	group, err := h.Service.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, service.ErrGroupNotFound) {
			helper.Error(w, http.StatusNotFound, "Group not found")
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	membership, err := h.Service.GetMembershipStatus(groupID, user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group retrieved successfully", map[string]interface{}{
		"group":      group,
		"membership": membership,
	})
}

// GetAllGroups handles GET /api/v1/groups
func (h *Handler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.Service.GetAllGroups()
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Groups retrieved successfully", groups)
}

// GetGroupMembers GET /api/v1/groups/members?id=123 (accepted members only)
func (h *Handler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid or missing group ID parameter")
		return
	}

	// Only accepted members may see the member list
	membership, err := h.Service.GetMembershipStatus(groupID, user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if membership != "accepted" {
		helper.Error(w, http.StatusForbidden, service.ErrNotMember.Error())
		return
	}

	members, err := h.Service.GetGroupMembers(groupID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group members retrieved", members)
}

// UpdateGroup handles PUT /api/v1/groups?id=123
func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid or missing group ID parameter")
		return
	}

	var payload models.UpdateGroupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Malformed JSON request body")
		return
	}
	defer r.Body.Close()

	if validationErrs := requests.ValidateUpdateGroup(payload); len(validationErrs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, validationErrs)
		return
	}

	updatedGroup, err := h.Service.UpdateGroup(user.ID, groupID, payload)
	if err != nil {
		if errors.Is(err, service.ErrGroupNotFound) {
			helper.Error(w, http.StatusNotFound, "Group not found")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			helper.Error(w, http.StatusForbidden, "Only the group creator can update group details")
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group updated successfully", updatedGroup)
}

// DeleteGroup handles DELETE /api/v1/groups?id=123
func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Authentication context missing")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid or missing group ID parameter")
		return
	}

	err = h.Service.DeleteGroup(user.ID, groupID)
	if err != nil {
		if errors.Is(err, service.ErrGroupNotFound) {
			helper.Error(w, http.StatusNotFound, "Group not found")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			helper.Error(w, http.StatusForbidden, "Only the group creator can delete this group")
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group deleted successfully", nil)
}

// InviteMembers POST /api/v1/groups/invite
func (h *Handler) InviteMembers(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.InviteMembersPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if errs := requests.ValidateInviteMembers(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	if err := h.Service.InviteUsers(user.ID, payload); err != nil {
		if errors.Is(err, service.ErrNotMember) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Notify each invited user
	group, err := h.Service.GetGroupByID(payload.GroupID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	actorID := user.ID
	for _, targetID := range payload.TargetUserIDs {
		if targetID == user.ID {
			continue
		}
		if _, err := h.DispatchNotification(targetID, &models.Notification{
			RecipientID: targetID,
			ActorID:     &actorID,
			Type:        models.NotificationGroupInvitation,
			Title:       "Group invitation",
			Message:     helper.DisplayName(user) + " invited you to join " + group.Title,
			Payload:     models.JSONText{"group_id": group.ID},
			Actions: models.JSONText{
				"buttons": []interface{}{
					map[string]interface{}{"action": "accept", "label": "Accept"},
					map[string]interface{}{"action": "decline", "label": "Decline"},
				},
			},
		}); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.Success(w, http.StatusOK, "Invitations sent successfully", nil)
}

// AcceptInvitation POST /api/v1/groups/invitations/accept
func (h *Handler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	h.respondInvite(w, r, true)
}

// DeclineInvitation POST /api/v1/groups/invitations/decline
func (h *Handler) DeclineInvitation(w http.ResponseWriter, r *http.Request) {
	h.respondInvite(w, r, false)
}

func (h *Handler) respondInvite(w http.ResponseWriter, r *http.Request, accept bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.GroupActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.Service.RespondToInvitation(user.ID, payload.GroupID, accept); err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Expire the group_invitation notification once acted upon
	if err := h.Service.ExpireNotificationsByType(user.ID, models.NotificationGroupInvitation); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := "Invitation declined"
	if accept {
		msg = "Invitation accepted successfully"
	}
	helper.Success(w, http.StatusOK, msg, nil)
}

// JoinGroup POST /api/v1/groups/join
func (h *Handler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var payload models.GroupActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.Service.JoinGroup(user.ID, payload.GroupID); err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Only notify the creator when a pending join request was created
	// (private group); public joins are accepted immediately.
	membership, err := h.Service.GetMembershipStatus(payload.GroupID, user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if membership == "pending" {
		group, err := h.Service.GetGroupByID(payload.GroupID)
		if err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		actorID := user.ID
		if _, err := h.DispatchNotification(group.CreatorID, &models.Notification{
			RecipientID: group.CreatorID,
			ActorID:     &actorID,
			Type:        models.NotificationGroupJoinRequest,
			Title:       "Group join request",
			Message:     helper.DisplayName(user) + " requested to join " + group.Title,
			Payload:     models.JSONText{"group_id": group.ID, "target_user_id": user.ID},
			Actions: models.JSONText{
				"buttons": []interface{}{
					map[string]interface{}{"action": "accept", "label": "Accept"},
					map[string]interface{}{"action": "decline", "label": "Decline"},
				},
			},
		}); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	helper.Success(w, http.StatusOK, "Join request processed", nil)
}

// LeaveGroup POST /api/v1/groups/leave
func (h *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.GroupActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.Service.LeaveGroup(user.ID, payload.GroupID); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Successfully left the group", nil)
}

// GetGroupFeed GET /api/v1/groups/feed?id=123&limit=10&cursor=456
// Returns { posts, next_cursor, has_more } like the home feed.
func (h *Handler) GetGroupFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid or missing group ID parameter")
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil || parsed < 1 {
			helper.Error(w, http.StatusBadRequest, "Invalid limit parameter")
			return
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
			return
		}
		cursor = &parsed
	}

	posts, hasMore, err := h.Service.GetGroupFeed(user.ID, groupID, limit, cursor)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var nextCursor *int64
	if hasMore && len(posts) > 0 {
		lastID := posts[len(posts)-1].ID
		nextCursor = &lastID
	}

	helper.Success(w, http.StatusOK, "Group feed retrieved successfully", map[string]interface{}{
		"posts":       posts,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

// GetMyInvitations GET /api/v1/groups/invitations
func (h *Handler) GetMyInvitations(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	invites, err := h.Service.GetPendingInvitations(user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Pending invitations retrieved", invites)
}

// GetGroupEvents GET /api/v1/groups/events?id=123 (members only)
func (h *Handler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groupIDStr := r.URL.Query().Get("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil || groupID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid or missing group ID parameter")
		return
	}

	events, err := h.Service.GetGroupEvents(user.ID, groupID)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Group events retrieved", events)
}

// CreateGroupEvent POST /api/v1/groups/events (members only)
func (h *Handler) CreateGroupEvent(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.CreateGroupEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if errs := requests.ValidateCreateGroupEvent(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	event, err := h.Service.CreateGroupEvent(user.ID, payload)
	if err != nil {
		if errors.Is(err, service.ErrNotMember) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Notify other group members about the new event
	if group, err := h.Service.GetGroupByID(payload.GroupID); err == nil {
		if err := h.notifyEventCreated(user.ID, group, event); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.Success(w, http.StatusCreated, "Group event created", event)
}

// notifyEventCreated dispatches group_event_created notifications to other members
func (h *Handler) notifyEventCreated(actorID int64, group *models.Group, event *models.GroupEvent) error {
	memberIDs, err := h.Service.GetGroupMemberIDs(group.ID)
	if err != nil {
		return err
	}
	actor := actorID
	for _, memberID := range memberIDs {
		if memberID == actorID {
			continue
		}
		if _, err := h.DispatchNotification(memberID, &models.Notification{
			RecipientID: memberID,
			ActorID:     &actor,
			Type:        models.NotificationGroupEvent,
			Title:       "New group event",
			Message:     "A new event was created in " + group.Title + ": " + event.Title,
			Payload:     models.JSONText{"group_id": group.ID, "event_id": event.ID},
			Actions: models.JSONText{
				"buttons": []interface{}{
					map[string]interface{}{"action": "view", "label": "View event"},
				},
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

// CancelGroupEvent POST /api/v1/groups/events/cancel (creator only)
func (h *Handler) CancelGroupEvent(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload struct {
		EventID int64 `json:"event_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if payload.EventID <= 0 {
		helper.Error(w, http.StatusBadRequest, "event_id is required")
		return
	}

	if err := h.Service.CancelGroupEvent(user.ID, payload.EventID); err != nil {
		switch {
		case errors.Is(err, service.ErrEventNotFound):
			helper.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotGroupCreator):
			helper.Error(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrEventExpired):
			helper.Error(w, http.StatusBadRequest, err.Error())
		default:
			helper.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.Success(w, http.StatusOK, "Group event cancelled", nil)
}

// SetEventResponse POST /api/v1/groups/events/respond (members only)
func (h *Handler) SetEventResponse(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.EventResponsePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if errs := requests.ValidateEventResponse(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	updated, err := h.Service.SetEventResponse(user.ID, payload.EventID, payload.Status)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEventNotFound):
			helper.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotMember):
			helper.Error(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrEventExpired):
			helper.Error(w, http.StatusBadRequest, err.Error())
		default:
			helper.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.Success(w, http.StatusOK, "Event response recorded", updated)
}

// HandleJoinRequestAction POST /api/v1/groups/join/accept|decline (creator only)
func (h *Handler) HandleJoinRequestAction(w http.ResponseWriter, r *http.Request, accept bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.GroupActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if errs := requests.ValidateJoinRequestAction(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	if err := h.Service.HandleJoinRequest(user.ID, payload.GroupID, payload.TargetUserID, accept); err != nil {
		if errors.Is(err, service.ErrNotGroupCreator) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Expire the join request notification once acted upon
	if err := h.Service.ExpireNotificationsByType(user.ID, models.NotificationGroupJoinRequest); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := "Join request declined"
	if accept {
		msg = "Join request accepted"
	}
	helper.Success(w, http.StatusOK, msg, nil)
}
