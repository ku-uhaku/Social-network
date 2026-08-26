package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"kuu/internal/handler"
	"kuu/internal/helper"
	"kuu/internal/middleware"
)

func Register(h *handler.Handler, m *middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/media/")

		// Don't allow access to /media/ itself or directories.
		if path == "" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		filePath := filepath.Join(helper.MediaDir, filepath.Clean(path))
		info, err := os.Stat(filePath)
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, filePath)
	})

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
