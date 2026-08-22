package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/provider"
	"github.com/osmanmertacar/elci/backend/internal/repository"
)

var ErrNoTargets = errors.New("post must have at least one target platform")

// ValidationFailure carries every target's rejection reasons at once, so a
// caller can surface all of them rather than just the first.
type ValidationFailure struct {
	ByTarget map[domain.Platform][]provider.ValidationError
}

func (e *ValidationFailure) Error() string {
	return fmt.Sprintf("validation failed for %d target(s)", len(e.ByTarget))
}

type PostService struct {
	registry    *provider.Registry
	posts       repository.PostRepository
	targets     repository.PostTargetRepository
	connections repository.ConnectionRepository
	media       repository.MediaRepository
}

func NewPostService(
	registry *provider.Registry,
	posts repository.PostRepository,
	targets repository.PostTargetRepository,
	connections repository.ConnectionRepository,
	media repository.MediaRepository,
) *PostService {
	return &PostService{registry: registry, posts: posts, targets: targets, connections: connections, media: media}
}

type TargetInput struct {
	Platform          domain.Platform
	CaptionOverride   *string
	MediaKindOverride *domain.MediaKind
	MediaURLsOverride []string
	Settings          map[string]any
}

type CreatePostInput struct {
	DefaultCaption   string
	DefaultMediaKind *domain.MediaKind
	DefaultMediaURLs []string
	Targets          []TargetInput
}

// CreateDraft saves a new post and its targets in draft status. Nothing is
// checked against live platform rules until Schedule or PublishNow.
func (s *PostService) CreateDraft(ctx context.Context, userID int64, input CreatePostInput) (domain.Post, []domain.PostTarget, error) {
	if len(input.Targets) == 0 {
		return domain.Post{}, nil, ErrNoTargets
	}

	post, err := s.posts.Create(ctx, domain.Post{
		UserID:           userID,
		DefaultCaption:   input.DefaultCaption,
		DefaultMediaKind: input.DefaultMediaKind,
		DefaultMediaURLs: input.DefaultMediaURLs,
		Status:           domain.PostStatusDraft,
	})
	if err != nil {
		return domain.Post{}, nil, err
	}

	targets := make([]domain.PostTarget, 0, len(input.Targets))
	for _, ti := range input.Targets {
		t, err := s.targets.Create(ctx, domain.PostTarget{
			PostID:            post.ID,
			Platform:          ti.Platform,
			CaptionOverride:   ti.CaptionOverride,
			MediaKindOverride: ti.MediaKindOverride,
			MediaURLsOverride: ti.MediaURLsOverride,
			Settings:          ti.Settings,
		})
		if err != nil {
			return domain.Post{}, nil, err
		}
		targets = append(targets, t)
	}

	return post, targets, nil
}

// Schedule validates every target against the platform's live rules, then
// marks the post scheduled for the scheduler worker to pick up.
func (s *PostService) Schedule(ctx context.Context, userID, postID int64, scheduledAt time.Time) error {
	post, targets, err := s.loadForValidation(ctx, userID, postID)
	if err != nil {
		return err
	}
	if err := s.validateAll(ctx, userID, post, targets); err != nil {
		return err
	}
	return s.posts.Schedule(ctx, postID, scheduledAt)
}

// PublishNow validates then publishes immediately, without waiting for the
// scheduler's next tick.
func (s *PostService) PublishNow(ctx context.Context, userID, postID int64) error {
	post, targets, err := s.loadForValidation(ctx, userID, postID)
	if err != nil {
		return err
	}
	if err := s.validateAll(ctx, userID, post, targets); err != nil {
		return err
	}
	if err := s.posts.UpdateStatus(ctx, postID, domain.PostStatusPublishing); err != nil {
		return err
	}
	return s.PublishPost(ctx, post, targets)
}

