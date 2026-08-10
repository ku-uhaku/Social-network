package routes

import (
	"net/http"

	"kuu/internal/handler"
	"kuu/internal/middleware"
)

func registerGroupRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	// Create group
	mux.Handle(
		"/api/v1/groups",
		m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.CreateGroup))),
	)

	// Fetch all groups
	mux.Handle(
		"/api/v1/groups/all",
		m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetAllGroups))),
	)

	// Fetch single group details (/api/v1/groups/detail?id=123)
	mux.Handle(
		"/api/v1/groups/detail",
		m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetGroup))),
	)

	// Fetch a single group's post feed (/api/v1/groups/feed?id=123)
	mux.Handle(
		"/api/v1/groups/feed",
		m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetGroupFeed))),
	)

	// Update an existing group (/api/v1/groups/update?id=123)
	mux.Handle(
		"/api/v1/groups/update",
		m.AllowMethods(http.MethodPut)(m.RequireAuth(http.HandlerFunc(h.UpdateGroup))),
	)

	// Delete a group (/api/v1/groups/delete?id=123)
	mux.Handle(
		"/api/v1/groups/delete",
		m.AllowMethods(http.MethodDelete)(m.RequireAuth(http.HandlerFunc(h.DeleteGroup))),
	)
	// Invitations
	mux.Handle("/api/v1/groups/invite", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.InviteMembers))))
	mux.Handle("/api/v1/groups/invitations", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetMyInvitations))))
	mux.Handle("/api/v1/groups/invitations/accept", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.AcceptInvitation))))
	mux.Handle("/api/v1/groups/invitations/decline", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.DeclineInvitation))))

	// Join
	mux.Handle("/api/v1/groups/join", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.JoinGroup))))
	mux.Handle("/api/v1/groups/join/accept", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleJoinRequestAction(w, r, true)
	}))))
	mux.Handle("/api/v1/groups/join/decline", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleJoinRequestAction(w, r, false)
	}))))
	mux.Handle("/api/v1/groups/leave", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.LeaveGroup))))

	// event
	mux.Handle("/api/v1/groups/events", m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.GetGroupEvents))))
	mux.Handle("/api/v1/groups/events/create", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.CreateGroupEvent))))
	mux.Handle("/api/v1/groups/events/cancel", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.CancelGroupEvent))))
	mux.Handle("/api/v1/groups/events/respond", m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.SetEventResponse))))
}
