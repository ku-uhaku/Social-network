package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerChatRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	mux.Handle("/api/v1/chat/direct", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetDirectHistory))))
	mux.Handle("/api/v1/chat/conversations", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetConversations))))
	mux.Handle("/api/v1/chat/group", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetGroupHistory))))
}
