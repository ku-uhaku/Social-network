package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

// Register initializes the primary multiplexer and mounts feature sub-routes
func Register(h *handler.Handler, m *middleware.Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	// WebSocket / Realtime
	mux.HandleFunc("GET /ws", h.WebSocket)

	// Sub-route modules
	registerAuthRoutes(mux, h, m)
	registerUserRoutes(mux, h, m)
	registerGroupRoutes(mux, h, m)
	registerPostRoutes(mux, h, m)

	return mux
}
