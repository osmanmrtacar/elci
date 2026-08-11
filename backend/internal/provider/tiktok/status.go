package tiktok

import (
	"context"
	"fmt"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type statusResponse struct {
	Status     string `json:"status"`
	FailReason string `json:"fail_reason"`
}

func (p *Provider) PollStatus(ctx context.Context, token provider.Token, ref provider.PublishRef) (provider.Status, error) {
	var res statusResponse
	err := postJSON(ctx, token.AccessToken, "/post/publish/status/fetch/", map[string]any{
		"publish_id": ref.ID,
	}, &res)
	if err != nil {
		return provider.Status{}, err
	}

	switch res.Status {
	case "PUBLISH_COMPLETE", "SEND_TO_USER_INBOX":
		return provider.Status{State: "published"}, nil
	case "FAILED":
		return provider.Status{State: "failed", Message: fmt.Sprintf("TikTok: %s", res.FailReason)}, nil
	default: // PROCESSING_UPLOAD, PROCESSING_DOWNLOAD
		return provider.Status{State: "processing"}, nil
	}
}
