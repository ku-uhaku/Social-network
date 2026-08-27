package handler

import (
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/models"
	"kuu/internal/requests"
	ws "kuu/internal/websocket"
)

const (
	defaultNotificationLimit = 20
	maxNotificationLimit     = 50
)

// NotificationList GET /api/v1/notifications?limit=20&last_id=0
func (h *Handler) NotificationList(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	limit, ok := notificationLimit(w, r)
	if !ok {
		return
	}
	lastID, ok := notificationCursor(w, r)
	if !ok {
		return
	}

	resp, err := h.Service.GetUserNotifications(user.ID, limit, lastID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Notifications retrieved", resp)
}

// MarkNotificationRead POST /api/v1/notifications/read
func (h *Handler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.MarkNotificationReadPayload
	if !decodeJSON(w, r, &payload) {
		return
	}
	if writeValidationErrors(w, requests.ValidateMarkNotificationRead(payload)) {
		return
	}

	var notificationID int64
	if payload.NotificationID != nil {
		notificationID = *payload.NotificationID
		if err := h.Service.MarkNotificationRead(user.ID, notificationID); err != nil {
			helper.Error(w, http.StatusNotFound, err.Error())
			return
		}
	} else if err := h.Service.MarkAllNotificationsRead(user.ID); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Sync other tabs in real time
	h.Hub.BroadcastToUser(user.ID, ws.EventNotificationRead, map[string]interface{}{
		"notification_id": notificationID,
		"all":             payload.NotificationID == nil,
	})

	helper.Success(w, http.StatusOK, "Notification marked as read", nil)
}

// ExpireNotification POST /api/v1/notifications/expire
func (h *Handler) ExpireNotification(w http.ResponseWriter, r *http.Request) {
	user, ok := h.authUser(w, r)
	if !ok {
		return
	}

	var payload models.ExpireNotificationPayload
	if !decodeJSON(w, r, &payload) {
		return
	}
	if writeValidationErrors(w, requests.ValidateExpireNotification(payload)) {
		return
	}

	if err := h.Service.ExpireNotification(user.ID, payload.NotificationID); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Sync other tabs in real time
	h.Hub.BroadcastToUser(user.ID, ws.EventNotificationExpired, map[string]interface{}{
		"notification_id": payload.NotificationID,
	})

	helper.Success(w, http.StatusOK, "Notification expired", nil)
}

// DispatchNotification creates a notification and pushes it to the recipient
func (h *Handler) DispatchNotification(recipientID int64, n *models.Notification) (*models.Notification, error) {
	created, err := h.Service.CreateNotification(n)
	if err != nil {
		return nil, err
	}

	h.Hub.BroadcastToUser(recipientID, ws.EventNewNotification, created)
	return created, nil
}

// expireNotificationsByType expires every unread notification of a type. Only
// for events that really do resolve all of them at once.
func (h *Handler) expireNotificationsByType(w http.ResponseWriter, userID int64, notifType string) bool {
	ids, err := h.Service.ExpireNotificationsByType(userID, notifType)
	return h.pushExpired(w, userID, ids, err)
}

// expireActorNotifications expires the notifications one actor raised, so
// answering one request leaves the other senders' notifications alone.
func (h *Handler) expireActorNotifications(w http.ResponseWriter, userID, actorID int64, notifType string) bool {
	ids, err := h.Service.ExpireNotificationsByActorType(userID, actorID, notifType)
	return h.pushExpired(w, userID, ids, err)
}

// expireGroupNotifications expires the notifications raised for one group
func (h *Handler) expireGroupNotifications(w http.ResponseWriter, userID, groupID int64, notifType string) bool {
	ids, err := h.Service.ExpireGroupNotifications(userID, groupID, notifType)
	return h.pushExpired(w, userID, ids, err)
}

// pushExpired reports the expiry to the user's open tabs; writes 500 and
// returns false if the expiry itself failed.
func (h *Handler) pushExpired(w http.ResponseWriter, userID int64, ids []int64, err error) bool {
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return false
	}
	for _, id := range ids {
		h.Hub.BroadcastToUser(userID, ws.EventNotificationExpired, map[string]interface{}{
			"notification_id": id,
		})
	}
	return true
}

// acceptDeclineActions is the standard button pair for request-style notifications
func acceptDeclineActions() models.JSONText {
	return models.JSONText{
		"buttons": []interface{}{
			map[string]interface{}{"action": "accept", "label": "Accept"},
			map[string]interface{}{"action": "decline", "label": "Decline"},
		},
	}
}

// notificationLimit reads ?limit, falling back to the default when absent
func notificationLimit(w http.ResponseWriter, r *http.Request) (int, bool) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return defaultNotificationLimit, true
	}

	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 {
		helper.Error(w, http.StatusBadRequest, "Invalid limit parameter")
		return 0, false
	}
	if limit > maxNotificationLimit {
		limit = maxNotificationLimit
	}
	return limit, true
}

// notificationCursor reads ?last_id, where 0 or absent means the newest page
func notificationCursor(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := r.URL.Query().Get("last_id")
	if raw == "" {
		return 0, true
	}

	lastID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || lastID < 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid last_id parameter")
		return 0, false
	}
	return lastID, true
}
