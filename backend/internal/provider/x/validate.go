package x

import (
	"fmt"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

// maxCaptionLen is X's standard post length. Premium subscribers get a much
// higher limit, but that isn't reliably detectable from the scopes this app
// requests, so every account is held to the safe default.
const maxCaptionLen = 280

func (p *Provider) ValidateSettings(content provider.Content, _ map[string]any, _ provider.AccountInfo) []provider.ValidationError {
	var errs []provider.ValidationError
	add := func(field, format string, a ...any) {
		errs = append(errs, provider.ValidationError{Field: field, Message: fmt.Sprintf(format, a...)})
	}

	if len(content.Caption) > maxCaptionLen {
		add("caption", "post must be %d characters or fewer", maxCaptionLen)
	}

	switch content.MediaKind {
	case provider.MediaVideo:
		if len(content.MediaURLs) > 1 {
			add("media", "only one video is allowed per post")
		}
	case provider.MediaImage:
		if len(content.MediaURLs) > 4 {
			add("media", "only up to 4 images are allowed per post")
		}
	}

	return errs
}
