package domain

import "time"

type TargetStatus string

const (
	TargetStatusPending    TargetStatus = "pending"
	TargetStatusProcessing TargetStatus = "processing"
	TargetStatusPublished  TargetStatus = "published"
	TargetStatusFailed     TargetStatus = "failed"
)

// PostTarget is one platform a post is (or will be) published to. Overrides
// are nil unless this platform diverges from the post's default content, so
// a new per-platform field only ever grows Settings — never a migration.
type PostTarget struct {
	ID                int64
	PostID            int64
	Platform          Platform
	CaptionOverride   *string
	MediaKindOverride *MediaKind
	MediaURLsOverride []string
	Settings          map[string]any
	Status            TargetStatus
	PlatformPostID    string
	Error             string
	PublishedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Content resolves this target's effective caption/media, falling back to
// the parent post's defaults for anything not overridden.
func (t PostTarget) Content(post Post) (caption string, kind MediaKind, urls []string) {
	caption = post.DefaultCaption
	if t.CaptionOverride != nil {
		caption = *t.CaptionOverride
	}

	switch {
	case t.MediaKindOverride != nil:
		kind = *t.MediaKindOverride
	case post.DefaultMediaKind != nil:
		kind = *post.DefaultMediaKind
	}

	urls = post.DefaultMediaURLs
	if t.MediaURLsOverride != nil {
		urls = t.MediaURLsOverride
	}
	return caption, kind, urls
}
