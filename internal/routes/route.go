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

	// --- Auth Group Endpoints ---
	mux.Handle(
		"/api/v1/auth/register",
		m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Register)),
	)
	mux.Handle(
		"/api/v1/auth/login",
		m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Login)),
	)
	// mux.Handle("/api/v1/auth/logout",
	// 	m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.Logout))),
	// )

	return mux
}
