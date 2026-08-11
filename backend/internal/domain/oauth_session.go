package domain

import "time"

// OAuthSession tracks a single in-flight OAuth attempt so the callback can be
// matched back to its request (CSRF state) and, for PKCE platforms, the
// original code verifier. UserID is set when an already-logged-in user is
// connecting an additional platform, and nil for a fresh sign-in.
type OAuthSession struct {
	ID           int64
	State        string
	CodeVerifier string
	Platform     Platform
	UserID       *int64
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
