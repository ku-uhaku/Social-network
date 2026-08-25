package routes

import (
	"net/http"
	"time"

	"kuu/internal/handler"
	"kuu/internal/helper"
	"kuu/internal/middleware"
)

func Register(h *handler.Handler, m *middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(helper.MediaDir))))
	mux.Handle("/ws", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.WebSocket))))

	registerAuthRoutes(mux, h, m)
	registerUserRoutes(mux, h, m)
	registerGroupRoutes(mux, h, m)
	registerPostRoutes(mux, h, m)
	registerChatRoutes(mux, h, m)
	registerNotificationRoutes(mux, h, m)

	limiter := helper.Neewratelimeter(time.Minute)
	return limiter.Wraponall("api", mux)
}
