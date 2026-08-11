package domain

import "time"

type PostStatus string

const (
	PostStatusDraft      PostStatus = "draft"
	PostStatusScheduled  PostStatus = "scheduled"
	PostStatusPublishing PostStatus = "publishing"
	PostStatusPublished  PostStatus = "published"
	PostStatusPartial    PostStatus = "partial"
	PostStatusFailed     PostStatus = "failed"
)

// Post is the platform-agnostic draft: its caption/media are the defaults
// every target falls back to unless that target overrides them.
type Post struct {
	ID               int64
	UserID           int64
	DefaultCaption   string
	DefaultMediaKind *MediaKind
	DefaultMediaURLs []string
	Status           PostStatus
	ScheduledAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
