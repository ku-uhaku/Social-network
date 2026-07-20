package handler

import (
	"encoding/json"
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

	// Reader Loop (Inbound messages from the client's browser)
	go func() {
		defer func() {
			h.Hub.Unregister <- client
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// Parse whatever JSON the frontend sends into your flexible Event struct
			var incomingEvent ws.Event
			if err := json.Unmarshal(msg, &incomingEvent); err != nil {
				log.Printf("[WS] Failed to parse incoming JSON payload from user %d: %v", user.ID, err)
				continue
			}

			// Optional safety check: Ensure the payload always tracks who actually sent it
			if incomingEvent.Payload == nil {
				incomingEvent.Payload = make(map[string]interface{})
			}
			if m, ok := incomingEvent.Payload.(map[string]interface{}); ok {
				m["sender_id"] = user.ID
			}

			// Forward the event to the Hub broadcast logic cleanly
			h.Hub.Broadcast <- incomingEvent
		}
	}()

	// Writer Loop (Outbound messages pushing data to client's browser)
	go func() {
		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()
}