package models

import (
	"database/sql"
	"fmt"
	"time"
)

type PostStatus string

const (
	PostStatusPending    PostStatus = "pending"
	PostStatusProcessing PostStatus = "processing"
	PostStatusPublished  PostStatus = "published"
	PostStatusFailed     PostStatus = "failed"
)

type Post struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	Platform       Platform   `json:"platform"`
	TikTokPostID   string     `json:"tiktok_post_id,omitempty"` // Deprecated: Use PlatformPostID
	PlatformPostID string     `json:"platform_post_id,omitempty"`
	VideoURL       string     `json:"video_url"`
	Caption        string     `json:"caption"`
	MediaType      string     `json:"media_type"` // video, image, text
	Status         PostStatus `json:"status"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`
}

type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// Create creates a new post
func (r *PostRepository) Create(post *Post) error {
	query := `
		INSERT INTO posts (user_id, video_url, caption, status, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.DB.Exec(query, post.UserID, post.VideoURL, post.Caption, post.Status, now)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	post.ID = id
	post.CreatedAt = now
	return nil
}

// GetByID retrieves a post by ID
func (r *PostRepository) GetByID(id int64) (*Post, error) {
	query := `
		SELECT id, user_id, tiktok_post_id, video_url, caption, status, error_message, created_at, published_at
		FROM posts WHERE id = ?
	`
	post := &Post{}
	var tiktokPostID, errorMessage sql.NullString
	var publishedAt sql.NullTime

	err := r.DB.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &tiktokPostID, &post.VideoURL, &post.Caption,
		&post.Status, &errorMessage, &post.CreatedAt, &publishedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if tiktokPostID.Valid {
		post.TikTokPostID = tiktokPostID.String
	}
	if errorMessage.Valid {
		post.ErrorMessage = errorMessage.String
	}
	if publishedAt.Valid {
		post.PublishedAt = &publishedAt.Time
	}

	return post, nil
}

// GetByUserID retrieves all posts for a user
func (r *PostRepository) GetByUserID(userID int64, limit, offset int) ([]*Post, error) {
	query := `
		SELECT id, user_id, tiktok_post_id, video_url, caption, status, error_message, created_at, published_at
		FROM posts
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		var tiktokPostID, errorMessage sql.NullString
		var publishedAt sql.NullTime

		err := rows.Scan(
			&post.ID, &post.UserID, &tiktokPostID, &post.VideoURL, &post.Caption,
			&post.Status, &errorMessage, &post.CreatedAt, &publishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		if tiktokPostID.Valid {
			post.TikTokPostID = tiktokPostID.String
		}
		if errorMessage.Valid {
			post.ErrorMessage = errorMessage.String
		}
		if publishedAt.Valid {
			post.PublishedAt = &publishedAt.Time
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

// UpdateStatus updates the status of a post
func (r *PostRepository) UpdateStatus(id int64, status PostStatus, errorMessage string) error {
	query := `
		UPDATE posts
		SET status = ?, error_message = ?
		WHERE id = ?
	`
	_, err := r.DB.Exec(query, status, errorMessage, id)
	if err != nil {
		return fmt.Errorf("failed to update post status: %w", err)
	}
	return nil
}

// MarkPublished marks a post as published
func (r *PostRepository) MarkPublished(id int64, tiktokPostID string) error {
	query := `
		UPDATE posts
		SET status = ?, tiktok_post_id = ?, published_at = ?, error_message = NULL
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.DB.Exec(query, PostStatusPublished, tiktokPostID, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark post as published: %w", err)
	}
	return nil
}

// Delete deletes a post
func (r *PostRepository) Delete(id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// CountByUserID counts posts for a user
func (r *PostRepository) CountByUserID(userID int64) (int, error) {
	query := "SELECT COUNT(*) FROM posts WHERE user_id = ?"
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}
	return count, nil
}

// GetByUserIDAndPlatform retrieves all posts for a user and specific platform
func (r *PostRepository) GetByUserIDAndPlatform(userID int64, platform Platform, limit, offset int) ([]*Post, error) {
	query := `
		SELECT id, user_id, platform, tiktok_post_id, platform_post_id, video_url, caption, media_type, status, error_message, created_at, published_at
		FROM posts
		WHERE user_id = ? AND platform = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.DB.Query(query, userID, platform, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		var tiktokPostID, platformPostID, errorMessage sql.NullString
		var publishedAt sql.NullTime

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Platform, &tiktokPostID, &platformPostID, &post.VideoURL, &post.Caption,
			&post.MediaType, &post.Status, &errorMessage, &post.CreatedAt, &publishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		if tiktokPostID.Valid {
			post.TikTokPostID = tiktokPostID.String
		}
		if platformPostID.Valid {
			post.PlatformPostID = platformPostID.String
		}
		if errorMessage.Valid {
			post.ErrorMessage = errorMessage.String
		}
		if publishedAt.Valid {
			post.PublishedAt = &publishedAt.Time
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

// MarkPublishedWithPlatform marks a post as published with platform-specific post ID
func (r *PostRepository) MarkPublishedWithPlatform(id int64, platformPostID string) error {
	query := `
		UPDATE posts
		SET status = ?, platform_post_id = ?, published_at = ?, error_message = NULL
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.DB.Exec(query, PostStatusPublished, platformPostID, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark post as published: %w", err)
	}
	return nil
}
