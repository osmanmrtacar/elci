package instagram

import (
	"fmt"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

const maxCaptionLen = 2200

// maxCarouselItems is Instagram's own limit on how many images a single
// carousel post may contain.
const maxCarouselItems = 10

// ValidateSettings allows exactly one video, or one to maxCarouselItems
// images — a single image publishes directly, more than one publishes as a
// carousel (see Publish).
func (p *Provider) ValidateSettings(content provider.Content, _ map[string]any, _ provider.AccountInfo) []provider.ValidationError {
	var errs []provider.ValidationError
	add := func(field, format string, a ...any) {
		errs = append(errs, provider.ValidationError{Field: field, Message: fmt.Sprintf(format, a...)})
	}

	if len(content.Caption) > maxCaptionLen {
		add("caption", "caption must be %d characters or fewer", maxCaptionLen)
	}

	switch {
	case content.MediaKind == provider.MediaVideo && len(content.MediaURLs) != 1:
		add("media", "exactly one video is required")
	case content.MediaKind == provider.MediaImage && len(content.MediaURLs) == 0:
		add("media", "at least one image is required")
	case content.MediaKind == provider.MediaImage && len(content.MediaURLs) > maxCarouselItems:
		add("media", "a carousel post allows at most %d images", maxCarouselItems)
	}

	return errs
}
