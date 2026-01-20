package platform

import (
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
)

// PlatformService defines the interface that all platform services must implement
// This allows us to support multiple social media platforms (TikTok, X, Instagram, etc.)
// in a uniform way without code duplication
type PlatformService interface {
	// Auth methods
	GenerateAuthURL() (AuthURLResponse, error)
	ExchangeCodeForTokens(code string, additionalParams map[string]string) (*TokenResponse, error)
	RefreshAccessToken(refreshToken string) (*TokenResponse, error)
	GetUserInfo(accessToken string) (*UserInfo, error)

	// Media methods
	UploadMedia(accessToken string, mediaURL string) (string, error)

	// Post methods
	CreatePost(accessToken string, content PostContent) (*PostResponse, error)
	GetPostStatus(accessToken string, postID string) (*PostStatusResponse, error)

	// Metadata
	GetPlatformName() models.Platform
	GetRequiredScopes() []string
}

// AuthURLResponse contains the OAuth authorization URL and associated data
type AuthURLResponse struct {
	URL          string // The OAuth authorization URL to redirect user to
	State        string // CSRF protection state parameter
	CodeVerifier string // PKCE code verifier (for platforms like X that use PKCE)
}

// TokenResponse contains OAuth tokens received from the platform
type TokenResponse struct {
	AccessToken  string // Access token for API requests
	RefreshToken string // Refresh token for obtaining new access tokens
	ExpiresIn    int    // Token expiration time in seconds
	TokenType    string // Token type (usually "Bearer")
	Scope        string // Granted scopes
}

// UserInfo contains basic user information from the platform
type UserInfo struct {
	PlatformUserID string // Platform-specific user ID
	Username       string // Username/handle
	DisplayName    string // Display name
	AvatarURL      string // Profile picture URL
	Email          string // Email (if available)
}

// TikTokSettings represents TikTok-specific post settings (required by TikTok UX Guidelines)
type TikTokSettings struct {
	Title          string // Video title
	PrivacyLevel   string // PUBLIC_TO_EVERYONE, MUTUAL_FOLLOW_FRIENDS, FOLLOWER_OF_CREATOR, SELF_ONLY
	AllowComment   bool   // Allow comments (default: false per UX guidelines)
	AllowDuet      bool   // Allow duet (default: false per UX guidelines)
	AllowStitch    bool   // Allow stitch (default: false per UX guidelines)
	IsBrandContent bool   // Promoting own brand
	IsBrandOrganic bool   // Paid partnership (branded content)
	AutoAddMusic   bool   // Auto-add trending music to photo posts (only for photos)
}

// PostContent represents the content to be posted
type PostContent struct {
	Text           string          // Post text/caption
	MediaURL       string          // Primary URL of media to download and upload
	MediaURLs      []string        // Multiple media URLs (for carousel/multi-image)
	MediaIDs       []string        // Pre-uploaded media IDs (for platforms like X)
	TikTokSettings *TikTokSettings // TikTok-specific settings (optional)
}

// PostResponse contains the result of creating a post
type PostResponse struct {
	PostID    string // Platform-specific post ID (immediate for X, after processing for TikTok)
	PublishID string // Async publish ID (for TikTok)
	Status    string // Post status (pending, processing, published)
	ShareURL  string // URL to view the post on the platform
	ErrorMsg  string // Error message if post creation failed
}

// PostStatusResponse contains the current status of a post
type PostStatusResponse struct {
	Status          string // pending, processing, published, failed
	PostID          string // Platform-specific post ID
	ShareURL        string // URL to view the post on the platform
	FailReason      string // Reason for failure if status is failed
	ProgressPercent int    // Progress percentage (for video processing)
}
