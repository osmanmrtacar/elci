package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/osmanmertacar/sosyal/backend/internal/config"
)

const (
	tiktokTokenURL       = "https://open.tiktokapis.com/v2/oauth/token/"
	tiktokUserInfoURL    = "https://open.tiktokapis.com/v2/user/info/"
	tiktokPublishURL     = "https://open.tiktokapis.com/v2/post/publish/video/init/"
	tiktokInboxURL       = "https://open.tiktokapis.com/v2/post/publish/inbox/video/init/"
	tiktokPublishStatus  = "https://open.tiktokapis.com/v2/post/publish/status/fetch/"
	tiktokContentInitURL  = "https://open.tiktokapis.com/v2/post/publish/content/init/"
	tiktokCreatorInfoURL  = "https://open.tiktokapis.com/v2/post/publish/creator_info/query/"
)

type TikTokService struct {
	config     *config.Config
	httpClient *http.Client
}

// TikTokUserInfo represents user information from TikTok
type TikTokUserInfo struct {
	OpenID      string `json:"open_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// TikTokUserInfoResponse represents the response from TikTok's user info endpoint
type TikTokUserInfoResponse struct {
	Data struct {
		User TikTokUserInfo `json:"user"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// NewTikTokService creates a new TikTok service
func NewTikTokService(cfg *config.Config) *TikTokService {
	return &TikTokService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExchangeCodeForTokens exchanges an authorization code for access and refresh tokens
func (s *TikTokService) ExchangeCodeForTokens(code string) (*TikTokTokenResponse, error) {
	// TikTok requires application/x-www-form-urlencoded
	formData := url.Values{}
	formData.Set("client_key", s.config.TikTok.ClientKey)
	formData.Set("client_secret", s.config.TikTok.ClientSecret)
	formData.Set("code", code)
	formData.Set("grant_type", "authorization_code")
	formData.Set("redirect_uri", s.config.TikTok.RedirectURI)

	req, err := http.NewRequest("POST", tiktokTokenURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	// DEBUG: Log the raw response
	fmt.Printf("DEBUG: Token exchange raw response: %s\n", string(responseBody))

	var tokenResponse TikTokTokenResponse
	if err := json.Unmarshal(responseBody, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// DEBUG: Log parsed token response
	fmt.Printf("DEBUG: Parsed token - AccessToken length: %d, RefreshToken length: %d\n",
		len(tokenResponse.AccessToken), len(tokenResponse.RefreshToken))

	return &tokenResponse, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *TikTokService) RefreshAccessToken(refreshToken string) (*TikTokTokenResponse, error) {
	// TikTok requires application/x-www-form-urlencoded
	formData := url.Values{}
	formData.Set("client_key", s.config.TikTok.ClientKey)
	formData.Set("client_secret", s.config.TikTok.ClientSecret)
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tiktokTokenURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	var tokenResponse TikTokTokenResponse
	if err := json.Unmarshal(responseBody, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tokenResponse, nil
}

// GetUserInfo fetches user information from TikTok
func (s *TikTokService) GetUserInfo(accessToken string) (*TikTokUserInfo, error) {
	// Use query parameters for fields (username is not available in v2 API)
	url := tiktokUserInfoURL + "?fields=open_id,display_name,avatar_url"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// TikTok API v2 requires "Bearer " prefix with space
	authHeader := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	// Log request details for debugging
	tokenPreview := accessToken
	if len(tokenPreview) > 10 {
		tokenPreview = tokenPreview[:10] + "..."
	}
	fmt.Printf("DEBUG: Making GET request to: %s\n", url)
	fmt.Printf("DEBUG: Authorization header: Bearer %s\n", tokenPreview)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// DEBUG: Log the full response
	fmt.Printf("DEBUG: User info response status: %d\n", resp.StatusCode)
	fmt.Printf("DEBUG: User info response body: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	var userInfoResponse TikTokUserInfoResponse
	if err := json.Unmarshal(responseBody, &userInfoResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w - body: %s", err, string(responseBody))
	}

	// Check for API errors (TikTok uses "ok" to indicate success)
	if userInfoResponse.Error.Code != "" && userInfoResponse.Error.Code != "ok" {
		return nil, fmt.Errorf("TikTok API error: %s - %s", userInfoResponse.Error.Code, userInfoResponse.Error.Message)
	}

	return &userInfoResponse.Data.User, nil
}

// TikTokPostSettings represents settings for TikTok posts (required by TikTok UX Guidelines)
type TikTokPostSettings struct {
	Title          string // Video title
	PrivacyLevel   string // PUBLIC_TO_EVERYONE, MUTUAL_FOLLOW_FRIENDS, FOLLOWER_OF_CREATOR, SELF_ONLY
	AllowComment   bool   // Allow comments
	AllowDuet      bool   // Allow duet
	AllowStitch    bool   // Allow stitch
	IsBrandContent bool   // Promoting own brand
	IsBrandOrganic bool   // Paid partnership (branded content)
	AutoAddMusic   bool   // Auto-add trending music to photo posts (only for photos)
	DirectPost     bool   // Direct Post (true) vs Send to Inbox (false)
}

// PublishVideoRequest represents the request to publish a video
type PublishVideoRequest struct {
	VideoURL string
	Caption  string
}

// PublishVideoResponse represents the response from TikTok's publish endpoint
type PublishVideoResponse struct {
	Data struct {
		PublishID string `json:"publish_id"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// PublishStatusResponse represents the response from TikTok's status endpoint
type PublishStatusResponse struct {
	Data struct {
		Status       string `json:"status"`
		PublishID    string `json:"publish_id"`
		ShareID      string `json:"share_id"`
		FailReason   string `json:"fail_reason"`
		PrivacyLevel string `json:"privacy_level"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// PublishVideoFromURL publishes a video to TikTok from a URL
func (s *TikTokService) PublishVideoFromURL(accessToken string, videoURL string, caption string, settings *TikTokPostSettings) (*PublishVideoResponse, error) {
	// Determine whether to use Direct Post or Send to Inbox
	useDirectPost := true
	if settings != nil {
		useDirectPost = settings.DirectPost
	}

	// Build post_info with settings
	postInfo := map[string]interface{}{
		"title": caption, // Use caption as default title
	}

	// Apply TikTok settings if provided (required by TikTok UX Guidelines)
	if settings != nil {
		// Use provided title if available
		if settings.Title != "" {
			postInfo["title"] = settings.Title
		}

		// Privacy level (required - user must select, no default)
		if settings.PrivacyLevel != "" {
			postInfo["privacy_level"] = settings.PrivacyLevel
		}

		// Interaction settings (default: disabled per UX guidelines)
		// TikTok API uses disable_ prefix (true = disabled)
		postInfo["disable_comment"] = !settings.AllowComment
		postInfo["disable_duet"] = !settings.AllowDuet
		postInfo["disable_stitch"] = !settings.AllowStitch

		// Commercial content disclosure (brand_content_toggle is required by Direct Post API only)
		if useDirectPost {
			postInfo["brand_content_toggle"] = settings.IsBrandContent || settings.IsBrandOrganic
			if settings.IsBrandOrganic {
				postInfo["brand_organic_toggle"] = true
			}
		}
	} else if useDirectPost {
		// brand_content_toggle is required by the Direct Post API, default to false
		postInfo["brand_content_toggle"] = false
	}

	requestBody := map[string]interface{}{
		"post_info": postInfo,
		"source_info": map[string]interface{}{
			"source":    "PULL_FROM_URL",
			"video_url": videoURL,
		},
	}

	// Choose endpoint based on publish mode
	publishURL := tiktokPublishURL
	if !useDirectPost {
		publishURL = tiktokInboxURL
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fmt.Printf("DEBUG PublishVideo: POST %s body=%s\n", publishURL, string(body))

	req, err := http.NewRequest("POST", publishURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("DEBUG PublishVideo: status=%d response=%s\n", resp.StatusCode, string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	var publishResponse PublishVideoResponse
	if err := json.Unmarshal(responseBody, &publishResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors (TikTok uses "ok" to indicate success)
	if publishResponse.Error.Code != "" && publishResponse.Error.Code != "ok" {
		return nil, fmt.Errorf("TikTok API error: %s - %s", publishResponse.Error.Code, publishResponse.Error.Message)
	}

	return &publishResponse, nil
}

// GetPublishStatus checks the status of a video publish
func (s *TikTokService) GetPublishStatus(accessToken string, publishID string) (*PublishStatusResponse, error) {
	requestBody := map[string]interface{}{
		"publish_id": publishID,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fmt.Printf("DEBUG GetPublishStatus: POST %s publish_id=%s\n", tiktokPublishStatus, publishID)

	req, err := http.NewRequest("POST", tiktokPublishStatus, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("DEBUG GetPublishStatus: status=%d response=%s\n", resp.StatusCode, string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	var statusResponse PublishStatusResponse
	if err := json.Unmarshal(responseBody, &statusResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors (TikTok uses "ok" to indicate success)
	if statusResponse.Error.Code != "" && statusResponse.Error.Code != "ok" {
		return nil, fmt.Errorf("TikTok API error: %s - %s", statusResponse.Error.Code, statusResponse.Error.Message)
	}

	return &statusResponse, nil
}

// CreatorInfoResponse represents the creator info from TikTok's Creator Info API
type CreatorInfoResponse struct {
	PrivacyLevelOptions     []string `json:"privacy_level_options"`
	MaxVideoPostDurationSec int      `json:"max_video_post_duration_sec"`
	StitchDisabled          bool     `json:"stitch_disabled"`
	CommentDisabled         bool     `json:"comment_disabled"`
	DuetDisabled            bool     `json:"duet_disabled"`
}

// GetCreatorInfo fetches creator posting capabilities from TikTok's Creator Info API
func (s *TikTokService) GetCreatorInfo(accessToken string) (*CreatorInfoResponse, error) {
	body, err := json.Marshal(map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", tiktokCreatorInfoURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("DEBUG GetCreatorInfo: status=%d response=%s\n", resp.StatusCode, string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	var raw struct {
		Data  CreatorInfoResponse `json:"data"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(responseBody, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if raw.Error.Code != "" && raw.Error.Code != "ok" {
		return nil, fmt.Errorf("TikTok API error: %s - %s", raw.Error.Code, raw.Error.Message)
	}

	return &raw.Data, nil
}

// PublishPhotoFromURL publishes a photo post to TikTok from one or more image URLs
func (s *TikTokService) PublishPhotoFromURL(accessToken string, imageURLs []string, caption string, settings *TikTokPostSettings) (*PublishVideoResponse, error) {
	if len(imageURLs) == 0 {
		return nil, fmt.Errorf("at least one image URL is required")
	}

	// Build post_info with settings
	postInfo := map[string]interface{}{
		"title":       caption,
		"description": caption,
	}

	// Apply TikTok settings if provided (required by TikTok UX Guidelines)
	if settings != nil {
		// Use provided title if available
		if settings.Title != "" {
			postInfo["title"] = settings.Title
		}

		// Privacy level (required - user must select, no default)
		if settings.PrivacyLevel != "" {
			postInfo["privacy_level"] = settings.PrivacyLevel
		}

		// Interaction settings (default: disabled per UX guidelines)
		postInfo["disable_comment"] = !settings.AllowComment

		// Commercial content disclosure
		if settings.IsBrandContent || settings.IsBrandOrganic {
			postInfo["brand_content_toggle"] = true
			postInfo["brand_organic_toggle"] = settings.IsBrandOrganic
		}

		// Auto-add music for photo posts (TikTok auto-selects trending music)
		postInfo["auto_add_music"] = settings.AutoAddMusic
	}

	requestBody := map[string]interface{}{
		"post_info": postInfo,
		"source_info": map[string]interface{}{
			"source":            "PULL_FROM_URL",
			"photo_cover_index": 0,
			"photo_images":      imageURLs,
		},
		"post_mode":  "DIRECT_POST",
		"media_type": "PHOTO",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", tiktokContentInitURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TikTok API error: %s - %s", resp.Status, string(responseBody))
	}

	var publishResponse PublishVideoResponse
	if err := json.Unmarshal(responseBody, &publishResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors (TikTok uses "ok" to indicate success)
	if publishResponse.Error.Code != "" && publishResponse.Error.Code != "ok" {
		return nil, fmt.Errorf("TikTok API error: %s - %s", publishResponse.Error.Code, publishResponse.Error.Message)
	}

	return &publishResponse, nil
}
