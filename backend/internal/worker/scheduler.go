// Package worker runs the background loops that turn a scheduled post into
// a published one: Scheduler dispatches due posts, StatusPoller finishes any
// still-processing target (e.g. TikTok) started by a since-restarted process.
package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

// Scheduler polls for posts whose scheduled time has arrived and publishes
// them through a small bounded worker pool, so a burst of due posts never
// spins up unbounded concurrent platform calls.
type Scheduler struct {
	posts    repository.PostRepository
	postSvc  *service.PostService
	interval time.Duration
	workers  int
	logger   *slog.Logger
}

func NewScheduler(posts repository.PostRepository, postSvc *service.PostService, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		posts:    posts,
		postSvc:  postSvc,
		interval: 30 * time.Second,
		workers:  4,
		logger:   logger,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	due, err := s.posts.ClaimDue(ctx, time.Now())
	if err != nil {
		s.logger.Error("scheduler: claim due posts", "error", err)
		return
	}
	if len(due) == 0 {
		return
	}

	jobs := make(chan domain.Post)
	go func() {
		defer close(jobs)
		for _, p := range due {
			select {
			case jobs <- p:
			case <-ctx.Done():
				return
			}
		}
	}()

	var wg sync.WaitGroup
	for range s.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if err := s.postSvc.PublishScheduled(ctx, p); err != nil {
					s.logger.Error("scheduler: publish post", "post_id", p.ID, "error", err)
				}
			}
		}()
	}
	wg.Wait()
}
