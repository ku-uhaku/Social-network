package handler

import (
	"log"
	"net/http"

	"kuu/internal/middleware"
	ws "kuu/internal/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Resolve user authorization details from the middleware context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Upgrade the connection to a persistent WebSocket pipe
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade failure: %v", err)
		return
	}

	// 3. Initialize client with the authenticated UserID
	client := &ws.Client{
		UserID: user.ID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.Hub.Register <- client

	go func() {
		onlineSnapshot := h.Hub.GetOnlineUsers()

		initEvent := ws.Event{
			Type:         "online_users_list",
			TargetUserID: user.ID, // Targeted exclusively to this user's connection
			Payload: map[string]interface{}{
				"online_user_ids": onlineSnapshot,
			},
		}

		h.Hub.Broadcast <- initEvent
	}()
	// ---------------------------------------------------------------------

	// Writer Loop (Outbound messages pushing data to client's browser)
	go func() {
		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()
}
