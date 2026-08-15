// Package provider defines the platform-agnostic contract every social
// platform integration implements (TikTok, Instagram, X, ...) and a registry
// to look them up by identifier. Adding a new platform means implementing
// this interface in its own subpackage and registering one instance — no
// other code in the app needs to change.
package provider

import (
	"context"
	"time"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	Scope        string
}

type AccountInfo struct {
	PlatformUserID string
	Username       string
	DisplayName    string
	AvatarURL      string
	// Extra carries platform-specific data the frontend needs before posting,
	// e.g. TikTok's privacy_level_options and max_video_post_duration_sec.
	Extra map[string]any
}

type MediaKind string

const (
	MediaImage MediaKind = "image"
	MediaVideo MediaKind = "video"
)

type Content struct {
	Caption   string
	MediaKind MediaKind
	MediaURLs []string
	// MediaDurationSec is the video's duration, if known (populated from the
	// media asset's stored metadata), so providers can enforce their own
	// max-duration rules server-side rather than trusting the client.
	MediaDurationSec *int
	// MediaContentTypes, MediaWidths, and MediaHeights are parallel to
	// MediaURLs (populated from each URL's stored media_assets row, where
	// known) so providers can reject files their own platform would bounce
	// — wrong format, oversized dimensions — before Publish ever calls out.
	// A zero value means unknown, not "no restriction": providers should
	// treat 0 as unverifiable rather than as a passing check.
	MediaContentTypes []string
	MediaWidths       []int
	MediaHeights      []int
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PublishRef struct {
	// ID is the provider's own tracking handle for this publish attempt
	// (e.g. TikTok's publish_id), used to poll for status later.
	ID string
}

type Status struct {
	State   string // "processing" | "published" | "failed"
	Message string
}

type Provider interface {
	Identifier() string
	UsesPKCE() bool

	// AuthURL and Exchange use the redirect URI the provider was
	// constructed with; codeVerifier is only meaningful when UsesPKCE().
	AuthURL(state, codeVerifier string) string
	Exchange(ctx context.Context, code, codeVerifier string) (Token, AccountInfo, error)
	RefreshToken(ctx context.Context, token Token) (Token, error)
	AccountInfo(ctx context.Context, token Token) (AccountInfo, error)

	// ValidateSettings is the authoritative, server-side check — it must
	// reject anything the platform's own live rules (returned in info.Extra)
	// wouldn't allow, not just what the UI happened to grey out.
	ValidateSettings(content Content, settings map[string]any, info AccountInfo) []ValidationError
	Publish(ctx context.Context, token Token, content Content, settings map[string]any) (PublishRef, error)
	PollStatus(ctx context.Context, token Token, ref PublishRef) (Status, error)
}
