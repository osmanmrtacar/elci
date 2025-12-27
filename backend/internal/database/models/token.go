package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Token struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Platform     Platform  `json:"platform"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope"`
	CodeVerifier string    `json:"code_verifier,omitempty"` // For PKCE (X platform)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TokenRepository struct {
	DB *sql.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

// Create creates a new token
func (r *TokenRepository) Create(token *Token) error {
	query := `
		INSERT INTO tokens (user_id, platform, access_token, refresh_token, token_type, expires_at, scope, code_verifier, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.DB.Exec(query,
		token.UserID, token.Platform, token.AccessToken, token.RefreshToken,
		token.TokenType, token.ExpiresAt, token.Scope, token.CodeVerifier, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	token.ID = id
	token.CreatedAt = now
	token.UpdatedAt = now
	return nil
}

// GetByUserID retrieves the most recent token for a user
func (r *TokenRepository) GetByUserID(userID int64) (*Token, error) {
	query := `
		SELECT id, user_id, platform, access_token, refresh_token, token_type, expires_at, scope, code_verifier, created_at, updated_at
		FROM tokens
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	token := &Token{}
	err := r.DB.QueryRow(query, userID).Scan(
		&token.ID, &token.UserID, &token.Platform, &token.AccessToken, &token.RefreshToken,
		&token.TokenType, &token.ExpiresAt, &token.Scope, &token.CodeVerifier, &token.CreatedAt, &token.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("token not found for user")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return token, nil
}

// Update updates a token
func (r *TokenRepository) Update(token *Token) error {
	query := `
		UPDATE tokens
		SET platform = ?, access_token = ?, refresh_token = ?, expires_at = ?, scope = ?, code_verifier = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.DB.Exec(query,
		token.Platform, token.AccessToken, token.RefreshToken, token.ExpiresAt,
		token.Scope, token.CodeVerifier, now, token.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}
	token.UpdatedAt = now
	return nil
}

// Delete deletes a token
func (r *TokenRepository) Delete(id int64) error {
	query := "DELETE FROM tokens WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

// DeleteByUserID deletes all tokens for a user
func (r *TokenRepository) DeleteByUserID(userID int64) error {
	query := "DELETE FROM tokens WHERE user_id = ?"
	_, err := r.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens for user: %w", err)
	}
	return nil
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// CreateOrUpdate creates a new token or updates if one exists for the user
func (r *TokenRepository) CreateOrUpdate(token *Token) error {
	existingToken, err := r.GetByUserID(token.UserID)
	if err != nil && err.Error() != "token not found for user" {
		return err
	}

	if existingToken == nil {
		return r.Create(token)
	}

	token.ID = existingToken.ID
	return r.Update(token)
}

// GetByUserIDAndPlatform retrieves the most recent token for a user and platform
func (r *TokenRepository) GetByUserIDAndPlatform(userID int64, platform Platform) (*Token, error) {
	query := `
		SELECT id, user_id, platform, access_token, refresh_token, token_type, expires_at, scope, code_verifier, created_at, updated_at
		FROM tokens
		WHERE user_id = ? AND platform = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	token := &Token{}
	err := r.DB.QueryRow(query, userID, platform).Scan(
		&token.ID, &token.UserID, &token.Platform, &token.AccessToken, &token.RefreshToken,
		&token.TokenType, &token.ExpiresAt, &token.Scope, &token.CodeVerifier, &token.CreatedAt, &token.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("token not found for user and platform")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return token, nil
}

// CreateOrUpdateForPlatform creates a new token or updates if one exists for the user and platform
func (r *TokenRepository) CreateOrUpdateForPlatform(token *Token) error {
	existingToken, err := r.GetByUserIDAndPlatform(token.UserID, token.Platform)
	if err != nil && err.Error() != "token not found for user and platform" {
		return err
	}

	if existingToken == nil {
		return r.Create(token)
	}

	token.ID = existingToken.ID
	return r.Update(token)
}

// DeleteByUserIDAndPlatform deletes all tokens for a user and platform
func (r *TokenRepository) DeleteByUserIDAndPlatform(userID int64, platform Platform) error {
	query := "DELETE FROM tokens WHERE user_id = ? AND platform = ?"
	_, err := r.DB.Exec(query, userID, platform)
	if err != nil {
		return fmt.Errorf("failed to delete tokens for user and platform: %w", err)
	}
	return nil
}
