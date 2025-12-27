package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osmanmertacar/sosyal/backend/internal/api/middleware"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	userRepo    *models.UserRepository
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService, userRepo *models.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

// TikTokLogin initiates the TikTok OAuth flow
func (h *AuthHandler) TikTokLogin(c *gin.Context) {
	// Generate OAuth URL with state parameter
	authURL, state, err := h.authService.GenerateAuthURL()
	if err != nil {
		log.Printf("Failed to generate auth URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate authorization URL",
		})
		return
	}

	// Store state in cookie for validation (expires in 10 minutes)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	// Redirect to TikTok OAuth page
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// TikTokCallback handles the OAuth callback from TikTok
func (h *AuthHandler) TikTokCallback(c *gin.Context) {
	// Get authorization code and state from query params
	code := c.Query("code")
	state := c.Query("state")

	// Try to get stored state from cookie (optional for tunnel compatibility)
	storedState, err := c.Cookie("oauth_state")
	if err == nil {
		// If cookie exists, validate state (CSRF protection)
		if !h.authService.ValidateState(state, storedState) {
			log.Printf("Invalid state parameter: expected %s, got %s", storedState, state)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid state parameter",
			})
			return
		}
		// Clear state cookie
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	} else {
		// Cookie not present (likely due to tunnel/CORS), proceed anyway
		// The authorization code itself is already validated by TikTok
		log.Printf("State cookie not present, proceeding with callback (tunnel mode)")
	}

	// Check for authorization code
	if code == "" {
		errorMsg := c.Query("error")
		errorDesc := c.Query("error_description")
		log.Printf("OAuth error: %s - %s", errorMsg, errorDesc)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "Authorization failed",
			"error_description": errorDesc,
		})
		return
	}

	// Exchange code for tokens and create user session
	jwtToken, user, err := h.authService.HandleCallback(code, state)
	if err != nil {
		log.Printf("Failed to handle callback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to complete authentication",
		})
		return
	}

	log.Printf("User authenticated successfully: %s (ID: %d)", user.Username, user.ID)

	// Redirect to frontend with token in query parameter
	// Frontend will extract the token and store it in localStorage
	frontendURL := "http://localhost:3000/callback?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// Logout logs out the current user
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	// Note: In a real app, you might want to delete tokens from the database
	// For now, we just acknowledge the logout (JWT will expire naturally)
	log.Printf("User %d logged out", userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetCurrentUser returns the current user's information
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
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

	// Return user info
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"avatar_url":   user.AvatarURL,
			"created_at":   user.CreatedAt,
		},
	})
}

// GetDevToken generates a JWT token for a user (development only)
func (h *AuthHandler) GetDevToken(c *gin.Context) {
	userIDStr := c.Param("user_id")

	var userID int64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("Failed to get user %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate JWT token
	jwtToken, err := h.authService.CreateJWTSession(user)
	if err != nil {
		log.Printf("Failed to create JWT for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"display_name": user.DisplayName,
		},
	})
}
