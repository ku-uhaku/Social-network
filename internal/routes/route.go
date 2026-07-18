package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func Register(h *handler.Handler, m *middleware.Middleware) *http.ServeMux {
	mux := http.NewServeMux()
	onlyPost := m.AllowMethods(http.MethodPost)

	loginHandler := http.HandlerFunc(h.Login)

	mux.Handle("/login", onlyPost(loginHandler))
	mux.HandleFunc("GET /ws", h.WebSocket)
	mux.Handle("/api/profile", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetProfile))))
	return mux
}
