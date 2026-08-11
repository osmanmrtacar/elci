package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, p domain.Post) (domain.Post, error)
	Get(ctx context.Context, userID, id int64) (domain.Post, error)
	// GetByID skips the owner check — for internal use by workers, which
	// operate on a target's post_id without an authenticated user in scope.
	GetByID(ctx context.Context, id int64) (domain.Post, error)
	ListByUser(ctx context.Context, userID int64) ([]domain.Post, error)
	UpdateStatus(ctx context.Context, id int64, status domain.PostStatus) error
	Schedule(ctx context.Context, id int64, scheduledAt time.Time) error
	// ClaimDue atomically flips each due post from "scheduled" to
	// "publishing" one at a time, so a single scheduler tick never
	// dispatches the same post twice.
	ClaimDue(ctx context.Context, before time.Time) ([]domain.Post, error)
	Delete(ctx context.Context, userID, id int64) error
}

type sqlitePostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &sqlitePostRepository{db: db}
}

const postColumns = `
	id, user_id, default_caption, default_media_kind, default_media_urls,
	status, scheduled_at, created_at, updated_at
`

func scanPost(row interface{ Scan(...any) error }) (domain.Post, error) {
	var p domain.Post
	var mediaKind sql.NullString
	var mediaURLsJSON string
	var scheduledAt sql.NullTime
	err := row.Scan(
		&p.ID, &p.UserID, &p.DefaultCaption, &mediaKind, &mediaURLsJSON,
		&p.Status, &scheduledAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return domain.Post{}, err
	}
	if mediaKind.Valid {
		kind := domain.MediaKind(mediaKind.String)
		p.DefaultMediaKind = &kind
	}
	if scheduledAt.Valid {
		p.ScheduledAt = &scheduledAt.Time
	}
	if err := json.Unmarshal([]byte(mediaURLsJSON), &p.DefaultMediaURLs); err != nil {
		return domain.Post{}, err
	}
	return p, nil
}

func (r *sqlitePostRepository) Create(ctx context.Context, p domain.Post) (domain.Post, error) {
	mediaURLs := p.DefaultMediaURLs
	if mediaURLs == nil {
		mediaURLs = []string{}
	}
	mediaURLsJSON, err := json.Marshal(mediaURLs)
	if err != nil {
		return domain.Post{}, err
	}
	var mediaKind *string
	if p.DefaultMediaKind != nil {
		s := string(*p.DefaultMediaKind)
		mediaKind = &s
	}
	if p.Status == "" {
		p.Status = domain.PostStatusDraft
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO posts (user_id, default_caption, default_media_kind, default_media_urls, status, scheduled_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, p.UserID, p.DefaultCaption, mediaKind, string(mediaURLsJSON), p.Status, p.ScheduledAt)
	if err != nil {
		return domain.Post{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Post{}, err
	}
	return r.Get(ctx, p.UserID, id)
}

func (r *sqlitePostRepository) Get(ctx context.Context, userID, id int64) (domain.Post, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+postColumns+` FROM posts WHERE id = ? AND user_id = ?`, id, userID)
	p, err := scanPost(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Post{}, ErrNotFound
	}
	return p, err
}

func (r *sqlitePostRepository) GetByID(ctx context.Context, id int64) (domain.Post, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+postColumns+` FROM posts WHERE id = ?`, id)
	p, err := scanPost(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Post{}, ErrNotFound
	}
	return p, err
}

func (r *sqlitePostRepository) ListByUser(ctx context.Context, userID int64) ([]domain.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+postColumns+` FROM posts WHERE user_id = ? ORDER BY COALESCE(scheduled_at, created_at) DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *sqlitePostRepository) UpdateStatus(ctx context.Context, id int64, status domain.PostStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE posts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (r *sqlitePostRepository) Schedule(ctx context.Context, id int64, scheduledAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE posts SET status = ?, scheduled_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, domain.PostStatusScheduled, scheduledAt, id)
	return err
}

func (r *sqlitePostRepository) ClaimDue(ctx context.Context, before time.Time) ([]domain.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id FROM posts WHERE status = ? AND scheduled_at <= ?
	`, domain.PostStatusScheduled, before)
	if err != nil {
		return nil, err
	}
	type candidate struct{ id, userID int64 }
	var candidates []candidate
	for rows.Next() {
		var c candidate
		if err := rows.Scan(&c.id, &c.userID); err != nil {
			rows.Close()
			return nil, err
		}
		candidates = append(candidates, c)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var claimed []domain.Post
	for _, c := range candidates {
		res, err := r.db.ExecContext(ctx, `
			UPDATE posts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND status = ?
		`, domain.PostStatusPublishing, c.id, domain.PostStatusScheduled)
		if err != nil {
			return nil, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if n == 0 {
			continue
		}
		p, err := r.Get(ctx, c.userID, c.id)
		if err != nil {
			return nil, err
		}
		claimed = append(claimed, p)
	}
	return claimed, nil
}

func (r *sqlitePostRepository) Delete(ctx context.Context, userID, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ? AND user_id = ?`, id, userID)
	return err
}
