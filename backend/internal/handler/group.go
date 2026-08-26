package handler

import (
	"errors"
	"fmt"
	"net/http"

	"kuu/internal/helper"
	"kuu/internal/models"
	"kuu/internal/requests"
	"kuu/internal/service"
)

// CreateGroup POST /api/v1/groups
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.CreateGroupPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if writeValidationErrors(w, requests.ValidateCreateGroup(payload)) {
		return
	}

	group, err := h.Service.CreateGroup(user.ID, payload)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusCreated, "Group created successfully", group)
}

// GetGroup GET /api/v1/groups/detail?id=123
func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	groupID, ok := queryPositiveInt64(w, r, "id", "Invalid or missing group ID parameter")
	if !ok {
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

// GetAllGroups GET /api/v1/groups
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	groupID, ok := queryPositiveInt64(w, r, "id", "Invalid or missing group ID parameter")
	if !ok {
		return
	}

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

// UpdateGroup PUT /api/v1/groups?id=123
func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	groupID, ok := queryPositiveInt64(w, r, "id", "Invalid or missing group ID parameter")
	if !ok {
		return
	}

	var payload models.UpdateGroupPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if writeValidationErrors(w, requests.ValidateUpdateGroup(payload)) {
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

// DeleteGroup DELETE /api/v1/groups?id=123
func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	groupID, ok := queryPositiveInt64(w, r, "id", "Invalid or missing group ID parameter")
	if !ok {
		return
	}

	err := h.Service.DeleteGroup(user.ID, groupID)
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.InviteMembersPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if writeValidationErrors(w, requests.ValidateInviteMembers(payload)) {
		return
	}

	// Fetch group before notifying
	group, err := h.Service.GetGroupByID(payload.GroupID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
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
			Actions:     acceptDeclineActions(),
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.GroupActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if err := h.Service.RespondToInvitation(user.ID, payload.GroupID, accept); err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Expire invitation notification
	if !h.expireNotificationsByType(w, user.ID, models.NotificationGroupInvitation) {
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}
	var payload models.GroupActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if err := h.Service.JoinGroup(user.ID, payload.GroupID); err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Notify creator only if join request is pending
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
			Actions:     acceptDeclineActions(),
		}); err != nil {
			helper.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	helper.Success(w, http.StatusOK, "Join request processed", nil)
}

// LeaveGroup POST /api/v1/groups/leave
func (h *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.GroupActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if err := h.Service.LeaveGroup(user.ID, payload.GroupID); err != nil {
		fmt.Println("errrr:::",err)
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Successfully left the group", nil)
}

// GetGroupFeed GET /api/v1/groups/feed?id=123&limit=10&cursor=456
// Returns { posts, next_cursor, has_more }
func (h *Handler) GetGroupFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	groupID, ok := queryPositiveInt64(w, r, "id", "Invalid or missing group ID parameter")
	if !ok {
		return
	}

	limit, cursor, ok := parseFeedParams(w, r)
	if !ok {
		return
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
	user, ok := h.authUser(w, r)
	if !ok {
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	groupID, ok := queryPositiveInt64(w, r, "id", "Invalid or missing group ID parameter")
	if !ok {
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
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.CreateGroupEventPayload
	if !decodeJSON(w, r, &payload) {
		return
	}

	if writeValidationErrors(w, requests.ValidateCreateGroupEvent(payload)) {
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

func writeEventActionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrEventNotFound):
		helper.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrNotGroupCreator), errors.Is(err, service.ErrNotMember):
		helper.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrEventExpired):
		helper.Error(w, http.StatusBadRequest, err.Error())
	default:
		helper.Error(w, http.StatusInternalServerError, err.Error())
	}
}

// CancelGroupEvent POST /api/v1/groups/events/cancel (creator only)
func (h *Handler) CancelGroupEvent(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload struct {
		EventID int64 `json:"event_id"`
	}
	if !decodeJSON(w, r, &payload) {
		return
	}
	if payload.EventID <= 0 {
		helper.Error(w, http.StatusBadRequest, "event_id is required")
		return
	}

	if err := h.Service.CancelGroupEvent(user.ID, payload.EventID); err != nil {
		writeEventActionError(w, err)
		return
	}

	helper.Success(w, http.StatusOK, "Group event cancelled", nil)
}

// SetEventResponse POST /api/v1/groups/events/respond (members only)
func (h *Handler) SetEventResponse(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.EventResponsePayload
	if !decodeJSON(w, r, &payload) {
		return
	}
	if writeValidationErrors(w, requests.ValidateEventResponse(payload)) {
		return
	}

	updated, err := h.Service.SetEventResponse(user.ID, payload.EventID, payload.Status)
	if err != nil {
		writeEventActionError(w, err)
		return
	}

	helper.Success(w, http.StatusOK, "Event response recorded", updated)
}

// HandleJoinRequestAction POST /api/v1/groups/join/accept|decline (creator only)
func (h *Handler) HandleJoinRequestAction(w http.ResponseWriter, r *http.Request, accept bool) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.GroupActionPayload
	if !decodeJSON(w, r, &payload) {
		return
	}
	if writeValidationErrors(w, requests.ValidateJoinRequestAction(payload)) {
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

	// Expire join request notification
	if !h.expireNotificationsByType(w, user.ID, models.NotificationGroupJoinRequest) {
		return
	}

	msg := "Join request declined"
	if accept {
		msg = "Join request accepted"
	}
	helper.Success(w, http.StatusOK, msg, nil)
}
