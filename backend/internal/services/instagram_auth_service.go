package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// InstagramAuthService handles Instagram OAuth via Facebook
type InstagramAuthService struct {
	appID       string
	appSecret   string
	redirectURI string
	httpClient  *http.Client
}

// NewInstagramAuthService creates a new Instagram auth service
func NewInstagramAuthService(appID, appSecret, redirectURI string) *InstagramAuthService {
	return &InstagramAuthService{
		appID:       appID,
		appSecret:   appSecret,
		redirectURI: redirectURI,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

// InstagramAuthResponse contains the OAuth authorization URL and state
type InstagramAuthResponse struct {
	URL   string
	State string
}

// InstagramTokenResponse represents the token response from Instagram
type InstagramTokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

// InstagramUserInfoResponse represents user info from Facebook/Instagram
type InstagramUserInfoResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// FacebookPageResponse represents a Facebook page with Instagram account
type FacebookPageResponse struct {
	Data []struct {
		ID                       string                    `json:"id"`
		Name                     string                    `json:"name"`
		AccessToken              string                    `json:"access_token"`
		InstagramBusinessAccount *InstagramBusinessAccount `json:"instagram_business_account,omitempty"`
	} `json:"data"`
}

// InstagramBusinessAccount represents an Instagram business account
type InstagramBusinessAccount struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username"`
}

// GenerateState generates a random state for CSRF protection
func (s *InstagramAuthService) GenerateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// GenerateAuthURL generates the Instagram OAuth URL (via Facebook)
func (s *InstagramAuthService) GenerateAuthURL() (*InstagramAuthResponse, error) {
	state := s.GenerateState()

	// Instagram requires OAuth via Facebook
	scopes := []string{
		"instagram_business_basic",
		"instagram_business_content_publish",
	}

	params := url.Values{}
	params.Set("client_id", s.appID)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("scope", strings.Join(scopes, ","))
	params.Set("response_type", "code")
	params.Set("state", state)

	fmt.Println(params)

	authURL := fmt.Sprintf("https://www.instagram.com/oauth/authorize?%s", params.Encode())

	fmt.Println(authURL)

	return &InstagramAuthResponse{
		URL:   authURL,
		State: state,
	}, nil
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (s *InstagramAuthService) ExchangeCodeForToken(code string) (*InstagramTokenResponse, error) {
	tokenURL := "https://api.instagram.com/oauth/access_token"

	params := url.Values{}
	params.Set("client_id", s.appID)
	params.Set("client_secret", s.appSecret)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)

	resp, err := s.httpClient.PostForm(tokenURL, params)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp InstagramTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// InstagramLongLivedTokenResponse represents the long-lived token response
type InstagramLongLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // Typically 5184000 (60 days)
}

// ExchangeLongLivedToken exchanges a short-lived token for a long-lived token
// Short-lived tokens expire in 1 hour, long-lived tokens last 60 days
func (s *InstagramAuthService) ExchangeLongLivedToken(shortLivedToken string) (*InstagramLongLivedTokenResponse, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/access_token?grant_type=ig_exchange_token&client_secret=%s&access_token=%s",
		s.appSecret, shortLivedToken)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange for long-lived token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("long-lived token exchange failed: %s", string(body))
	}

	var tokenResp InstagramLongLivedTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse long-lived token response: %w", err)
	}

	fmt.Printf("DEBUG: Long-lived token obtained, expires in %d seconds (%.1f days)\n", tokenResp.ExpiresIn, float64(tokenResp.ExpiresIn)/86400)

	return &tokenResp, nil
}

// RefreshLongLivedToken refreshes a long-lived token before it expires
// Can only be refreshed if the token is at least 24 hours old and not expired
// Returns a new long-lived token valid for 60 days
func (s *InstagramAuthService) RefreshLongLivedToken(longLivedToken string) (*InstagramLongLivedTokenResponse, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/refresh_access_token?grant_type=ig_refresh_token&access_token=%s",
		longLivedToken)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh long-lived token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed: %s", string(body))
	}

	var tokenResp InstagramLongLivedTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh token response: %w", err)
	}

	fmt.Printf("DEBUG: Token refreshed, new expiry in %d seconds (%.1f days)\n", tokenResp.ExpiresIn, float64(tokenResp.ExpiresIn)/86400)

	return &tokenResp, nil
}

// GetInstagramUserInfo retrieves Instagram user information using the access token
func (s *InstagramAuthService) GetInstagramUserInfo(accessToken string) (*InstagramUserInfoResponse, error) {
	// Use "me" endpoint to get current user info
	apiURL := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username,account_type&access_token=%s",
		accessToken)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user info failed: %s", string(body))
	}

	var userInfo InstagramUserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}
