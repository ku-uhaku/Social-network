package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func Register(h *handler.Handler, m *middleware.Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /health", h.Health)

	// Protected routes
	// mux.Handle("GET /profile", m.Auth(http.HandlerFunc(h.Profile)))
	mux.HandleFunc("GET /ws", h.WebSocket)
	return mux
}
