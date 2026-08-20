package routes

import (
	"net/http"
	"time"

	"kuu/internal/handler"
	"kuu/internal/helper"
	"kuu/internal/middleware"
)

func registerPostRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	ratelimeer := helper.Neewratelimeter(time.Minute)
	// Post Creation, retrieval, & feed
	mux.Handle("/api/v1/posts", m.AllowMethods(http.MethodPost, http.MethodGet)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreatePost(w, r)
			return
		}
		h.GetPost(w, r)
	}))))
	// shooould i apply ratelimmet first
	// mux.Handle("/api/v1/posts/feed", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetFeed))))
	mux.Handle("/api/v1/posts/feed", ratelimeer.Wraponall("api", m.AllowMethods(http.MethodGet)(
		m.RequireAuth(http.HandlerFunc(h.GetFeed)),),
	),
	)

	// Comments 
	mux.Handle("/api/v1/posts/comments",ratelimeer.Wraponall("api", m.AllowMethods(http.MethodPost, http.MethodGet)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateComment(w, r)
		} else {
			h.GetComments(w, r)
		}
	})))))
}
