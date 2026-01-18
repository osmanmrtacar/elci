package platform

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/osmanmertacar/sosyal/backend/internal/config"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
)

// TikTokPlatformService implements PlatformService for TikTok
type TikTokPlatformService struct {
	tiktokService *services.TikTokService
	config        *config.Config
}

// NewTikTokPlatformService creates a new TikTok platform service
func NewTikTokPlatformService(cfg *config.Config, tiktokService *services.TikTokService) *TikTokPlatformService {
	return &TikTokPlatformService{
		tiktokService: tiktokService,
		config:        cfg,
	}
}

// GetPlatformName returns the platform identifier
func (s *TikTokPlatformService) GetPlatformName() models.Platform {
	return models.PlatformTikTok
}

// GetRequiredScopes returns the OAuth scopes required by TikTok
func (s *TikTokPlatformService) GetRequiredScopes() []string {
	return s.config.TikTok.Scopes
}

// GenerateAuthURL generates the OAuth authorization URL
// Note: TikTok doesn't use PKCE, so CodeVerifier will be empty
func (s *TikTokPlatformService) GenerateAuthURL() (AuthURLResponse, error) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		return AuthURLResponse{}, fmt.Errorf("failed to generate state: %w", err)
	}

	// Build OAuth URL
	baseURL := "https://www.tiktok.com/v2/auth/authorize/"
	params := url.Values{}
	params.Add("client_key", s.config.TikTok.ClientKey)
	params.Add("scope", strings.Join(s.config.TikTok.Scopes, ","))
	params.Add("response_type", "code")
	params.Add("redirect_uri", s.config.TikTok.RedirectURI)
	params.Add("state", state)

	authURL := baseURL + "?" + params.Encode()

	return AuthURLResponse{
		URL:          authURL,
		State:        state,
		CodeVerifier: "", // TikTok doesn't use PKCE
	}, nil
}

// generateRandomState generates a random state parameter for OAuth
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ExchangeCodeForTokens exchanges authorization code for access token
// additionalParams is not used for TikTok (no PKCE)
func (s *TikTokPlatformService) ExchangeCodeForTokens(code string, additionalParams map[string]string) (*TokenResponse, error) {
	resp, err := s.tiktokService.ExchangeCodeForTokens(code)
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
func (s *TikTokPlatformService) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	resp, err := s.tiktokService.RefreshAccessToken(refreshToken)
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

// GetUserInfo fetches user information from TikTok
func (s *TikTokPlatformService) GetUserInfo(accessToken string) (*UserInfo, error) {
	resp, err := s.tiktokService.GetUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// TikTok v2 API doesn't provide username field, use DisplayName or OpenID
	username := resp.DisplayName
	if username == "" {
		username = resp.OpenID
	}

	return &UserInfo{
		PlatformUserID: resp.OpenID,
		Username:       username,
		DisplayName:    resp.DisplayName,
		AvatarURL:      resp.AvatarURL,
		Email:          "", // TikTok doesn't provide email
	}, nil
}

// UploadMedia is not applicable for TikTok (videos are published directly from URL)
// This method returns the mediaURL as-is for validation
func (s *TikTokPlatformService) UploadMedia(accessToken, mediaURL string) (string, error) {
	// TikTok publishes directly from URL, no separate upload step needed
	// Just return the URL for use in CreatePost
	if mediaURL == "" {
		return "", fmt.Errorf("media URL is required for TikTok")
	}
	return mediaURL, nil
}

// CreatePost publishes a video or photo to TikTok
// Automatically detects media type from URL and uses appropriate method
func (s *TikTokPlatformService) CreatePost(accessToken string, content PostContent) (*PostResponse, error) {
	mediaURL := content.MediaURL
	if mediaURL == "" && len(content.MediaIDs) > 0 {
		// If MediaIDs provided (from UploadMedia), use the first one as URL
		mediaURL = content.MediaIDs[0]
	}

	if mediaURL == "" {
		return nil, fmt.Errorf("media URL is required for TikTok posts")
	}

	var resp *services.PublishVideoResponse
	var err error

	// Detect media type and publish accordingly
	if services.IsImageURL(mediaURL) {
		// Photo post - TikTok accepts array of image URLs
		resp, err = s.tiktokService.PublishPhotoFromURL(accessToken, []string{mediaURL}, content.Text)
	} else {
		// Video post
		resp, err = s.tiktokService.PublishVideoFromURL(accessToken, mediaURL, content.Text)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &PostResponse{
		PostID: resp.Data.PublishID,
		Status: "processing", // TikTok posts are processed asynchronously
	}, nil
}

// GetPostStatus gets the status of a TikTok post
func (s *TikTokPlatformService) GetPostStatus(accessToken, postID string) (*PostStatusResponse, error) {
	resp, err := s.tiktokService.GetPublishStatus(accessToken, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post status: %w", err)
	}

	// Map TikTok status to our status format
	// TikTok statuses: PUBLISH_COMPLETE, PROCESSING_UPLOAD, PROCESSING_DOWNLOAD, FAILED, etc.
	status := "processing"
	shareURL := ""

	switch resp.Data.Status {
	case "PUBLISH_COMPLETE":
		status = "published"
		if resp.Data.ShareID != "" {
			shareURL = fmt.Sprintf("https://www.tiktok.com/@user/video/%s", resp.Data.ShareID)
		}
	case "FAILED":
		status = "failed"
	case "PROCESSING_UPLOAD", "PROCESSING_DOWNLOAD", "SEND_TO_USER_INBOX":
		status = "processing"
	default:
		status = "processing"
	}

	return &PostStatusResponse{
		Status:     status,
		PostID:     resp.Data.PublishID,
		ShareURL:   shareURL,
		FailReason: resp.Data.FailReason,
	}, nil
}
