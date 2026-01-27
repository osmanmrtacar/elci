package services

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
)

// PlatformService interface (duplicated to avoid circular dependency)
type PlatformService interface {
	GetPlatformName() models.Platform
	GetRequiredScopes() []string
	GenerateAuthURL() (AuthURLResponse, error)
	ExchangeCodeForTokens(code string, additionalParams map[string]string) (*TokenResponse, error)
	RefreshAccessToken(refreshToken string) (*TokenResponse, error)
	GetUserInfo(accessToken string) (*UserInfo, error)
	UploadMedia(accessToken string, mediaURL string) (string, error)
	CreatePost(accessToken string, content PostContent) (*PostResponse, error)
	GetPostStatus(accessToken string, postID string) (*PostStatusResponse, error)
}

// platformServiceAdapter wraps an interface{} and implements PlatformService
// This allows us to work with platform services without importing the platform package
type platformServiceAdapter struct {
	service interface{}
}

func (a *platformServiceAdapter) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	svcValue := reflect.ValueOf(a.service)
	method := svcValue.MethodByName("RefreshAccessToken")

	if !method.IsValid() {
		return nil, fmt.Errorf("service does not implement RefreshAccessToken")
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(refreshToken)})

	if len(results) != 2 {
		return nil, fmt.Errorf("RefreshAccessToken method has wrong number of return values")
	}

	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	// The result might be *platform.TokenResponse, so we need to convert it
	respValue := results[0]
	if respValue.IsNil() {
		return nil, nil
	}

	// Extract fields using reflection
	respElem := respValue.Elem()
	return &TokenResponse{
		AccessToken:  respElem.FieldByName("AccessToken").String(),
		RefreshToken: respElem.FieldByName("RefreshToken").String(),
		ExpiresIn:    int(respElem.FieldByName("ExpiresIn").Int()),
		TokenType:    respElem.FieldByName("TokenType").String(),
		Scope:        respElem.FieldByName("Scope").String(),
	}, nil
}

func (a *platformServiceAdapter) UploadMedia(accessToken string, mediaURL string) (string, error) {
	svcValue := reflect.ValueOf(a.service)
	method := svcValue.MethodByName("UploadMedia")

	if !method.IsValid() {
		return "", fmt.Errorf("service does not implement UploadMedia")
	}

	results := method.Call([]reflect.Value{
		reflect.ValueOf(accessToken),
		reflect.ValueOf(mediaURL),
	})

	if len(results) != 2 {
		return "", fmt.Errorf("UploadMedia method has wrong number of return values")
	}

	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return "", err
	}

	return results[0].String(), nil
}

