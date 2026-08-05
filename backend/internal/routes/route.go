package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

// Register initializes the primary multiplexer and mounts feature sub-routes
func Register(h *handler.Handler, m *middleware.Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("media"))))

	// WebSocket / Realtime (Protected by Method Check & Auth Middleware)
	mux.Handle("/ws", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.WebSocket))))

	// Sub-route modules
	registerAuthRoutes(mux, h, m)
	registerUserRoutes(mux, h, m)
	registerGroupRoutes(mux, h, m)
	registerPostRoutes(mux, h, m)
	registerChatRoutes(mux, h, m)
	registerNotificationRoutes(mux, h, m)

	return mux
}
