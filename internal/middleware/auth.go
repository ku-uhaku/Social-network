package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

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
		// 1. Extract session token from cookies
		cookie, err := r.Cookie("kuu_token")
		if err != nil {
			helper.Error(w, http.StatusUnauthorized, "Authentication required: missing token")
			return
		}

		// 2. Fetch session details from the repository
		session, err := m.Repo.GetSession(cookie.Value) // Assumes your middleware has access to h.Repo
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helper.Error(w, http.StatusUnauthorized, "Authentication required: invalid session")
				return
			}
			helper.Error(w, http.StatusInternalServerError, "Database error verifying session")
			return
		}

		// 3. Confirm that the session hasn't expired yet
		if time.Now().UTC().After(session.ExpiresAt) {
			helper.Error(w, http.StatusUnauthorized, "Authentication required: session expired")
			return
		}

		// 4. Retrieve the actual user profile attached to this session
		user, err := m.Repo.GetUserByID(session.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helper.Error(w, http.StatusUnauthorized, "Authentication required: user no longer exists")
				return
			}
			helper.Error(w, http.StatusInternalServerError, "Database error retrieving user profile")
			return
		}

		// 5. Inject the user into the request context and proceed down the chain
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
