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

// WebSocket upgrades the request and keeps the connection open for the session.
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

// writeLoop pushes queued events to the browser until the hub closes Send.
func writeLoop(client *ws.Client) {
	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			// The connection is broken: close it so readLoop unblocks and
			// unregisters the client instead of leaving it in the hub.
			client.Conn.Close()
			return
		}
	}
}

// readLoop handles what the browser sends until the connection drops.
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

	// The sender gets it back too, so their other tabs stay in sync.
	h.Hub.BroadcastToUsers([]int64{senderID, *req.ReceiverID}, ws.EventNewDirectMessage, msg)
}

func (h *Handler) sendGroupMessage(senderID int64, payload interface{}) {
	req, ok := decodeMessage(payload)
	if !ok || req.GroupID == nil {
		return
	}

	// The service checks that the sender belongs to the group.
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

// decodeMessage re-reads the generic payload as a message and rejects empty text.
func decodeMessage(payload interface{}) (models.SendMessagePayload, bool) {
	var req models.SendMessagePayload

	raw, err := json.Marshal(payload)
	if err != nil || json.Unmarshal(raw, &req) != nil {
		return req, false
	}
	return req, strings.TrimSpace(req.Content) != ""
}
