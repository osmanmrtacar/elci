package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

// StatusPoller re-checks every still-"processing" target on a fixed
// interval. It re-scans from the database rather than tracking in-memory
// state, so a process restart never orphans a post mid-publish.
type StatusPoller struct {
	targets  repository.PostTargetRepository
	postSvc  *service.PostService
	interval time.Duration
	logger   *slog.Logger
}

func NewStatusPoller(targets repository.PostTargetRepository, postSvc *service.PostService, logger *slog.Logger) *StatusPoller {
	return &StatusPoller{
		targets:  targets,
		postSvc:  postSvc,
		interval: 20 * time.Second,
		logger:   logger,
	}
}

func (p *StatusPoller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.tick(ctx)
		}
	}
}

func (p *StatusPoller) tick(ctx context.Context) {
	processing, err := p.targets.ListProcessing(ctx)
	if err != nil {
		p.logger.Error("status poller: list processing targets", "error", err)
		return
	}
	for _, t := range processing {
		if err := p.postSvc.RefreshTargetStatus(ctx, t); err != nil {
			p.logger.Error("status poller: refresh target", "target_id", t.ID, "error", err)
		}
	}
}
