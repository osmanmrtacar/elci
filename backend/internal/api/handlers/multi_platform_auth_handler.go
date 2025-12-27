package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/osmanmertacar/sosyal/backend/internal/api/middleware"
	"github.com/osmanmertacar/sosyal/backend/internal/config"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services/platform"
)

type MultiPlatformAuthHandler struct {
	config                  *config.Config
	platformRegistry        *platform.PlatformRegistry
	userRepo                *models.UserRepository
	tokenRepo               *models.TokenRepository
	platformConnectionRepo  *models.PlatformConnectionRepository
	oauthSessionRepo        *models.OAuthSessionRepository
}

// NewMultiPlatformAuthHandler creates a new multi-platform auth handler
func NewMultiPlatformAuthHandler(
	cfg *config.Config,
	platformRegistry *platform.PlatformRegistry,
	userRepo *models.UserRepository,
	tokenRepo *models.TokenRepository,
	platformConnectionRepo *models.PlatformConnectionRepository,
	oauthSessionRepo *models.OAuthSessionRepository,
) *MultiPlatformAuthHandler {
	return &MultiPlatformAuthHandler{
		config:                 cfg,
		platformRegistry:       platformRegistry,
		userRepo:               userRepo,
		tokenRepo:              tokenRepo,
		platformConnectionRepo: platformConnectionRepo,
		oauthSessionRepo:       oauthSessionRepo,
	}
}

// TikTokLogin initiates the TikTok OAuth flow
func (h *MultiPlatformAuthHandler) TikTokLogin(c *gin.Context) {
	h.handlePlatformLogin(c, models.PlatformTikTok)
}

// XLogin initiates the X (Twitter) OAuth flow
func (h *MultiPlatformAuthHandler) XLogin(c *gin.Context) {
	h.handlePlatformLogin(c, models.PlatformX)
}

// InstagramLogin initiates the Instagram OAuth flow (via Facebook)
func (h *MultiPlatformAuthHandler) InstagramLogin(c *gin.Context) {
	h.handlePlatformLogin(c, models.PlatformInstagram)
}

// handlePlatformLogin is a generic handler for initiating OAuth flow for any platform
func (h *MultiPlatformAuthHandler) handlePlatformLogin(c *gin.Context, platformType models.Platform) {
	// Get platform service
	rawService, err := h.platformRegistry.Get(platformType)
	if err != nil {
		log.Printf("Platform %s not supported: %v", platformType, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Platform %s is not configured", platformType),
		})
		return
	}

	// Type assert to platform.PlatformService
	platformService, ok := rawService.(platform.PlatformService)
	if !ok {
		log.Printf("Platform %s service does not implement PlatformService", platformType)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid platform service",
		})
		return
	}

	// Generate OAuth URL
	authResp, err := platformService.GenerateAuthURL()
	if err != nil {
		log.Printf("Failed to generate %s auth URL: %v", platformType, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate authorization URL",
		})
		return
	}

	// Store OAuth session in database (for PKCE and state validation)
	expiresAt := time.Now().Add(10 * time.Minute)
	oauthSession := &models.OAuthSession{
		State:        authResp.State,
		CodeVerifier: authResp.CodeVerifier, // Empty for TikTok, populated for X
		Platform:     platformType,
		ExpiresAt:    expiresAt,
	}

	if err := h.oauthSessionRepo.Create(oauthSession); err != nil {
		log.Printf("Failed to store OAuth session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initialize OAuth session",
		})
		return
	}

	// Store state in cookie as backup (for backward compatibility)
	c.SetCookie("oauth_state", authResp.State, 600, "/", "", false, true)
	c.SetCookie("oauth_platform", string(platformType), 600, "/", "", false, true)

	log.Printf("Redirecting to %s OAuth: state=%s", platformType, authResp.State[:10]+"...")

	// Redirect to platform OAuth page
	c.Redirect(http.StatusTemporaryRedirect, authResp.URL)
}

// TikTokCallback handles the OAuth callback from TikTok
func (h *MultiPlatformAuthHandler) TikTokCallback(c *gin.Context) {
	h.handlePlatformCallback(c, models.PlatformTikTok)
}