func (a *platformServiceAdapter) CreatePost(accessToken string, content PostContent) (*PostResponse, error) {
	// Use reflection to call CreatePost with the correct type
	svcValue := reflect.ValueOf(a.service)
	method := svcValue.MethodByName("CreatePost")

	if !method.IsValid() {
		return nil, fmt.Errorf("service does not implement CreatePost")
	}

	// Get the type of the second parameter (PostContent) from the method
	methodType := method.Type()
	if methodType.NumIn() != 2 {
		return nil, fmt.Errorf("CreatePost method has wrong number of parameters")
	}

	// Create a new instance of the platform's PostContent type
	contentType := methodType.In(1)
	contentValue := reflect.New(contentType).Elem()

	// Set fields using reflection
	contentValue.FieldByName("Text").SetString(content.Text)
	contentValue.FieldByName("MediaURL").SetString(content.MediaURL)

	// Set MediaURLs slice (for carousel/multi-image)
	mediaURLsField := contentValue.FieldByName("MediaURLs")
	if mediaURLsField.IsValid() && mediaURLsField.CanSet() {
		mediaURLsSlice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(content.MediaURLs), len(content.MediaURLs))
		for i, url := range content.MediaURLs {
			mediaURLsSlice.Index(i).SetString(url)
		}
		mediaURLsField.Set(mediaURLsSlice)
	}

	// Set MediaIDs slice
	mediaIDsField := contentValue.FieldByName("MediaIDs")
	mediaIDsSlice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(content.MediaIDs), len(content.MediaIDs))
	for i, id := range content.MediaIDs {
		mediaIDsSlice.Index(i).SetString(id)
	}
	mediaIDsField.Set(mediaIDsSlice)

	// Set TikTok settings if provided (for TikTok platform)
	if content.TikTokSettings != nil {
		tiktokSettingsField := contentValue.FieldByName("TikTokSettings")
		if tiktokSettingsField.IsValid() && tiktokSettingsField.CanSet() {
			// Get the type of TikTokSettings struct in the platform package
			tiktokSettingsType := tiktokSettingsField.Type().Elem()
			tiktokSettingsValue := reflect.New(tiktokSettingsType).Elem()

			// Set each field
			if f := tiktokSettingsValue.FieldByName("Title"); f.IsValid() {
				f.SetString(content.TikTokSettings.Title)
			}
			if f := tiktokSettingsValue.FieldByName("PrivacyLevel"); f.IsValid() {
				f.SetString(content.TikTokSettings.PrivacyLevel)
			}
			if f := tiktokSettingsValue.FieldByName("AllowComment"); f.IsValid() {
				f.SetBool(content.TikTokSettings.AllowComment)
			}
			if f := tiktokSettingsValue.FieldByName("AllowDuet"); f.IsValid() {
				f.SetBool(content.TikTokSettings.AllowDuet)
			}
			if f := tiktokSettingsValue.FieldByName("AllowStitch"); f.IsValid() {
				f.SetBool(content.TikTokSettings.AllowStitch)
			}
			if f := tiktokSettingsValue.FieldByName("IsBrandContent"); f.IsValid() {
				f.SetBool(content.TikTokSettings.IsBrandContent)
			}
			if f := tiktokSettingsValue.FieldByName("IsBrandOrganic"); f.IsValid() {
				f.SetBool(content.TikTokSettings.IsBrandOrganic)
			}
			if f := tiktokSettingsValue.FieldByName("AutoAddMusic"); f.IsValid() {
				f.SetBool(content.TikTokSettings.AutoAddMusic)
			}

			tiktokSettingsField.Set(tiktokSettingsValue.Addr())
		}
	}

	// Call the method
	results := method.Call([]reflect.Value{
		reflect.ValueOf(accessToken),
		contentValue,
	})

	// Check for error (second return value)
	if len(results) != 2 {
		return nil, fmt.Errorf("CreatePost method has wrong number of return values")
	}

	// Get error if any
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	// Get PostResponse (first return value) and convert it
	respValue := results[0]
	if respValue.IsNil() {
		return nil, nil
	}

	// Extract fields using reflection
	respElem := respValue.Elem()
	return &PostResponse{
		PostID:    respElem.FieldByName("PostID").String(),
		PublishID: respElem.FieldByName("PublishID").String(),
		Status:    respElem.FieldByName("Status").String(),
		ShareURL:  respElem.FieldByName("ShareURL").String(),
		ErrorMsg:  respElem.FieldByName("ErrorMsg").String(),
	}, nil
}

func (a *platformServiceAdapter) GetPlatformName() models.Platform {
	type namer interface {
		GetPlatformName() models.Platform
	}
	if svc, ok := a.service.(namer); ok {
		return svc.GetPlatformName()
	}
	return ""
}

func (a *platformServiceAdapter) GetRequiredScopes() []string {
	type scoper interface {
		GetRequiredScopes() []string
	}
	if svc, ok := a.service.(scoper); ok {
		return svc.GetRequiredScopes()
	}
	return nil
}

func (a *platformServiceAdapter) GenerateAuthURL() (AuthURLResponse, error) {
	type authURLGenerator interface {
		GenerateAuthURL() (AuthURLResponse, error)
	}
	if svc, ok := a.service.(authURLGenerator); ok {
		return svc.GenerateAuthURL()
	}
	return AuthURLResponse{}, fmt.Errorf("service does not implement GenerateAuthURL")
}

