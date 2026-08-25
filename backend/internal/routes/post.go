package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerPostRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	// Post Creation, retrieval & feed (rate limited at the mux level)
	mux.Handle("/api/v1/posts", m.AllowMethods(http.MethodPost, http.MethodGet)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreatePost(w, r)
			return
		}
		h.GetPost(w, r)
	}))))
	mux.Handle("/api/v1/posts/feed", m.AllowMethods(http.MethodGet)(
		m.RequireAuth(http.HandlerFunc(h.GetFeed)),
	))

	// Comments
	mux.Handle("/api/v1/posts/comments", m.AllowMethods(http.MethodPost, http.MethodGet)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateComment(w, r)
		} else {
			h.GetComments(w, r)
		}
	}))))
}
