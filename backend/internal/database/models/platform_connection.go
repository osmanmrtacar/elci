package models

import (
	"database/sql"
	"fmt"
	"time"
)

// PlatformConnection represents a user's connection to a social media platform
type PlatformConnection struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Platform       Platform  `json:"platform"`
	PlatformUserID string    `json:"platform_user_id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"display_name"`
	AvatarURL      string    `json:"avatar_url"`
	IsActive       bool      `json:"is_active"`
	ConnectedAt    time.Time `json:"connected_at"`
	LastUsedAt     time.Time `json:"last_used_at"`
}

type PlatformConnectionRepository struct {
	DB *sql.DB
}

// NewPlatformConnectionRepository creates a new platform connection repository
func NewPlatformConnectionRepository(db *sql.DB) *PlatformConnectionRepository {
	return &PlatformConnectionRepository{DB: db}
}

// Create creates a new platform connection
func (r *PlatformConnectionRepository) Create(conn *PlatformConnection) error {
	query := `
		INSERT INTO platform_connections (user_id, platform, platform_user_id, username, display_name, avatar_url, is_active, connected_at, last_used_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.DB.Exec(query, conn.UserID, conn.Platform, conn.PlatformUserID, conn.Username, conn.DisplayName, conn.AvatarURL, conn.IsActive, now, now)
	if err != nil {
		return fmt.Errorf("failed to create platform connection: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	conn.ID = id
	conn.ConnectedAt = now
	conn.LastUsedAt = now
	return nil
}

// GetByUserIDAndPlatform retrieves a platform connection by user ID and platform
func (r *PlatformConnectionRepository) GetByUserIDAndPlatform(userID int64, platform Platform) (*PlatformConnection, error) {
	query := `
		SELECT id, user_id, platform, platform_user_id, username, display_name, avatar_url, is_active, connected_at, last_used_at
		FROM platform_connections WHERE user_id = ? AND platform = ?
	`
	conn := &PlatformConnection{}
	err := r.DB.QueryRow(query, userID, platform).Scan(
		&conn.ID, &conn.UserID, &conn.Platform, &conn.PlatformUserID, &conn.Username, &conn.DisplayName,
		&conn.AvatarURL, &conn.IsActive, &conn.ConnectedAt, &conn.LastUsedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil if not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get platform connection: %w", err)
	}
	return conn, nil
}

// GetByUserID retrieves all platform connections for a user
func (r *PlatformConnectionRepository) GetByUserID(userID int64) ([]*PlatformConnection, error) {
	query := `
		SELECT id, user_id, platform, platform_user_id, username, display_name, avatar_url, is_active, connected_at, last_used_at
		FROM platform_connections WHERE user_id = ? ORDER BY connected_at DESC
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform connections: %w", err)
	}
	defer rows.Close()

	var connections []*PlatformConnection
	for rows.Next() {
		conn := &PlatformConnection{}
		err := rows.Scan(
			&conn.ID, &conn.UserID, &conn.Platform, &conn.PlatformUserID, &conn.Username, &conn.DisplayName,
			&conn.AvatarURL, &conn.IsActive, &conn.ConnectedAt, &conn.LastUsedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan platform connection: %w", err)
		}
		connections = append(connections, conn)
	}

	return connections, nil
}

// Update updates a platform connection
func (r *PlatformConnectionRepository) Update(conn *PlatformConnection) error {
	query := `
		UPDATE platform_connections
		SET username = ?, display_name = ?, avatar_url = ?, is_active = ?, last_used_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.DB.Exec(query, conn.Username, conn.DisplayName, conn.AvatarURL, conn.IsActive, now, conn.ID)
	if err != nil {
		return fmt.Errorf("failed to update platform connection: %w", err)
	}
	conn.LastUsedAt = now
	return nil
}

// UpdateLastUsed updates the last used timestamp
func (r *PlatformConnectionRepository) UpdateLastUsed(userID int64, platform Platform) error {
	query := `
		UPDATE platform_connections
		SET last_used_at = ?
		WHERE user_id = ? AND platform = ?
	`
	_, err := r.DB.Exec(query, time.Now(), userID, platform)
	if err != nil {
		return fmt.Errorf("failed to update last used: %w", err)
	}
	return nil
}

// Deactivate deactivates a platform connection
func (r *PlatformConnectionRepository) Deactivate(userID int64, platform Platform) error {
	query := `
		UPDATE platform_connections
		SET is_active = ?
		WHERE user_id = ? AND platform = ?
	`
	_, err := r.DB.Exec(query, false, userID, platform)
	if err != nil {
		return fmt.Errorf("failed to deactivate platform connection: %w", err)
	}
	return nil
}

// Delete deletes a platform connection
func (r *PlatformConnectionRepository) Delete(userID int64, platform Platform) error {
	query := "DELETE FROM platform_connections WHERE user_id = ? AND platform = ?"
	_, err := r.DB.Exec(query, userID, platform)
	if err != nil {
		return fmt.Errorf("failed to delete platform connection: %w", err)
	}
	return nil
}

// CreateOrUpdate creates a new platform connection or updates if exists
func (r *PlatformConnectionRepository) CreateOrUpdate(conn *PlatformConnection) error {
	existing, err := r.GetByUserIDAndPlatform(conn.UserID, conn.Platform)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(conn)
	}

	conn.ID = existing.ID
	conn.ConnectedAt = existing.ConnectedAt
	return r.Update(conn)
}