func (a *platformServiceAdapter) ExchangeCodeForTokens(code string, additionalParams map[string]string) (*TokenResponse, error) {
	type tokenExchanger interface {
		ExchangeCodeForTokens(string, map[string]string) (*TokenResponse, error)
	}
	if svc, ok := a.service.(tokenExchanger); ok {
		return svc.ExchangeCodeForTokens(code, additionalParams)
	}
	return nil, fmt.Errorf("service does not implement ExchangeCodeForTokens")
}

func (a *platformServiceAdapter) GetUserInfo(accessToken string) (*UserInfo, error) {
	type userInfoGetter interface {
		GetUserInfo(string) (*UserInfo, error)
	}
	if svc, ok := a.service.(userInfoGetter); ok {
		return svc.GetUserInfo(accessToken)
	}
	return nil, fmt.Errorf("service does not implement GetUserInfo")
}

func (a *platformServiceAdapter) GetPostStatus(accessToken string, postID string) (*PostStatusResponse, error) {
	type statusGetter interface {
		GetPostStatus(string, string) (*PostStatusResponse, error)
	}
	if svc, ok := a.service.(statusGetter); ok {
		return svc.GetPostStatus(accessToken, postID)
	}
	return nil, fmt.Errorf("service does not implement GetPostStatus")
}

// AuthURLResponse contains the OAuth authorization URL and associated data
type AuthURLResponse struct {
	URL          string
	State        string
	CodeVerifier string
}

// TokenResponse contains OAuth tokens received from the platform
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
	Scope        string
}

// UserInfo contains basic user information from the platform
type UserInfo struct {
	PlatformUserID string
	Username       string
	DisplayName    string
	AvatarURL      string
	Email          string
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
	Text           string
	MediaURL       string   // Primary media URL (for single media posts)
	MediaURLs      []string // Multiple media URLs (for carousel/multi-image posts)
	MediaIDs       []string
	TikTokSettings *TikTokSettings // TikTok-specific settings
}

// PostResponse contains the result of creating a post
type PostResponse struct {
	PostID    string
	PublishID string
	Status    string
	ShareURL  string
	ErrorMsg  string
}

// PostStatusResponse contains the current status of a post
type PostStatusResponse struct {
	Status          string
	PostID          string
	ShareURL        string
	FailReason      string
	ProgressPercent int
}

// PlatformRegistry interface to avoid circular dependency
// Note: This must match the actual platform.PlatformRegistry interface
type PlatformRegistry interface {
	Get(platform models.Platform) (interface{}, error)
	IsSupported(platform models.Platform) bool
}

type MultiPlatformPostService struct {
	postRepo               *models.PostRepository
	tokenRepo              *models.TokenRepository
	platformConnectionRepo *models.PlatformConnectionRepository
	platformRegistry       PlatformRegistry
}

// NewMultiPlatformPostService creates a new multi-platform post service
func NewMultiPlatformPostService(
	postRepo *models.PostRepository,
	tokenRepo *models.TokenRepository,
	platformConnectionRepo *models.PlatformConnectionRepository,
	platformRegistry PlatformRegistry,
) *MultiPlatformPostService {
	return &MultiPlatformPostService{
		postRepo:               postRepo,
		tokenRepo:              tokenRepo,
		platformConnectionRepo: platformConnectionRepo,
		platformRegistry:       platformRegistry,
	}
}

// CreateMultiPlatformPostRequest represents a request to create a post on multiple platforms
type CreateMultiPlatformPostRequest struct {
	Platforms      []models.Platform `json:"platforms"`                 // ["tiktok", "x"]
	MediaURL       string            `json:"media_url"`                 // Primary video/image URL (for single media)
	MediaURLs      []string          `json:"media_urls"`                // Multiple media URLs (for carousel/multi-image)
	Caption        string            `json:"caption"`                   // Post text/caption
	TikTokSettings *TikTokSettings   `json:"tiktok_settings,omitempty"` // TikTok-specific settings
}

// CreateMultiPlatformPostResponse represents the response after creating posts
type CreateMultiPlatformPostResponse struct {
	Posts  []*models.Post    `json:"posts"`
	Errors map[string]string `json:"errors,omitempty"`
}