// XCallback handles the OAuth callback from X (Twitter)
func (h *MultiPlatformAuthHandler) XCallback(c *gin.Context) {
	h.handlePlatformCallback(c, models.PlatformX)
}

// InstagramCallback handles the OAuth callback from Instagram (via Facebook)
func (h *MultiPlatformAuthHandler) InstagramCallback(c *gin.Context) {
	h.handlePlatformCallback(c, models.PlatformInstagram)
}

// handlePlatformCallback is a generic handler for OAuth callback from any platform
func (h *MultiPlatformAuthHandler) handlePlatformCallback(c *gin.Context, platformType models.Platform) {
	// Get authorization code and state from query params
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		errorMsg := c.Query("error")
		errorDesc := c.Query("error_description")
		log.Printf("%s OAuth error: %s - %s", platformType, errorMsg, errorDesc)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "Authorization failed",
			"error_description": errorDesc,
		})
		return
	}

	// Retrieve OAuth session from database
	oauthSession, err := h.oauthSessionRepo.GetByState(state)
	if err != nil {
		log.Printf("Invalid or expired OAuth state: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired OAuth session",
		})
		return
	}

	// Verify platform matches
	if oauthSession.Platform != platformType {
		log.Printf("Platform mismatch: expected %s, got %s", oauthSession.Platform, platformType)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Platform mismatch",
		})
		return
	}

	// Get platform service
	rawService, err := h.platformRegistry.Get(platformType)
	if err != nil {
		log.Printf("Platform %s not found: %v", platformType, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Platform service not available",
		})
		return
	}

	// Type assert to platform.PlatformService
	platformService, ok := rawService.(platform.PlatformService)
	if !ok {
		log.Printf("Platform %s service does not implement PlatformService", platformType)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid platform service",
		})
		return
	}

	// Exchange code for tokens (with PKCE code_verifier for X)
	additionalParams := make(map[string]string)
	if oauthSession.CodeVerifier != "" {
		additionalParams["code_verifier"] = oauthSession.CodeVerifier
	}

	tokenResp, err := platformService.ExchangeCodeForTokens(code, additionalParams)
	if err != nil {
		log.Printf("Failed to exchange code for %s tokens: %v", platformType, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Delete OAuth session (one-time use)
	if err := h.oauthSessionRepo.DeleteByState(state); err != nil {
		log.Printf("Failed to delete OAuth session: %v", err)
	}

	// Get user info from platform
	userInfo, err := platformService.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Failed to get %s user info: %v", platformType, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user information",
		})
		return
	}

	// Check if user already exists by trying to get from JWT (if logged in)
	var userID int64
	existingUserID, err := middleware.GetUserID(c)
	if err == nil {
		// User is already logged in, add new platform connection
		userID = existingUserID
		log.Printf("Adding %s connection to existing user %d", platformType, userID)
	} else {
		// Check if this platform user already exists
		existingUser, err := h.userRepo.GetByPlatformUserID(platformType, userInfo.PlatformUserID)
		if err == nil && existingUser != nil {
			// User exists, use existing user
			userID = existingUser.ID
			log.Printf("Found existing user %d for %s user %s", userID, platformType, userInfo.PlatformUserID)
		} else {
			// Create new user
			user := &models.User{
				Platform:       platformType,
				PlatformUserID: userInfo.PlatformUserID,
				Username:       userInfo.Username,
				DisplayName:    userInfo.DisplayName,
				AvatarURL:      userInfo.AvatarURL,
				// Legacy fields (for backward compatibility)
				TikTokUserID: userInfo.PlatformUserID,
			}

			if err := h.userRepo.CreateOrUpdate(user); err != nil {
				log.Printf("Failed to create user: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to create user account",
				})
				return
			}
			userID = user.ID
			log.Printf("Created new user %d for %s", userID, platformType)
		}
	}

	// Create or update platform connection
	platformConnection := &models.PlatformConnection{
		UserID:         userID,
		Platform:       platformType,
		PlatformUserID: userInfo.PlatformUserID,
		Username:       userInfo.Username,
		DisplayName:    userInfo.DisplayName,
		AvatarURL:      userInfo.AvatarURL,
		IsActive:       true,
	}

	if err := h.platformConnectionRepo.CreateOrUpdate(platformConnection); err != nil {
		log.Printf("Failed to create platform connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save platform connection",
		})
		return
	}

	// Save tokens to database
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	token := &models.Token{
		UserID:       userID,
		Platform:     platformType,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := h.tokenRepo.CreateOrUpdateForPlatform(token); err != nil {
		log.Printf("Failed to save tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save authentication tokens",
		})
		return
	}

	// Generate JWT session token for frontend
	user, _ := h.userRepo.GetByID(userID)
	jwtToken, err := h.createJWTSession(user)
	if err != nil {
		log.Printf("Failed to create JWT session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create session",
		})
		return
	}

	log.Printf("%s authentication successful for user %d", platformType, userID)

	// Clear OAuth cookies
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	c.SetCookie("oauth_platform", "", -1, "/", "", false, true)

	// Redirect to frontend with token
	frontendURL := h.config.Server.FrontendURL + "/callback?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// GetConnectedPlatforms returns all platforms connected by the current user
