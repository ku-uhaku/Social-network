package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"kuu/internal/middleware"
	"kuu/internal/models"
	ws "kuu/internal/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade failure: %v", err)
		return
	}

	client := &ws.Client{
		UserID: user.ID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.Hub.Register <- client

	// Send initial list of online users to this connection
	go func() {
		onlineSnapshot := h.Hub.GetOnlineUsers()
		initEvent := ws.Event{
			Type:         ws.EventOnlineUsersList,
			TargetUserID: user.ID,
			Payload: map[string]interface{}{
				"online_user_ids": onlineSnapshot,
			},
		}
		h.Hub.Broadcast <- initEvent
	}()

	// Outbound Writer Loop (Pushes messages to client browser)
	go func() {
		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()

	// Inbound Reader Loop (Reads incoming frames from client)
	go h.handleIncomingWSMessages(r.Context(), client)
}

func (h *Handler) handleIncomingWSMessages(ctx context.Context, client *ws.Client) {
	defer func() {
		h.Hub.Unregister <- client
	}()

	for {
		_, rawMessage, err := client.Conn.ReadMessage()
		if err != nil {
			break // Connection closed or failed
		}

		var incoming ws.Event
		if err := json.Unmarshal(rawMessage, &incoming); err != nil {
			log.Printf("[WS] Invalid event payload from user %d: %v", client.UserID, err)
			continue
		}

		// Validate session on every client message
		token := ""
		if payloadMap, ok := incoming.Payload.(map[string]interface{}); ok {
			if v, ok := payloadMap["token"].(string); ok {
				token = v
			}
		}
		if token != "" {
			if _, err := h.Service.ValidateSession(ctx, token); err != nil {
				break
			}
		}

		switch incoming.Type {
		case ws.ClientSendDirectMessage:
			h.processDirectMessage(ctx, client.UserID, incoming.Payload)

		case ws.ClientSendGroupMessage:
			h.processGroupMessage(ctx, client.UserID, incoming.Payload)
		}
	}
}

func (h *Handler) processDirectMessage(ctx context.Context, senderID int64, rawPayload interface{}) {
	payloadBytes, _ := json.Marshal(rawPayload)
	var req models.SendMessagePayload
	if err := json.Unmarshal(payloadBytes, &req); err != nil || req.ReceiverID == nil || strings.TrimSpace(req.Content) == "" {
		return
	}

	// Save DM via Service
	msg, err := h.Service.SaveDirectMessage(ctx, senderID, *req.ReceiverID, req.Content)
	if err != nil {
		log.Printf("[WS] Failed to save DM: %v", err)
		return
	}

	// Dispatch real-time event to recipient and sender (for multi-tab sync)
	h.Hub.Broadcast <- ws.Event{
		Type:           ws.EventNewDirectMessage,
		TargetUsersIDs: []int64{senderID, *req.ReceiverID},
		Payload:        msg,
	}
}

func (h *Handler) processGroupMessage(ctx context.Context, senderID int64, rawPayload interface{}) {
	payloadBytes, _ := json.Marshal(rawPayload)
	var req models.SendMessagePayload
	if err := json.Unmarshal(payloadBytes, &req); err != nil || req.GroupID == nil || strings.TrimSpace(req.Content) == "" {
		return
	}

	// Save Group Message via Service (validates group membership inside Service)
	msg, err := h.Service.SaveGroupMessage(ctx, senderID, *req.GroupID, req.Content)
	if err != nil {
		log.Printf("[WS] Failed to save group message: %v", err)
		return
	}

	// Fetch active group members to route event
	memberIDs, err := h.Service.GetGroupMemberIDs(ctx, *req.GroupID)
	if err != nil {
		return
	}

	h.Hub.Broadcast <- ws.Event{
		Type:           ws.EventNewGroupMessage,
		TargetGroupID:  *req.GroupID,
		TargetUsersIDs: memberIDs,
		Payload:        msg,
	}
}