// CreateMultiPlatformPost creates a post on multiple platforms simultaneously
func (s *MultiPlatformPostService) CreateMultiPlatformPost(userID int64, req CreateMultiPlatformPostRequest) (*CreateMultiPlatformPostResponse, error) {
	if len(req.Platforms) == 0 {
		return nil, fmt.Errorf("at least one platform must be specified")
	}

	// Ensure we have at least one media URL
	if req.MediaURL == "" && len(req.MediaURLs) == 0 {
		return nil, fmt.Errorf("media URL is required")
	}

	// Use MediaURLs if provided, otherwise fall back to MediaURL
	mediaURLs := req.MediaURLs
	if len(mediaURLs) == 0 && req.MediaURL != "" {
		mediaURLs = []string{req.MediaURL}
	}

	// Validate that user has connected all requested platforms
	connectedPlatforms, err := s.platformConnectionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected platforms: %w", err)
	}

	// Create a map of connected platforms for quick lookup
	platformMap := make(map[models.Platform]bool)
	for _, conn := range connectedPlatforms {
		if conn.IsActive {
			platformMap[conn.Platform] = true
		}
	}

	// Validate all requested platforms are connected
	var notConnected []string
	for _, plt := range req.Platforms {
		if !platformMap[plt] {
			notConnected = append(notConnected, string(plt))
		}
	}

	if len(notConnected) > 0 {
		return nil, fmt.Errorf("platforms not connected: %v", notConnected)
	}

	// Create post records for each platform
	posts := make([]*models.Post, 0, len(req.Platforms))
	errors := make(map[string]string)
	var mu sync.Mutex

	// Detect media type from primary URL
	primaryMediaURL := mediaURLs[0]
	mediaType := "video"
	if IsImageURL(primaryMediaURL) {
		mediaType = "image"
	}
	// If multiple images, it's a carousel
	if len(mediaURLs) > 1 && mediaType == "image" {
		mediaType = "carousel"
	}

	for _, plt := range req.Platforms {
		post := &models.Post{
			UserID:    userID,
			Platform:  plt,
			VideoURL:  primaryMediaURL, // Store primary URL in existing field
			Caption:   req.Caption,
			Status:    models.PostStatusPending,
			MediaType: mediaType,
		}

		if err := s.postRepo.Create(post); err != nil {
			log.Printf("Failed to create post record for platform %s: %v", plt, err)
			errors[string(plt)] = fmt.Sprintf("Failed to create post record: %v", err)
			continue
		}

		posts = append(posts, post)

		// Process post asynchronously for each platform
		// Pass TikTok settings only for TikTok platform
		var tiktokSettings *TikTokSettings
		if plt == models.PlatformTikTok && req.TikTokSettings != nil {
			tiktokSettings = req.TikTokSettings
		}
		go s.processPlatformPost(post.ID, userID, plt, tiktokSettings, mediaURLs, &mu, errors)
	}

	response := &CreateMultiPlatformPostResponse{
		Posts: posts,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	return response, nil
}

