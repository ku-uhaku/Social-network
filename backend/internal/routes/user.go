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

	mux.Handle("/api/v1/user/follow", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.FollowUser))))
	mux.Handle("/api/v1/user/unfollow", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.UnfollowUser))))

	// Follow Request Approvals
	mux.Handle("/api/v1/user/follow/accept", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.AcceptFollowRequest))))
	mux.Handle("/api/v1/user/follow/decline", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.DeclineFollowRequest))))
	mux.Handle("/api/v1/user/follow/requests", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetFollowRequests))))

	// Lists
	mux.Handle("/api/v1/user/followers", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetFollowers))))
	mux.Handle("/api/v1/user/following", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetFollowing))))
}
