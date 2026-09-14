package routes

import (
	"net/http"
	"time"

	"kuu/internal/handler"
	"kuu/internal/helper"
	"kuu/internal/middleware"
)

func registerAuthRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	limiter := helper.NewRateLimiter(time.Minute)

	mux.Handle("/api/v1/auth/register",
		limiter.Wrap("auth", m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Register))),
	)
	mux.Handle("/api/v1/auth/login",
		limiter.Wrap("auth", m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Login))),
	)
	mux.Handle("/api/v1/auth/logout",
		limiter.Wrap("auth", m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Logout))),
	)
	mux.Handle(
		"/api/v1/auth/me",
		m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.Me))),
	)
}
