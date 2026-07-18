package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func Register(h *handler.Handler, m *middleware.Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws", h.WebSocket)
	// --- Public Group Routes ---

	return mux
}
