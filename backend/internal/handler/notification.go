package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/models"
	"kuu/internal/requests"
	ws "kuu/internal/websocket"
)

// NotificationList GET /api/v1/notifications?limit=20&last_id=0
func (h *Handler) NotificationList(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	lastID, _ := strconv.ParseInt(r.URL.Query().Get("last_id"), 10, 64)
	if lastID < 0 {
		lastID = 0
	}

	resp, err := h.Service.GetUserNotifications(r.Context(), user.ID, limit, lastID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Notifications retrieved", resp)
}

// MarkNotificationRead POST /api/v1/notifications/read
func (h *Handler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.MarkNotificationReadPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if errs := requests.ValidateMarkNotificationRead(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	var notificationID int64
	if payload.NotificationID != nil {
		notificationID = *payload.NotificationID
	}

	id, err := h.Service.MarkNotificationRead(r.Context(), user.ID, notificationID)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Sync other tabs in real time
	h.Hub.BroadcastToUser(user.ID, ws.EventNotificationRead, map[string]interface{}{
		"notification_id": id,
		"all":             payload.NotificationID == nil,
	})

	helper.Success(w, http.StatusOK, "Notification marked as read", nil)
}

// ExpireNotification POST /api/v1/notifications/expire
func (h *Handler) ExpireNotification(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var payload models.ExpireNotificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if errs := requests.ValidateExpireNotification(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	if err := h.Service.ExpireNotification(r.Context(), user.ID, payload.NotificationID); err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Sync other tabs in real time
	h.Hub.BroadcastToUser(user.ID, ws.EventNotificationExpired, map[string]interface{}{
		"notification_id": payload.NotificationID,
	})

	helper.Success(w, http.StatusOK, "Notification expired", nil)
}

// DispatchNotification is the single contract phase-2 handlers call to create + broadcast
func (h *Handler) DispatchNotification(ctx context.Context, recipientID int64, n *models.Notification) (*models.Notification, error) {
	created, err := h.Service.CreateNotification(ctx, n)
	if err != nil {
		return nil, err
	}
	h.Hub.BroadcastToUser(recipientID, ws.EventNewNotification, created)
	return created, nil
}

// ExpireAndPush is the single contract phase-2 handlers call to expire + broadcast
func (h *Handler) ExpireAndPush(ctx context.Context, recipientID, notificationID int64) error {
	if err := h.Service.ExpireNotification(ctx, recipientID, notificationID); err != nil {
		return err
	}
	h.Hub.BroadcastToUser(recipientID, ws.EventNotificationExpired, map[string]interface{}{
		"notification_id": notificationID,
	})
	return nil
}

// TODO: might need to refactor each notification type into its file
// ExpireFollowRequest finds and expires a follow_request notification for a recipient+requester pair
func (h *Handler) ExpireFollowRequest(ctx context.Context, recipientID, requesterID int64) error {
	n, err := h.Service.GetNotificationByActorType(ctx, recipientID, requesterID, models.NotificationFollowRequest)
	if err != nil {
		return err
	}
	if n == nil {
		return nil
	}
	return h.ExpireAndPush(ctx, recipientID, n.ID)
}
