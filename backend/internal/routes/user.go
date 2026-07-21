package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerUserRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	mux.Handle(
		"/api/v1/user/profile/update",
		m.AllowMethods(http.MethodPut)(m.RequireAuth(http.HandlerFunc(h.UpdateProfile))),
	)
}
