package models

import (
	"database/sql"
	"fmt"
	"time"
)

// OAuthSession stores temporary OAuth state for PKCE flow
type OAuthSession struct {
	ID           int64     `json:"id"`
	State        string    `json:"state"`
	CodeVerifier string    `json:"code_verifier"`
	Platform     Platform  `json:"platform"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type OAuthSessionRepository struct {
	DB *sql.DB
}

// NewOAuthSessionRepository creates a new OAuth session repository
func NewOAuthSessionRepository(db *sql.DB) *OAuthSessionRepository {
	return &OAuthSessionRepository{DB: db}
}

// Create creates a new OAuth session
func (r *OAuthSessionRepository) Create(session *OAuthSession) error {
	query := `
		INSERT INTO oauth_sessions (state, code_verifier, platform, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.DB.Exec(query, session.State, session.CodeVerifier, session.Platform, now, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create oauth session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	session.ID = id
	session.CreatedAt = now
	return nil
}

// GetByState retrieves an OAuth session by state
func (r *OAuthSessionRepository) GetByState(state string) (*OAuthSession, error) {
	query := `
		SELECT id, state, code_verifier, platform, created_at, expires_at
		FROM oauth_sessions WHERE state = ?
	`
	session := &OAuthSession{}
	err := r.DB.QueryRow(query, state).Scan(
		&session.ID, &session.State, &session.CodeVerifier, &session.Platform, &session.CreatedAt, &session.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("oauth session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth session: %w", err)
	}

	// Check if expired
	if time.Now().After(session.ExpiresAt) {
		r.Delete(session.ID) // Clean up expired session
		return nil, fmt.Errorf("oauth session expired")
	}

	return session, nil
}

// Delete deletes an OAuth session
func (r *OAuthSessionRepository) Delete(id int64) error {
	query := "DELETE FROM oauth_sessions WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete oauth session: %w", err)
	}
	return nil
}

// DeleteByState deletes an OAuth session by state
func (r *OAuthSessionRepository) DeleteByState(state string) error {
	query := "DELETE FROM oauth_sessions WHERE state = ?"
	_, err := r.DB.Exec(query, state)
	if err != nil {
		return fmt.Errorf("failed to delete oauth session: %w", err)
	}
	return nil
}

// CleanupExpired removes all expired OAuth sessions
func (r *OAuthSessionRepository) CleanupExpired() error {
	query := "DELETE FROM oauth_sessions WHERE expires_at < ?"
	_, err := r.DB.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired oauth sessions: %w", err)
	}
	return nil
}