// PublishScheduled publishes a post the scheduler has already claimed (its
// status is "publishing" and its targets were validated at Schedule time).
func (s *PostService) PublishScheduled(ctx context.Context, post domain.Post) error {
	targets, err := s.targets.ListByPost(ctx, post.ID)
	if err != nil {
		return err
	}
	return s.PublishPost(ctx, post, targets)
}

func (s *PostService) loadForValidation(ctx context.Context, userID, postID int64) (domain.Post, []domain.PostTarget, error) {
	post, err := s.posts.Get(ctx, userID, postID)
	if err != nil {
		return domain.Post{}, nil, err
	}
	targets, err := s.targets.ListByPost(ctx, post.ID)
	if err != nil {
		return domain.Post{}, nil, err
	}
	if len(targets) == 0 {
		return domain.Post{}, nil, ErrNoTargets
	}
	return post, targets, nil
}

// validateAll re-fetches each target's *live* account info and re-checks its
// settings server-side — the frontend's own checks are never trusted alone.
func (s *PostService) validateAll(ctx context.Context, userID int64, post domain.Post, targets []domain.PostTarget) error {
	failures := ValidationFailure{ByTarget: map[domain.Platform][]provider.ValidationError{}}

	for _, t := range targets {
		p, ok := s.registry.Get(string(t.Platform))
		if !ok {
			failures.ByTarget[t.Platform] = []provider.ValidationError{{Field: "platform", Message: "unsupported platform"}}
			continue
		}

		conn, err := s.connections.Get(ctx, userID, t.Platform)
		if err != nil {
			failures.ByTarget[t.Platform] = []provider.ValidationError{{Field: "platform", Message: "not connected"}}
			continue
		}

		token, err := refreshTokenIfNeeded(ctx, p, s.connections, conn)
		if err != nil {
			failures.ByTarget[t.Platform] = []provider.ValidationError{{Field: "platform", Message: "token refresh failed: " + err.Error()}}
			continue
		}

		info, err := p.AccountInfo(ctx, token)
		if err != nil {
			failures.ByTarget[t.Platform] = []provider.ValidationError{{Field: "platform", Message: "could not load account info: " + err.Error()}}
			continue
		}

		if errs := p.ValidateSettings(s.targetContent(ctx, post, t), t.Settings, info); len(errs) > 0 {
			failures.ByTarget[t.Platform] = errs
		}
	}

	if len(failures.ByTarget) > 0 {
		return &failures
	}
	return nil
}

// PublishPost publishes every target of an already-validated post — the one
// code path shared by PublishNow and the scheduler.
func (s *PostService) PublishPost(ctx context.Context, post domain.Post, targets []domain.PostTarget) error {
	for _, t := range targets {
		s.publishTarget(ctx, post, t)
	}
	return s.recomputePostStatus(ctx, post.ID)
}

func (s *PostService) publishTarget(ctx context.Context, post domain.Post, t domain.PostTarget) {
	p, ok := s.registry.Get(string(t.Platform))
	if !ok {
		_ = s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusFailed, "", "unsupported platform")
		return
	}

	conn, err := s.connections.Get(ctx, post.UserID, t.Platform)
	if err != nil {
		_ = s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusFailed, "", "not connected: "+err.Error())
		return
	}

	token, err := refreshTokenIfNeeded(ctx, p, s.connections, conn)
	if err != nil {
		_ = s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusFailed, "", "token refresh failed: "+err.Error())
		return
	}

	ref, err := p.Publish(ctx, token, s.targetContent(ctx, post, t), t.Settings)
	if err != nil {
		_ = s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusFailed, "", err.Error())
		return
	}

	// X and Instagram are already done by the time Publish returns; TikTok
	// usually still needs the status poller worker to finish the job.
	status, err := p.PollStatus(ctx, token, ref)
	if err != nil {
		_ = s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusProcessing, ref.ID, "")
		return
	}
	s.applyStatus(ctx, t.ID, ref.ID, status)
}

