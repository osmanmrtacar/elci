package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	TikTokUserID   string    `json:"tiktok_user_id"` // Deprecated: Use PlatformUserID
	Platform       Platform  `json:"platform"`
	PlatformUserID string    `json:"platform_user_id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"display_name"`
	AvatarURL      string    `json:"avatar_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (tiktok_user_id, username, display_name, avatar_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.DB.Exec(query, user.TikTokUserID, user.Username, user.DisplayName, user.AvatarURL, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*User, error) {
	query := `
		SELECT id, tiktok_user_id, username, display_name, avatar_url, created_at, updated_at
		FROM users WHERE id = ?
	`
	user := &User{}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.TikTokUserID, &user.Username, &user.DisplayName,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByTikTokUserID retrieves a user by TikTok user ID
func (r *UserRepository) GetByTikTokUserID(tiktokUserID string) (*User, error) {
	query := `
		SELECT id, tiktok_user_id, username, display_name, avatar_url, created_at, updated_at
		FROM users WHERE tiktok_user_id = ?
	`
	user := &User{}
	err := r.DB.QueryRow(query, tiktokUserID).Scan(
		&user.ID, &user.TikTokUserID, &user.Username, &user.DisplayName,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil if not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET username = ?, display_name = ?, avatar_url = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.DB.Exec(query, user.Username, user.DisplayName, user.AvatarURL, now, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	user.UpdatedAt = now
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// CreateOrUpdate creates a new user or updates if exists
func (r *UserRepository) CreateOrUpdate(user *User) error {
	existingUser, err := r.GetByTikTokUserID(user.TikTokUserID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return r.Create(user)
	}

	user.ID = existingUser.ID
	return r.Update(user)
}

// GetByPlatformUserID retrieves a user by platform and platform user ID
func (r *UserRepository) GetByPlatformUserID(platform Platform, platformUserID string) (*User, error) {
	query := `
		SELECT id, tiktok_user_id, platform, platform_user_id, username, display_name, avatar_url, created_at, updated_at
		FROM users WHERE platform = ? AND platform_user_id = ?
	`
	user := &User{}
	err := r.DB.QueryRow(query, platform, platformUserID).Scan(
		&user.ID, &user.TikTokUserID, &user.Platform, &user.PlatformUserID, &user.Username, &user.DisplayName,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil if not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