// processPlatformPost handles posting to a specific platform asynchronously
func (s *MultiPlatformPostService) processPlatformPost(postID int64, userID int64, plt models.Platform, tiktokSettings *TikTokSettings, mediaURLs []string, mu *sync.Mutex, errors map[string]string) {
	// Update status to processing
	if err := s.postRepo.UpdateStatus(postID, models.PostStatusProcessing, ""); err != nil {
		log.Printf("Failed to update post %d status to processing: %v", postID, err)
		return
	}

	// Get post details
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		log.Printf("Failed to get post %d: %v", postID, err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, "Failed to retrieve post details")
		return
	}

	// Get platform service
	rawService, err := s.platformRegistry.Get(plt)
	if err != nil {
		log.Printf("Platform %s not found: %v", plt, err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, fmt.Sprintf("Platform %s not available", plt))
		mu.Lock()
		errors[string(plt)] = fmt.Sprintf("Platform not available: %v", err)
		mu.Unlock()
		return
	}

	// Wrap rawService with adapter to implement PlatformService interface
	platformService := &platformServiceAdapter{service: rawService}

	// Get valid access token for this platform
	token, err := s.tokenRepo.GetByUserIDAndPlatform(userID, plt)
	if err != nil {
		log.Printf("Failed to get token for user %d on platform %s: %v", userID, plt, err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, "Failed to get access token")
		mu.Lock()
		errors[string(plt)] = "Failed to get access token"
		mu.Unlock()
		return
	}

	// Check if token needs refresh
	// Refresh proactively if token will expire within 7 days (for Instagram long-lived tokens)
	// or if it's already expired
	tokenExpiresSoon := time.Now().Add(7 * 24 * time.Hour).After(token.ExpiresAt)
	tokenExpired := time.Now().After(token.ExpiresAt)

	if tokenExpired {
		// Token is already expired - for Instagram, this means re-authentication is required
		log.Printf("Token already expired for user %d on %s at %v", userID, plt, token.ExpiresAt)

		// Try to refresh anyway (works for TikTok, may fail for Instagram)
		tokenResp, err := platformService.RefreshAccessToken(token.RefreshToken)
		if err != nil {
			log.Printf("Failed to refresh expired token: %v", err)
			errorMsg := "Access token has expired. Please reconnect your account."
			if plt == models.PlatformInstagram {
				errorMsg = "Instagram access token has expired. Please disconnect and reconnect your Instagram account to continue posting."
			}
			s.postRepo.UpdateStatus(postID, models.PostStatusFailed, errorMsg)
			mu.Lock()
			errors[string(plt)] = errorMsg
			mu.Unlock()
			return
		}

		// Token refresh succeeded
		token.AccessToken = tokenResp.AccessToken
		token.RefreshToken = tokenResp.RefreshToken
		token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		if err := s.tokenRepo.CreateOrUpdateForPlatform(token); err != nil {
			log.Printf("Failed to update token: %v", err)
		}
		log.Printf("Successfully refreshed expired token for user %d on %s", userID, plt)
	} else if tokenExpiresSoon {
		// Token will expire soon - refresh proactively
		log.Printf("Token expiring soon for user %d on %s (expires at %v), refreshing proactively...", userID, plt, token.ExpiresAt)
		tokenResp, err := platformService.RefreshAccessToken(token.RefreshToken)
		if err != nil {
			// Log warning but continue with existing token since it's still valid
			log.Printf("Warning: Failed to proactively refresh token for %s: %v", plt, err)
		} else {
			// Update token in database
			token.AccessToken = tokenResp.AccessToken
			token.RefreshToken = tokenResp.RefreshToken
			token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
			if err := s.tokenRepo.CreateOrUpdateForPlatform(token); err != nil {
				log.Printf("Failed to update token: %v", err)
			}
			log.Printf("Successfully refreshed token proactively for user %d on %s, new expiry: %v", userID, plt, token.ExpiresAt)
		}
	}

	// Upload media if needed (for platforms like X that require upload before posting)
	var mediaIDs []string
	if plt == models.PlatformX {
		log.Printf("Uploading %d media file(s) to %s for post %d", len(mediaURLs), plt, postID)
		for i, mediaURL := range mediaURLs {
			mediaID, err := platformService.UploadMedia(token.AccessToken, mediaURL)
			if err != nil {
				log.Printf("Failed to upload media %d to %s: %v", i+1, plt, err)
				s.postRepo.UpdateStatus(postID, models.PostStatusFailed, fmt.Sprintf("Media upload failed: %v", err))
				mu.Lock()
				errors[string(plt)] = fmt.Sprintf("Media upload failed: %v", err)
				mu.Unlock()
				return
			}
			mediaIDs = append(mediaIDs, mediaID)
			log.Printf("Media %d uploaded to %s: %s", i+1, plt, mediaID)
		}
	}

	// Create post on platform
	postContent := PostContent{
		Text:           post.Caption,
		MediaURL:       mediaURLs[0], // Primary URL
		MediaURLs:      mediaURLs,    // All URLs for carousel/multi-image
		MediaIDs:       mediaIDs,
		TikTokSettings: tiktokSettings,
	}

	postResp, err := platformService.CreatePost(token.AccessToken, postContent)
	if err != nil {
		log.Printf("Failed to create post on %s: %v", plt, err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, fmt.Sprintf("Post creation failed: %v", err))
		mu.Lock()
		errors[string(plt)] = fmt.Sprintf("Post creation failed: %v", err)
		mu.Unlock()
		return
	}

	log.Printf("Post created on %s with ID: %s (status: %s)", plt, postResp.PostID, postResp.Status)

	// Handle platform-specific processing
	if plt == models.PlatformTikTok {
		// TikTok requires polling for status
		s.pollTikTokStatus(postID, userID, postResp.PostID, token.AccessToken, platformService)
	} else if plt == models.PlatformX {
		// X posts are published immediately
		if err := s.postRepo.MarkPublishedWithPlatform(postID, postResp.PostID); err != nil {
			log.Printf("Failed to mark post %d as published: %v", postID, err)
		} else {
			log.Printf("Post %d successfully published to %s", postID, plt)
		}
	} else if plt == models.PlatformInstagram {
		// Instagram posts are published immediately after successful creation
		if err := s.postRepo.MarkPublishedWithPlatform(postID, postResp.PostID); err != nil {
			log.Printf("Failed to mark post %d as published: %v", postID, err)
		} else {
			log.Printf("Post %d successfully published to %s with permalink: %s", postID, plt, postResp.ShareURL)
		}
	}
}

