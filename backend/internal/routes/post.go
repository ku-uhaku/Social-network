package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerPostRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	// Post Creation & Feed
	mux.Handle("/api/v1/posts", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.CreatePost))))
	mux.Handle("/api/v1/posts/feed", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetFeed))))

	// Comments
	mux.Handle("/api/v1/posts/comments", m.AllowMethods(http.MethodPost, http.MethodGet)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateComment(w, r)
		} else {
			h.GetComments(w, r)
		}
	}))))
}
