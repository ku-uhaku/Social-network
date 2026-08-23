package routes

import (
	"net/http"
	"time"

	"kuu/internal/handler"
	"kuu/internal/helper"
	"kuu/internal/middleware"
)

func registerAuthRoutes(mux *http.ServeMux, h *handler.Handler, m *middleware.Middleware) {
	ratelimeer:=helper.Neewratelimeter(time.Minute)
	// mux.Handle("/profile/", ratelimeer.Wraponall("api", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

	// })))

	mux.Handle(
		"/api/v1/auth/register",
		m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Register)),
	)
	// mux.Handle(
	// 	"/api/v1/auth/login",
	// 	m.AllowMethods(http.MethodPost)(http.HandlerFunc(h.Login)),
	// )
	mux.Handle("/api/v1/auth/login",ratelimeer.Wraponall("authonti",m.AllowMethods(http.MethodPost)(
				http.HandlerFunc(h.Login),
			),
		),
	)
	// mux.Handle("/api/v1/auth/me",ratelimeer.Wraponall("api",m.AllowMethods(http.MethodGet)(
	// 			http.HandlerFunc(h.Me),
	// 		),
	// 	),
	// )
	mux.Handle("/api/v1/auth/logout",ratelimeer.Wraponall("authonti",m.AllowMethods(http.MethodPost)(
				http.HandlerFunc(h.Logout),
			),
		),
	)
	mux.Handle(
		"/api/v1/auth/me",
		m.AllowMethods(http.MethodGet)(m.RequireAuth(http.HandlerFunc(h.Me))),
	)
	// mux.Handle(
	// 	"/api/v1/auth/logout",
	// 	m.AllowMethods(http.MethodPost)(m.RequireAuth(http.HandlerFunc(h.Logout))),
	// )
}
