package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerNotificationRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	mux.Handle("/api/v1/notifications", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.NotificationList))))
	mux.Handle("/api/v1/notifications/read", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.MarkNotificationRead))))
	mux.Handle("/api/v1/notifications/expire", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.ExpireNotification))))
}