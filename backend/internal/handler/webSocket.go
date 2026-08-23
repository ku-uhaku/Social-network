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

// WebSocket upgrades the request and keeps the connection open for the session.
func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] upgrade failed: %v", err)
		return
	}

	client := &ws.Client{
		UserID: user.ID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	h.Hub.Register(client)

	go writeLoop(client)
	go h.readLoop(client)
}

// writeLoop pushes queued events to the browser until the hub closes Send.
func writeLoop(client *ws.Client) {
	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// readLoop handles what the browser sends until the connection drops.
func (h *Handler) readLoop(client *ws.Client) {
	defer h.Hub.Unregister(client)

	ctx := context.Background()
	for {
		_, raw, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		var event ws.Event
		if err := json.Unmarshal(raw, &event); err != nil {
			log.Printf("[WS] invalid event from user %d: %v", client.UserID, err)
			continue
		}

		switch event.Type {
		case ws.ClientSendDirectMessage:
			h.sendDirectMessage(ctx, client.UserID, event.Payload)
		case ws.ClientSendGroupMessage:
			h.sendGroupMessage(ctx, client.UserID, event.Payload)
		}
	}
}

func (h *Handler) sendDirectMessage(ctx context.Context, senderID int64, payload interface{}) {
	req, ok := decodeMessage(payload)
	if !ok || req.ReceiverID == nil {
		return
	}

	msg, err := h.Service.SaveDirectMessage(ctx, senderID, *req.ReceiverID, req.Content)
	if err != nil {
		log.Printf("[WS] cannot save direct message: %v", err)
		return
	}

	// The sender gets it back too, so their other tabs stay in sync.
	h.Hub.BroadcastToUsers([]int64{senderID, *req.ReceiverID}, ws.EventNewDirectMessage, msg)
}

func (h *Handler) sendGroupMessage(ctx context.Context, senderID int64, payload interface{}) {
	req, ok := decodeMessage(payload)
	if !ok || req.GroupID == nil {
		return
	}

	// The service checks that the sender belongs to the group.
	msg, err := h.Service.SaveGroupMessage(ctx, senderID, *req.GroupID, req.Content)
	if err != nil {
		log.Printf("[WS] cannot save group message: %v", err)
		return
	}

	memberIDs, err := h.Service.GetGroupMemberIDs(ctx, *req.GroupID)
	if err != nil {
		log.Printf("[WS] cannot load members of group %d: %v", *req.GroupID, err)
		return
	}

	h.Hub.BroadcastToUsers(memberIDs, ws.EventNewGroupMessage, msg)
}

// decodeMessage re-reads the generic payload as a message and rejects empty text.
func decodeMessage(payload interface{}) (models.SendMessagePayload, bool) {
	var req models.SendMessagePayload

	raw, err := json.Marshal(payload)
	if err != nil || json.Unmarshal(raw, &req) != nil {
		return req, false
	}
	return req, strings.TrimSpace(req.Content) != ""
}
