package handler

import (
	"encoding/json"
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

// Upgrades the request to a WebSocket connection
func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
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

// Pushes queued events to the browser
func writeLoop(client *ws.Client) {
	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			// Connection broken: close to unblock readLoop
			client.Conn.Close()
			return
		}
	}
}

// Handles browser messages until disconnect
func (h *Handler) readLoop(client *ws.Client) {
	defer h.Hub.Unregister(client)

	for {
		_, raw, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		var event ws.Event
		if err := json.Unmarshal(raw, &event); err != nil {
			continue
		}

		switch event.Type {
		case ws.ClientSendDirectMessage:
			h.sendDirectMessage(client.UserID, event.Payload)
		case ws.ClientSendGroupMessage:
			h.sendGroupMessage(client.UserID, event.Payload)
		}
	}
}

func (h *Handler) sendDirectMessage(senderID int64, payload interface{}) {
	req, ok := decodeMessage(payload)
	if !ok || req.ReceiverID == nil {
		return
	}

	msg, err := h.Service.SaveDirectMessage(senderID, *req.ReceiverID, req.Content)
	if err != nil {
		return
	}

	// Echo to sender's other tabs
	h.Hub.BroadcastToUsers([]int64{senderID, *req.ReceiverID}, ws.EventNewDirectMessage, msg)
}

func (h *Handler) sendGroupMessage(senderID int64, payload interface{}) {
	req, ok := decodeMessage(payload)
	if !ok || req.GroupID == nil {
		return
	}

	msg, err := h.Service.SaveGroupMessage(senderID, *req.GroupID, req.Content)
	if err != nil {
		return
	}

	memberIDs, err := h.Service.GetGroupMemberIDs(*req.GroupID)
	if err != nil {
		return
	}

	h.Hub.BroadcastToUsers(memberIDs, ws.EventNewGroupMessage, msg)
}

// Re-reads payload as a message; rejects empty text
func decodeMessage(payload interface{}) (models.SendMessagePayload, bool) {
	var req models.SendMessagePayload

	raw, err := json.Marshal(payload)
	if err != nil || json.Unmarshal(raw, &req) != nil {
		return req, false
	}
	return req, strings.TrimSpace(req.Content) != ""
}
