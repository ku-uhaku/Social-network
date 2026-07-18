package handler

import (
	"net/http"

	ws "kuu/internal/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.Hub.Register <- client

	go func() {
		defer func() {
			h.Hub.Unregister <- client
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			h.Hub.Broadcast <- msg
		}
	}()

	go func() {
		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()
}
