package domain

import "time"

type Platform string

const (
	PlatformTikTok    Platform = "tiktok"
	PlatformInstagram Platform = "instagram"
	PlatformX         Platform = "x"
)

// PlatformConnection is one user's authorization to post to a specific
// platform account. Access/refresh tokens live here since they're always
// one-to-one with a connection.
type PlatformConnection struct {
	ID             int64
	UserID         int64
	Platform       Platform
	PlatformUserID string
	Username       string
	DisplayName    string
	AvatarURL      string
	AccessToken    string
	RefreshToken   string
	TokenExpiresAt time.Time
	Scope          string
	IsActive       bool
	ConnectedAt    time.Time
	LastUsedAt     time.Time
}
