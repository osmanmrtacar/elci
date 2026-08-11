package instagram

import (
	"fmt"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

const maxCaptionLen = 2200

// ValidateSettings covers what this MVP actually supports: exactly one
// image or video per post. Carousels are a real Instagram feature but
// aren't implemented yet, so multiple media URLs are rejected rather than
// silently only posting the first one.
func (p *Provider) ValidateSettings(content provider.Content, _ map[string]any, _ provider.AccountInfo) []provider.ValidationError {
	var errs []provider.ValidationError
	add := func(field, format string, a ...any) {
		errs = append(errs, provider.ValidationError{Field: field, Message: fmt.Sprintf(format, a...)})
	}

	if len(content.Caption) > maxCaptionLen {
		add("caption", "caption must be %d characters or fewer", maxCaptionLen)
	}
	if len(content.MediaURLs) != 1 {
		add("media", "exactly one image or video is required")
	}

	return errs
}
