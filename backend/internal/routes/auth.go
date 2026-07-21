package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerAuthRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	mux.Handle(
		"/api/v1/auth/register",
		m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Register)),
	)
	mux.Handle(
		"/api/v1/auth/login",
		m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Login)),
	)
	mux.Handle(
		"/api/v1/auth/me",
		m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.Me))),
	)
	mux.Handle(
		"/api/v1/auth/logout",
		m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.Logout))),
	)
}
