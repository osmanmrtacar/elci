package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/osmanmertacar/sosyal/backend/internal/config"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
)

type TokenService struct {
	config       *config.Config
	tokenRepo    *models.TokenRepository
	tiktokService *TikTokService
}

// TikTokTokenResponse represents the response from TikTok's token endpoint
type TikTokTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	OpenID       string `json:"open_id"`
}

// NewTokenService creates a new token service
func NewTokenService(cfg *config.Config, tokenRepo *models.TokenRepository) *TokenService {
	return &TokenService{
		config:    cfg,
		tokenRepo: tokenRepo,
	}
}

// SetTikTokService sets the TikTok service (to avoid circular dependency)
func (s *TokenService) SetTikTokService(tiktokService *TikTokService) {
	s.tiktokService = tiktokService
}

// SaveTokens saves TikTok OAuth tokens to the database (encrypted)
func (s *TokenService) SaveTokens(userID int64, tokenResponse *TikTokTokenResponse) error {
	// Encrypt tokens before storing
	encryptedAccessToken, err := s.encryptToken(tokenResponse.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	encryptedRefreshToken, err := s.encryptToken(tokenResponse.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	// Create token model
	token := &models.Token{
		UserID:       userID,
		AccessToken:  encryptedAccessToken,
		RefreshToken: encryptedRefreshToken,
		TokenType:    tokenResponse.TokenType,
		ExpiresAt:    expiresAt,
		Scope:        tokenResponse.Scope,
	}

	// Save to database (creates or updates)
	if err := s.tokenRepo.CreateOrUpdate(token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// GetValidToken retrieves a valid access token for a user (refreshes if expired)
func (s *TokenService) GetValidToken(userID int64) (string, error) {
	// Get token from database
	token, err := s.tokenRepo.GetByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	// Decrypt tokens
	accessToken, err := s.decryptToken(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// Check if token is expired (with 5 minute buffer)
	if time.Now().Add(5 * time.Minute).After(token.ExpiresAt) {
		// Token is expired or about to expire, refresh it
		refreshToken, err := s.decryptToken(token.RefreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt refresh token: %w", err)
		}

		// Refresh the token
		newTokenResponse, err := s.RefreshToken(refreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}

		// Save new tokens
		if err := s.SaveTokens(userID, newTokenResponse); err != nil {
			return "", fmt.Errorf("failed to save refreshed tokens: %w", err)
		}

		return newTokenResponse.AccessToken, nil
	}

	return accessToken, nil
}

// RefreshToken refreshes an expired access token using the refresh token
func (s *TokenService) RefreshToken(refreshToken string) (*TikTokTokenResponse, error) {
	if s.tiktokService == nil {
		return nil, fmt.Errorf("TikTok service not initialized")
	}

	return s.tiktokService.RefreshAccessToken(refreshToken)
}

// DeleteTokens deletes all tokens for a user (logout)
func (s *TokenService) DeleteTokens(userID int64) error {
	if err := s.tokenRepo.DeleteByUserID(userID); err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}
	return nil
}

// encryptToken encrypts a token using AES-GCM
func (s *TokenService) encryptToken(plaintext string) (string, error) {
	// Use JWT secret as encryption key (first 32 bytes for AES-256)
	key := []byte(s.config.JWT.Secret)
	if len(key) > 32 {
		key = key[:32]
	} else if len(key) < 32 {
		// Pad key if too short
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptToken decrypts a token using AES-GCM
func (s *TokenService) decryptToken(ciphertext string) (string, error) {
	// Use JWT secret as encryption key (first 32 bytes for AES-256)
	key := []byte(s.config.JWT.Secret)
	if len(key) > 32 {
		key = key[:32]
	} else if len(key) < 32 {
		// Pad key if too short
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], string(data[nonceSize:])
	plaintext, err := gcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
