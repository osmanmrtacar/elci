package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/osmanmertacar/elci/backend/internal/auth"
)

type contextKey int

const userIDKey contextKey = 0

// RequireAuth rejects the request unless it carries a valid bearer JWT.
func RequireAuth(tokens auth.TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := parseBearer(r, tokens)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, userID)))
		})
	}
}

// OptionalAuth attaches the user id if a valid bearer JWT is present, but
// never rejects the request. Used by the "start OAuth" endpoint, which
// works both for a fresh sign-in and for an already-logged-in user
// connecting an additional platform. Since that endpoint is reached by a
// top-level browser redirect (no custom headers possible), it also accepts
// the token as a "?token=" query parameter as a fallback.
func OptionalAuth(tokens auth.TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := parseBearer(r, tokens)
			if !ok {
				if t := r.URL.Query().Get("token"); t != "" {
					if uid, err := tokens.Parse(t); err == nil {
						userID, ok = uid, true
					}
				}
			}
			if ok {
				r = r.WithContext(context.WithValue(r.Context(), userIDKey, userID))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserID(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}

func parseBearer(r *http.Request, tokens auth.TokenIssuer) (int64, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return 0, false
	}
	userID, err := tokens.Parse(strings.TrimPrefix(header, prefix))
	if err != nil {
		return 0, false
	}
	return userID, true
}
