package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"kuu/internal/helper"
	"kuu/internal/models"
)

type ContextKey string

const UserContextKey ContextKey = "current_user"

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.User)

	return user, ok
}

// RequireAuth validates the session token and injects the user into the request context
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract session token from cookies (Matches your login cookie name)
		cookie, err := r.Cookie("session_token")
		if err != nil {
			helper.Error(w, http.StatusUnauthorized, "Authentication required: missing token")
			return
		}

		// 2. Validate session and fetch user in one database step via Service layer
		user, err := m.Service.ValidateSession(cookie.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// If token is invalid or expired, clear the cookie from their browser
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				helper.Error(w, http.StatusUnauthorized, "Authentication required: invalid or expired session")
				return
			}
			helper.Error(w, http.StatusInternalServerError, "Internal error verifying session")
			return
		}

		// 3. Inject the user into the request context and proceed down the chain
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