// pollTikTokStatus polls TikTok for post status until complete
func (s *MultiPlatformPostService) pollTikTokStatus(postID int64, userID int64, publishID string, accessToken string, platformService PlatformService) {
	maxAttempts := 60 // Poll for up to 5 minutes
	attempt := 0

	for attempt < maxAttempts {
		time.Sleep(5 * time.Second)
		attempt++

		statusResp, err := platformService.GetPostStatus(accessToken, publishID)
		if err != nil {
			log.Printf("Failed to get publish status: %v", err)
			continue
		}

		log.Printf("TikTok publish status for %s: %s", publishID, statusResp.Status)

		switch statusResp.Status {
		case "published":
			// Video successfully published
			if err := s.postRepo.MarkPublishedWithPlatform(postID, publishID); err != nil {
				log.Printf("Failed to mark post %d as published: %v", postID, err)
			} else {
				log.Printf("Post %d successfully published to TikTok", postID)
			}
			return

		case "failed":
			// Publishing failed
			failReason := statusResp.FailReason
			if failReason == "" {
				failReason = "Unknown error"
			}
			s.postRepo.UpdateStatus(postID, models.PostStatusFailed, fmt.Sprintf("TikTok publish failed: %s", failReason))
			log.Printf("Post %d failed: %s", postID, failReason)
			return

		case "processing":
			log.Printf("Post %d: TikTok is processing the video (%d%%)", postID, statusResp.ProgressPercent)
		}
	}

	// Timeout reached
	s.postRepo.UpdateStatus(postID, models.PostStatusFailed, "Publishing timeout - took too long")
	log.Printf("Post %d timed out after %d attempts", postID, maxAttempts)
}

// GetPostByID retrieves a post by ID
func (s *MultiPlatformPostService) GetPostByID(postID int64, userID int64) (*models.Post, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// Verify post belongs to user
	if post.UserID != userID {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

// GetUserPosts retrieves all posts for a user (optionally filtered by platform)
func (s *MultiPlatformPostService) GetUserPosts(userID int64, platformFilter models.Platform, limit, offset int) ([]*models.Post, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if platformFilter != "" {
		return s.postRepo.GetByUserIDAndPlatform(userID, platformFilter, limit, offset)
	}

	return s.postRepo.GetByUserID(userID, limit, offset)
}

// GetPostStatus retrieves the current status of a post
func (s *MultiPlatformPostService) GetPostStatus(postID int64, userID int64) (*models.Post, error) {
	return s.GetPostByID(postID, userID)
}
