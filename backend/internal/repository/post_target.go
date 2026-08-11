package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
)

type PostTargetRepository interface {
	Create(ctx context.Context, t domain.PostTarget) (domain.PostTarget, error)
	Get(ctx context.Context, id int64) (domain.PostTarget, error)
	ListByPost(ctx context.Context, postID int64) ([]domain.PostTarget, error)
	UpdateStatus(ctx context.Context, id int64, status domain.TargetStatus, platformPostID, errMsg string) error
	// ListProcessing feeds the status poller worker: it re-scans from the
	// database rather than tracking in-memory state, so a restart never
	// orphans a target mid-publish.
	ListProcessing(ctx context.Context) ([]domain.PostTarget, error)
}

type sqlitePostTargetRepository struct {
	db *sql.DB
}

func NewPostTargetRepository(db *sql.DB) PostTargetRepository {
	return &sqlitePostTargetRepository{db: db}
}

const postTargetColumns = `
	id, post_id, platform, caption_override, media_kind_override, media_urls_override,
	settings_json, status, platform_post_id, error, published_at, created_at, updated_at
`

func scanPostTarget(row interface{ Scan(...any) error }) (domain.PostTarget, error) {
	var t domain.PostTarget
	var captionOverride sql.NullString
	var mediaKindOverride sql.NullString
	var mediaURLsOverride sql.NullString
	var settingsJSON string
	var publishedAt sql.NullTime
	err := row.Scan(
		&t.ID, &t.PostID, &t.Platform, &captionOverride, &mediaKindOverride, &mediaURLsOverride,
		&settingsJSON, &t.Status, &t.PlatformPostID, &t.Error, &publishedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return domain.PostTarget{}, err
	}
	if captionOverride.Valid {
		t.CaptionOverride = &captionOverride.String
	}
	if mediaKindOverride.Valid {
		kind := domain.MediaKind(mediaKindOverride.String)
		t.MediaKindOverride = &kind
	}
	if mediaURLsOverride.Valid {
		if err := json.Unmarshal([]byte(mediaURLsOverride.String), &t.MediaURLsOverride); err != nil {
			return domain.PostTarget{}, err
		}
	}
	if err := json.Unmarshal([]byte(settingsJSON), &t.Settings); err != nil {
		return domain.PostTarget{}, err
	}
	if publishedAt.Valid {
		t.PublishedAt = &publishedAt.Time
	}
	return t, nil
}

func (r *sqlitePostTargetRepository) Create(ctx context.Context, t domain.PostTarget) (domain.PostTarget, error) {
	var mediaURLsOverride *string
	if t.MediaURLsOverride != nil {
		b, err := json.Marshal(t.MediaURLsOverride)
		if err != nil {
			return domain.PostTarget{}, err
		}
		s := string(b)
		mediaURLsOverride = &s
	}
	var mediaKindOverride *string
	if t.MediaKindOverride != nil {
		s := string(*t.MediaKindOverride)
		mediaKindOverride = &s
	}
	settings := t.Settings
	if settings == nil {
		settings = map[string]any{}
	}
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return domain.PostTarget{}, err
	}
	if t.Status == "" {
		t.Status = domain.TargetStatusPending
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO post_targets (post_id, platform, caption_override, media_kind_override, media_urls_override, settings_json, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, t.PostID, t.Platform, t.CaptionOverride, mediaKindOverride, mediaURLsOverride, string(settingsJSON), t.Status)
	if err != nil {
		return domain.PostTarget{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.PostTarget{}, err
	}
	return r.Get(ctx, id)
}

func (r *sqlitePostTargetRepository) Get(ctx context.Context, id int64) (domain.PostTarget, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+postTargetColumns+` FROM post_targets WHERE id = ?`, id)
	t, err := scanPostTarget(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PostTarget{}, ErrNotFound
	}
	return t, err
}

func (r *sqlitePostTargetRepository) ListByPost(ctx context.Context, postID int64) ([]domain.PostTarget, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+postTargetColumns+` FROM post_targets WHERE post_id = ? ORDER BY id`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PostTarget
	for rows.Next() {
		t, err := scanPostTarget(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *sqlitePostTargetRepository) UpdateStatus(ctx context.Context, id int64, status domain.TargetStatus, platformPostID, errMsg string) error {
	var publishedAt *time.Time
	if status == domain.TargetStatusPublished {
		now := time.Now()
		publishedAt = &now
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE post_targets
		SET status = ?, platform_post_id = ?, error = ?, published_at = COALESCE(?, published_at), updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, status, platformPostID, errMsg, publishedAt, id)
	return err
}

func (r *sqlitePostTargetRepository) ListProcessing(ctx context.Context) ([]domain.PostTarget, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+postTargetColumns+` FROM post_targets WHERE status = ?`, domain.TargetStatusProcessing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PostTarget
	for rows.Next() {
		t, err := scanPostTarget(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
