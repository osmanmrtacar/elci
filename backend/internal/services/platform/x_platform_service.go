package platform

import (
	"fmt"

	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
)

// XPlatformService implements PlatformService for X (Twitter)
type XPlatformService struct {
	authService  *services.XAuthService
	mediaService *services.XMediaService
	postService  *services.XPostService
	clientID     string
	clientSecret string
	redirectURI  string
}

// NewXPlatformService creates a new X platform service
func NewXPlatformService(clientID, clientSecret, redirectURI string) *XPlatformService {
	return &XPlatformService{
		authService:  services.NewXAuthService(clientID, clientSecret, redirectURI),
		mediaService: services.NewXMediaService(),
		postService:  services.NewXPostService(),
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetPlatformName returns the platform identifier
func (s *XPlatformService) GetPlatformName() models.Platform {
	return models.PlatformX
}

// GetRequiredScopes returns the OAuth scopes required by X
func (s *XPlatformService) GetRequiredScopes() []string {
	return []string{"tweet.read", "tweet.write", "users.read", "offline.access", "media.write"}
}

// GenerateAuthURL generates the OAuth authorization URL with PKCE
func (s *XPlatformService) GenerateAuthURL() (AuthURLResponse, error) {
	url, state, codeVerifier, err := s.authService.GenerateAuthURL()
	if err != nil {
		return AuthURLResponse{}, fmt.Errorf("failed to generate auth URL: %w", err)
	}

	return AuthURLResponse{
		URL:          url,
		State:        state,
		CodeVerifier: codeVerifier,
	}, nil
}

// ExchangeCodeForTokens exchanges authorization code for access token
// additionalParams must contain "code_verifier" for PKCE
func (s *XPlatformService) ExchangeCodeForTokens(code string, additionalParams map[string]string) (*TokenResponse, error) {
	codeVerifier, ok := additionalParams["code_verifier"]
	if !ok {
		return nil, fmt.Errorf("code_verifier is required for X OAuth")
	}

	resp, err := s.authService.ExchangeCodeForToken(code, codeVerifier)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	return &TokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		Scope:        resp.Scope,
	}, nil
}

// RefreshAccessToken refreshes an expired access token
func (s *XPlatformService) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	resp, err := s.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh access token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		Scope:        resp.Scope,
	}, nil
}

// GetUserInfo fetches user information from X
func (s *XPlatformService) GetUserInfo(accessToken string) (*UserInfo, error) {
	resp, err := s.authService.GetUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &UserInfo{
		PlatformUserID: resp.Data.ID,
		Username:       resp.Data.Username,
		DisplayName:    resp.Data.Name,
		AvatarURL:      resp.Data.ProfileImageURL,
		Email:          "", // X doesn't provide email in this endpoint
	}, nil
}

// UploadMedia downloads and uploads media to X
func (s *XPlatformService) UploadMedia(accessToken, mediaURL string) (string, error) {
	mediaID, err := s.mediaService.UploadFromURL(accessToken, mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to upload media: %w", err)
	}
	return mediaID, nil
}

// CreatePost creates a tweet with optional media
func (s *XPlatformService) CreatePost(accessToken string, content PostContent) (*PostResponse, error) {
	var mediaIDs []string

	// If media URL provided, upload it first
	if content.MediaURL != "" {
		mediaID, err := s.UploadMedia(accessToken, content.MediaURL)
		if err != nil {
			return nil, err
		}
		mediaIDs = []string{mediaID}
	} else if len(content.MediaIDs) > 0 {
		// Use pre-uploaded media IDs
		mediaIDs = content.MediaIDs
	}

	// Create tweet
	resp, err := s.postService.CreatePost(accessToken, content.Text, mediaIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &PostResponse{
		PostID:   resp.Data.ID,
		Status:   "published", // X posts are published immediately
		ShareURL: fmt.Sprintf("https://twitter.com/i/web/status/%s", resp.Data.ID),
	}, nil
}

// GetPostStatus gets the status of a post (X posts are immediate, so always returns published)
func (s *XPlatformService) GetPostStatus(accessToken, postID string) (*PostStatusResponse, error) {
	// X posts are published immediately, no async processing
	return &PostStatusResponse{
		Status: "published",
		PostID: postID,
	}, nil
}