// RefreshTargetStatus re-polls one still-processing target and, if it has
// resolved, updates the target and recomputes the parent post's status.
// Used by the status poller worker, so a restart never leaves a post stuck.
func (s *PostService) RefreshTargetStatus(ctx context.Context, t domain.PostTarget) error {
	p, ok := s.registry.Get(string(t.Platform))
	if !ok {
		return s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusFailed, t.PlatformPostID, "unsupported platform")
	}

	post, err := s.posts.GetByID(ctx, t.PostID)
	if err != nil {
		return err
	}

	conn, err := s.connections.Get(ctx, post.UserID, t.Platform)
	if err != nil {
		return s.targets.UpdateStatus(ctx, t.ID, domain.TargetStatusFailed, t.PlatformPostID, "not connected: "+err.Error())
	}

	token, err := refreshTokenIfNeeded(ctx, p, s.connections, conn)
	if err != nil {
		return err
	}

	status, err := p.PollStatus(ctx, token, provider.PublishRef{ID: t.PlatformPostID})
	if err != nil {
		return err
	}
	if status.State == "processing" {
		return nil
	}

	s.applyStatus(ctx, t.ID, t.PlatformPostID, status)
	return s.recomputePostStatus(ctx, post.ID)
}

func (s *PostService) applyStatus(ctx context.Context, targetID int64, platformPostID string, status provider.Status) {
	switch status.State {
	case "published":
		_ = s.targets.UpdateStatus(ctx, targetID, domain.TargetStatusPublished, platformPostID, "")
	case "failed":
		_ = s.targets.UpdateStatus(ctx, targetID, domain.TargetStatusFailed, platformPostID, status.Message)
	default:
		_ = s.targets.UpdateStatus(ctx, targetID, domain.TargetStatusProcessing, platformPostID, "")
	}
}

func (s *PostService) recomputePostStatus(ctx context.Context, postID int64) error {
	targets, err := s.targets.ListByPost(ctx, postID)
	if err != nil {
		return err
	}

	var published, failed, pending int
	for _, t := range targets {
		switch t.Status {
		case domain.TargetStatusPublished:
			published++
		case domain.TargetStatusFailed:
			failed++
		default:
			pending++
		}
	}

	status := domain.PostStatusPublishing
	switch {
	case pending > 0:
		status = domain.PostStatusPublishing
	case failed == 0:
		status = domain.PostStatusPublished
	case published == 0:
		status = domain.PostStatusFailed
	default:
		status = domain.PostStatusPartial
	}

	return s.posts.UpdateStatus(ctx, postID, status)
}

// targetContent looks up each media URL's stored asset (uploads are recorded
// at confirm time, keyed by the same user) so providers can validate format
// and dimensions against what was actually uploaded, not just trust the URL.
// A lookup miss (asset deleted, or not one of this user's own uploads) just
// leaves that entry's content type/dimensions at zero — unknown, not passing.
func (s *PostService) targetContent(ctx context.Context, post domain.Post, t domain.PostTarget) provider.Content {
	caption, kind, urls := t.Content(post)

	contentTypes := make([]string, len(urls))
	widths := make([]int, len(urls))
	heights := make([]int, len(urls))
	var durationSec *int
	for i, url := range urls {
		asset, err := s.media.GetByPublicURL(ctx, post.UserID, url)
		if err != nil {
			continue
		}
		contentTypes[i] = asset.ContentType
		if asset.Width != nil {
			widths[i] = *asset.Width
		}
		if asset.Height != nil {
			heights[i] = *asset.Height
		}
		if asset.DurationSeconds != nil {
			durationSec = asset.DurationSeconds
		}
	}

	return provider.Content{
		Caption:           caption,
		MediaKind:         provider.MediaKind(kind),
		MediaURLs:         urls,
		MediaContentTypes: contentTypes,
		MediaWidths:       widths,
		MediaHeights:      heights,
		MediaDurationSec:  durationSec,
	}
}
