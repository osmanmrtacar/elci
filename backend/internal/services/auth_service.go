package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/osmanmertacar/sosyal/backend/internal/config"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
)

type AuthService struct {
	config        *config.Config
	tokenService  *TokenService
	tiktokService *TikTokService
	userRepo      *models.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(cfg *config.Config, tokenService *TokenService, tiktokService *TikTokService, userRepo *models.UserRepository) *AuthService {
	return &AuthService{
		config:        cfg,
		tokenService:  tokenService,
		tiktokService: tiktokService,
		userRepo:      userRepo,
	}
}

// GenerateAuthURL generates the TikTok OAuth authorization URL
func (s *AuthService) GenerateAuthURL() (string, string, error) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
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
	return authURL, state, nil
}

// HandleCallback processes the OAuth callback and creates a user session
func (s *AuthService) HandleCallback(code, state string) (string, *models.User, error) {
	// Exchange code for tokens
	tokenResponse, err := s.tiktokService.ExchangeCodeForTokens(code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	// Log token info for debugging (first 10 chars only for security)
	tokenPreview := tokenResponse.AccessToken
	if len(tokenPreview) > 10 {
		tokenPreview = tokenPreview[:10] + "..."
	}
	log.Printf("DEBUG: Received access token: %s\n", tokenPreview)
	log.Printf("DEBUG: Token type: %s\n", tokenResponse.TokenType)
	log.Printf("DEBUG: Expires in: %d seconds\n", tokenResponse.ExpiresIn)

	// Get user info from TikTok
	userInfo, err := s.tiktokService.GetUserInfo(tokenResponse.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Create or update user in database
	// Note: TikTok v2 API doesn't provide username, use display_name or open_id as fallback
	username := userInfo.DisplayName
	if username == "" {
		username = userInfo.OpenID
	}
	user := &models.User{
		TikTokUserID: userInfo.OpenID,
		Username:     username,
		DisplayName:  userInfo.DisplayName,
		AvatarURL:    userInfo.AvatarURL,
	}

	if err := s.userRepo.CreateOrUpdate(user); err != nil {
		return "", nil, fmt.Errorf("failed to create/update user: %w", err)
	}

	// Save tokens to database
	if err := s.tokenService.SaveTokens(user.ID, tokenResponse); err != nil {
		return "", nil, fmt.Errorf("failed to save tokens: %w", err)
	}

	// Generate JWT session token for frontend
	jwtToken, err := s.CreateJWTSession(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create JWT session: %w", err)
	}

	return jwtToken, user, nil
}

// CreateJWTSession creates a JWT token for the user session
func (s *AuthService) CreateJWTSession(user *models.User) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"tiktok_user_id": user.TikTokUserID,
		"username":       user.Username,
		"exp":            time.Now().Add(s.config.JWT.Expiration).Unix(),
		"iat":            time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// ValidateState validates the OAuth state parameter (CSRF protection)
func (s *AuthService) ValidateState(receivedState, expectedState string) bool {
	return receivedState != "" && receivedState == expectedState
}

// generateRandomState generates a random state parameter for OAuth
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
