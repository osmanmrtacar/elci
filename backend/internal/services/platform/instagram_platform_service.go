package platform

import (
	"fmt"

	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
)

// InstagramPlatformService implements PlatformService for Instagram
type InstagramPlatformService struct {
	authService  *services.InstagramAuthService
	mediaService *services.InstagramMediaService
	postService  *services.InstagramPostService
}

// NewInstagramPlatformService creates a new Instagram platform service
func NewInstagramPlatformService(appID, appSecret, redirectURI string) *InstagramPlatformService {
	authService := services.NewInstagramAuthService(appID, appSecret, redirectURI)
	mediaService := services.NewInstagramMediaService()
	postService := services.NewInstagramPostService(authService, mediaService)

	return &InstagramPlatformService{
		authService:  authService,
		mediaService: mediaService,
		postService:  postService,
	}
}

// GetPlatformName returns the platform name
func (s *InstagramPlatformService) GetPlatformName() models.Platform {
	return models.PlatformInstagram
}

// GetRequiredScopes returns the required OAuth scopes
func (s *InstagramPlatformService) GetRequiredScopes() []string {
	return []string{
		"business_management",
		"pages_read_engagement",
		"pages_show_list",
		"instagram_basic",
		"instagram_content_publish",
	}
}

// GenerateAuthURL generates the Instagram OAuth authorization URL (via Facebook)
func (s *InstagramPlatformService) GenerateAuthURL() (AuthURLResponse, error) {
	authResp, err := s.authService.GenerateAuthURL()
	if err != nil {
		return AuthURLResponse{}, fmt.Errorf("failed to generate auth URL: %w", err)
	}

	return AuthURLResponse{
		URL:          authResp.URL,
		State:        authResp.State,
		CodeVerifier: "", // Instagram doesn't use PKCE
	}, nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens
// This automatically exchanges the short-lived token for a long-lived token
func (s *InstagramPlatformService) ExchangeCodeForTokens(code string, additionalParams map[string]string) (*TokenResponse, error) {
	// Step 1: Exchange code for short-lived token (valid for 1 hour)
	tokenResp, err := s.authService.ExchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Step 2: Exchange short-lived token for long-lived token (valid for 60 days)
	longLivedResp, err := s.authService.ExchangeLongLivedToken(tokenResp.AccessToken)
	if err != nil {
		// If long-lived exchange fails, fall back to short-lived token
		fmt.Printf("WARNING: Failed to get long-lived token, using short-lived: %v\n", err)
		expiresIn := tokenResp.ExpiresIn
		if expiresIn == 0 {
			expiresIn = 3600 // 1 hour for short-lived token
		}
		return &TokenResponse{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.AccessToken, // Use same token as refresh token for retry
			ExpiresIn:    expiresIn,
			TokenType:    "Bearer",
			Scope:        "instagram_business_basic,instagram_business_content_publish",
		}, nil
	}

	// Use long-lived token expiry (typically 60 days = 5184000 seconds)
	expiresIn := longLivedResp.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 60 * 24 * 60 * 60 // 60 days in seconds
	}

	return &TokenResponse{
		AccessToken:  longLivedResp.AccessToken,
		RefreshToken: longLivedResp.AccessToken, // Instagram uses the same token for refresh
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		Scope:        "instagram_business_basic,instagram_business_content_publish",
	}, nil
}

// RefreshAccessToken refreshes an Instagram long-lived access token
// Instagram long-lived tokens can be refreshed as long as they haven't expired
// The refreshed token will be valid for another 60 days
func (s *InstagramPlatformService) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	// Refresh the long-lived token
	refreshedToken, err := s.authService.RefreshLongLivedToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh Instagram token: %w", err)
	}

	expiresIn := refreshedToken.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 60 * 24 * 60 * 60 // 60 days in seconds
	}

	return &TokenResponse{
		AccessToken:  refreshedToken.AccessToken,
		RefreshToken: refreshedToken.AccessToken, // Instagram uses the same token for refresh
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		Scope:        "instagram_business_basic,instagram_business_content_publish",
	}, nil
}

// GetUserInfo retrieves user information from Instagram
func (s *InstagramPlatformService) GetUserInfo(accessToken string) (*UserInfo, error) {
	userInfo, err := s.authService.GetInstagramUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get Instagram user info: %w", err)
	}

	return &UserInfo{
		PlatformUserID: userInfo.ID,
		Username:       userInfo.Username,
		DisplayName:    userInfo.Name,
		AvatarURL:      "", // Instagram Business API doesn't provide profile picture
		Email:          "", // Not available through Instagram Business API
	}, nil
}

// UploadMedia uploads media to Instagram (creates container and waits for processing)
func (s *InstagramPlatformService) UploadMedia(accessToken string, mediaURL string) (string, error) {
	// Get Instagram user info
	userInfo, err := s.authService.GetInstagramUserInfo(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get Instagram user info: %w", err)
	}

	// Create media container (Step 1)
	containerID, err := s.mediaService.CreateMediaContainer(
		accessToken,
		userInfo.ID,
		mediaURL,
		"", // No caption during upload
		"REELS",
	)
	if err != nil {
		return "", fmt.Errorf("failed to create media container: %w", err)
	}

	// Wait for processing (Step 2)
	success, err := s.mediaService.WaitForMediaProcessing(accessToken, containerID, 300)
	if err != nil {
		return "", fmt.Errorf("media processing failed: %w", err)
	}
	if !success {
		return "", fmt.Errorf("media processing did not complete successfully")
	}

	// Return container ID (will be used for publishing)
	return containerID, nil
}

// CreatePost creates and publishes a post to Instagram
func (s *InstagramPlatformService) CreatePost(accessToken string, content PostContent) (*PostResponse, error) {
	mediaID, permalink, err := s.postService.CreatePost(accessToken, content.MediaURL, content.Text)
	if err != nil {
		return &PostResponse{
			Status:   "failed",
			ErrorMsg: err.Error(),
		}, err
	}

	return &PostResponse{
		PostID:    mediaID,
		PublishID: "",
		Status:    "published",
		ShareURL:  permalink,
		ErrorMsg:  "",
	}, nil
}

// GetPostStatus retrieves the status of a post
func (s *InstagramPlatformService) GetPostStatus(accessToken string, postID string) (*PostStatusResponse, error) {
	// Instagram doesn't provide detailed status after publishing
	// Once published, the post is live
	return &PostStatusResponse{
		Status:          "published",
		PostID:          postID,
		ShareURL:        "", // Would need to fetch permalink
		FailReason:      "",
		ProgressPercent: 100,
	}, nil
}
