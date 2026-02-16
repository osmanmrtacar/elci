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

	// Convert platform TikTokSettings to services TikTokSettings
	var tiktokSettings *services.TikTokPostSettings
	if content.TikTokSettings != nil {
		// Validate TikTok settings according to UX Guidelines
		settings := content.TikTokSettings

		// Privacy level is required (Point 2b)
		if settings.PrivacyLevel == "" {
			return nil, fmt.Errorf("privacy level is required for TikTok posts")
		}

		// Validate privacy level is one of the allowed values
		validPrivacyLevels := map[string]bool{
			"PUBLIC_TO_EVERYONE":    true,
			"MUTUAL_FOLLOW_FRIENDS": true,
			"FOLLOWER_OF_CREATOR":   true,
			"SELF_ONLY":             true,
		}
		if !validPrivacyLevels[settings.PrivacyLevel] {
			return nil, fmt.Errorf("invalid privacy level: %s", settings.PrivacyLevel)
		}

		// Branded content cannot be private (Point 3b)
		if (settings.IsBrandContent || settings.IsBrandOrganic) && settings.PrivacyLevel == "SELF_ONLY" {
			return nil, fmt.Errorf("branded content cannot be set to private visibility")
		}

		// Title max length is 150 characters (Point 2a)
		if len(settings.Title) > 150 {
			return nil, fmt.Errorf("title cannot exceed 150 characters")
		}

		// If DirectPost is enabled, brand_content_toggle must be set (Point 3)
		isDirectPost := settings.DirectPost
		if isDirectPost && !settings.IsBrandContent && !settings.IsBrandOrganic {
			// This is handled by frontend, but we ensure it's false in API
			settings.IsBrandContent = false
			settings.IsBrandOrganic = false
		}

		tiktokSettings = &services.TikTokPostSettings{
			Title:          settings.Title,
			PrivacyLevel:   settings.PrivacyLevel,
			AllowComment:   settings.AllowComment,
			AllowDuet:      settings.AllowDuet,
			AllowStitch:    settings.AllowStitch,
			IsBrandContent: settings.IsBrandContent,
			IsBrandOrganic: settings.IsBrandOrganic,
			AutoAddMusic:   settings.AutoAddMusic,
			DirectPost:     settings.DirectPost,
		}
	} else {
		// If no settings provided, return error (privacy level is required)
		return nil, fmt.Errorf("TikTok settings are required including privacy level")
	}

	var resp *services.PublishVideoResponse
	var err error

	// Detect media type and publish accordingly
	if services.IsImageURL(mediaURL) {
		// Photo post - TikTok accepts array of image URLs
		// Use MediaURLs if provided, otherwise use single mediaURL
		imageURLs := content.MediaURLs
		if len(imageURLs) == 0 {
			imageURLs = []string{mediaURL}
		}
		resp, err = s.tiktokService.PublishPhotoFromURL(accessToken, imageURLs, content.Text, tiktokSettings)
	} else {
		// Video post - TikTok only supports single video
		resp, err = s.tiktokService.PublishVideoFromURL(accessToken, mediaURL, content.Text, tiktokSettings)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Inbox posts are done once TikTok accepts them; Direct Post needs polling
	status := string(models.PostStatusProcessing)
	if content.TikTokSettings != nil && !content.TikTokSettings.DirectPost {
		status = string(models.PostStatusSentToInbox)
	}

	return &PostResponse{
		PostID: resp.Data.PublishID,
		Status: status,
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
		status = string(models.PostStatusPublished)
		if resp.Data.ShareID != "" {
			shareURL = fmt.Sprintf("https://www.tiktok.com/@user/video/%s", resp.Data.ShareID)
		}
	case "SEND_TO_USER_INBOX":
		status = string(models.PostStatusSentToInbox)
	case "FAILED":
		status = string(models.PostStatusFailed)
	case "PROCESSING_UPLOAD", "PROCESSING_DOWNLOAD":
		status = string(models.PostStatusProcessing)
	default:
		status = string(models.PostStatusProcessing)
	}

	return &PostStatusResponse{
		Status:     status,
		PostID:     resp.Data.PublishID,
		ShareURL:   shareURL,
		FailReason: resp.Data.FailReason,
	}, nil
}
