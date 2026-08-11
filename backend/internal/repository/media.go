package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/osmanmertacar/elci/backend/internal/domain"
)

type MediaRepository interface {
	Create(ctx context.Context, m domain.MediaAsset) (domain.MediaAsset, error)
	Get(ctx context.Context, userID, id int64) (domain.MediaAsset, error)
}

type sqliteMediaRepository struct {
	db *sql.DB
}

func NewMediaRepository(db *sql.DB) MediaRepository {
	return &sqliteMediaRepository{db: db}
}

func (r *sqliteMediaRepository) Create(ctx context.Context, m domain.MediaAsset) (domain.MediaAsset, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO media_assets (user_id, kind, storage_key, public_url, content_type, size_bytes, duration_seconds)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, m.UserID, m.Kind, m.StorageKey, m.PublicURL, m.ContentType, m.SizeBytes, m.DurationSeconds)
	if err != nil {
		return domain.MediaAsset{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.MediaAsset{}, err
	}
	return r.Get(ctx, m.UserID, id)
}

func (r *sqliteMediaRepository) Get(ctx context.Context, userID, id int64) (domain.MediaAsset, error) {
	var m domain.MediaAsset
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, kind, storage_key, public_url, content_type, size_bytes, duration_seconds, created_at
		FROM media_assets WHERE id = ? AND user_id = ?
	`, id, userID).Scan(
		&m.ID, &m.UserID, &m.Kind, &m.StorageKey, &m.PublicURL, &m.ContentType, &m.SizeBytes, &m.DurationSeconds, &m.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.MediaAsset{}, ErrNotFound
	}
	return m, err
}