func (h *MultiPlatformAuthHandler) GetConnectedPlatforms(c *gin.Context) {
	// Get user ID from JWT
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	// Get platform connections
	connections, err := h.platformConnectionRepo.GetByUserID(userID)
	if err != nil {
		log.Printf("Failed to get platform connections for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve platform connections",
		})
		return
	}

	// Format response
	platformsData := make([]gin.H, 0, len(connections))
	for _, conn := range connections {
		platformsData = append(platformsData, gin.H{
			"platform":     conn.Platform,
			"username":     conn.Username,
			"display_name": conn.DisplayName,
			"avatar_url":   conn.AvatarURL,
			"is_active":    conn.IsActive,
			"connected_at": conn.ConnectedAt,
			"last_used_at": conn.LastUsedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"platforms": platformsData,
	})
}

// DisconnectPlatform disconnects a specific platform from the current user
func (h *MultiPlatformAuthHandler) DisconnectPlatform(c *gin.Context) {
	// Get user ID from JWT
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	// Get platform from URL param
	platformStr := c.Param("platform")
	platformType := models.Platform(platformStr)

	// Validate platform
	if !h.platformRegistry.IsSupported(platformType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid platform",
		})
		return
	}

	// Deactivate platform connection
	if err := h.platformConnectionRepo.Deactivate(userID, platformType); err != nil {
		log.Printf("Failed to deactivate platform connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to disconnect platform",
		})
		return
	}

	// Delete tokens for this platform
	if err := h.tokenRepo.DeleteByUserIDAndPlatform(userID, platformType); err != nil {
		log.Printf("Failed to delete tokens: %v", err)
	}

	log.Printf("User %d disconnected from %s", userID, platformType)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully disconnected from %s", platformType),
	})
}

// createJWTSession creates a JWT token for the user session
func (h *MultiPlatformAuthHandler) createJWTSession(user *models.User) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(h.config.JWT.Expiration).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(h.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// GetCurrentUser returns the current user's information with connected platforms
func (h *MultiPlatformAuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from JWT
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("Failed to get user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Get connected platforms
	connections, err := h.platformConnectionRepo.GetByUserID(userID)
	if err != nil {
		log.Printf("Failed to get platform connections: %v", err)
		connections = []*models.PlatformConnection{} // Empty array on error
	}

	// Format platforms
	platforms := make([]string, 0, len(connections))
	for _, conn := range connections {
		if conn.IsActive {
			platforms = append(platforms, string(conn.Platform))
		}
	}

	// Return user info
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":                 user.ID,
			"username":           user.Username,
			"display_name":       user.DisplayName,
			"avatar_url":         user.AvatarURL,
			"created_at":         user.CreatedAt,
			"connected_platforms": platforms,
		},
	})
}

// Logout logs out the current user
func (h *MultiPlatformAuthHandler) Logout(c *gin.Context) {
	// Get user ID from JWT
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	log.Printf("User %d logged out", userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
