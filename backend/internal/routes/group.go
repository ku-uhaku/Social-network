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
}
