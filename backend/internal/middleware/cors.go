package middleware

import "net/http"

// CORS intercepts requests to handle cross-origin constraints safely
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Explicitly allow your frontend server domain (No wildcards allowed with credentials)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // TODO: won't work with docker
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}
