package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// XAuthService handles X (Twitter) OAuth 2.0 authentication with PKCE
type XAuthService struct {
	clientID     string
	clientSecret string
	redirectURI  string
	httpClient   *http.Client
}

// NewXAuthService creates a new X authentication service
func NewXAuthService(clientID, clientSecret, redirectURI string) *XAuthService {
	return &XAuthService{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// XTokenResponse represents the OAuth token response from X API
type XTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// XUserInfo represents user information from X API
type XUserInfo struct {
	Data struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Username        string `json:"username"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

// GenerateCodeVerifier creates a random 32-byte base64url-encoded string for PKCE
// CRITICAL: Must use base64.RawURLEncoding (NO padding) to match TypeScript implementation
func (s *XAuthService) GenerateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Use RawURLEncoding (no padding) to match TypeScript base64url
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateCodeChallenge creates SHA256 hash of verifier, base64url encoded (NO padding)
// CRITICAL: Must use base64.RawURLEncoding to match TypeScript implementation
func (s *XAuthService) GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// GenerateState creates a random state parameter for CSRF protection
func (s *XAuthService) GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateAuthURL creates the OAuth authorization URL with PKCE
func (s *XAuthService) GenerateAuthURL() (authURL, state, codeVerifier string, err error) {
	// Generate PKCE parameters
	codeVerifier, err = s.GenerateCodeVerifier()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	codeChallenge := s.GenerateCodeChallenge(codeVerifier)

	// Generate state for CSRF protection
	state, err = s.GenerateState()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Build OAuth URL with PKCE
	scopes := []string{"tweet.read", "tweet.write", "users.read", "offline.access", "media.write"}
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", s.clientID)
	params.Add("redirect_uri", s.redirectURI)
	params.Add("scope", strings.Join(scopes, " "))
	params.Add("state", state)
	params.Add("code_challenge", codeChallenge)
	params.Add("code_challenge_method", "S256")

	authURL = "https://twitter.com/i/oauth2/authorize?" + params.Encode()

	return authURL, state, codeVerifier, nil
}

// ExchangeCodeForToken exchanges authorization code for access token using PKCE
func (s *XAuthService) ExchangeCodeForToken(code, codeVerifier string) (*XTokenResponse, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", s.clientID)
	formData.Set("redirect_uri", s.redirectURI)
	formData.Set("code_verifier", codeVerifier)

	// Create request
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/oauth2/token",
		strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - Basic Auth with client credentials
	credentials := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", s.clientID, s.clientSecret)),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+credentials)

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("X API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp XTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshAccessToken refreshes an expired access token
func (s *XAuthService) RefreshAccessToken(refreshToken string) (*XTokenResponse, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("refresh_token", refreshToken)
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", s.clientID)

	// Create request
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/oauth2/token",
		strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - Basic Auth
	credentials := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", s.clientID, s.clientSecret)),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+credentials)

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("X API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp XTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo fetches user information from X API
func (s *XAuthService) GetUserInfo(accessToken string) (*XUserInfo, error) {
	// Create request
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("user.fields", "profile_image_url")
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("X API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var userInfo XUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	return &userInfo, nil
}
