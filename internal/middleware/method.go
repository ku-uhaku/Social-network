package middleware

import (
	"net/http"
)

// AllowMethods ensures the request uses one of the specified HTTP methods
func (m *Middleware) AllowMethods(allowedMethods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAllowed := false
			for _, method := range allowedMethods {
				if r.Method == method {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
