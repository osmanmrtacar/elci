package tiktok

import (
	"fmt"
	"slices"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

// TikTok's own limits: a video's post_info.title allows up to 2200 UTF-16
// runes; a photo's post_info.description (where we put the caption) allows
// up to 4000.
const (
	maxVideoCaptionLen = 2200
	maxPhotoCaptionLen = 4000

	// TikTok's Content Posting API rejects anything else for photo posts —
	// confirmed against the real API via file_format_check_failed after PNG
	// uploads, and picture_size_check_failed above this pixel cap.
	maxPhotoDimensionPx = 1080
)

var allowedPhotoContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/webp": true,
}

// ValidateSettings re-checks everything TikTok's Direct Post UX guidelines
// require against the account's *live* creator_info (info.Extra), so a
// stale or bypassed frontend can never publish something the account
// wouldn't actually allow.
func (p *Provider) ValidateSettings(content provider.Content, settings map[string]any, info provider.AccountInfo) []provider.ValidationError {
	var errs []provider.ValidationError
	add := func(field, format string, a ...any) {
		errs = append(errs, provider.ValidationError{Field: field, Message: fmt.Sprintf(format, a...)})
	}

	maxCaptionLen := maxVideoCaptionLen
	if content.MediaKind == provider.MediaImage {
		maxCaptionLen = maxPhotoCaptionLen
	}
	if len(content.Caption) > maxCaptionLen {
		add("caption", "caption must be %d characters or fewer", maxCaptionLen)
	}

	privacyLevel, _ := settings["privacy_level"].(string)
	allowedLevels, _ := info.Extra["privacy_level_options"].([]string)
	switch {
	case privacyLevel == "":
		add("privacy_level", "a privacy level must be selected")
	case !slices.Contains(allowedLevels, privacyLevel):
		add("privacy_level", "%q is not one of this account's currently allowed privacy levels", privacyLevel)
	}

	comment, _ := settings["comment"].(bool)
	if commentDisabled, _ := info.Extra["comment_disabled"].(bool); commentDisabled && comment {
		add("comment", "comments are disabled on this account and cannot be enabled")
	}

	if content.MediaKind == provider.MediaVideo {
		duet, _ := settings["duet"].(bool)
		if duetDisabled, _ := info.Extra["duet_disabled"].(bool); duetDisabled && duet {
			add("duet", "duet is disabled on this account and cannot be enabled")
		}
		stitch, _ := settings["stitch"].(bool)
		if stitchDisabled, _ := info.Extra["stitch_disabled"].(bool); stitchDisabled && stitch {
			add("stitch", "stitch is disabled on this account and cannot be enabled")
		}

		if maxDuration, ok := info.Extra["max_video_post_duration_sec"].(int); ok && maxDuration > 0 && content.MediaDurationSec != nil {
			if *content.MediaDurationSec > maxDuration {
				add("media", "video is longer than the %ds this account allows", maxDuration)
			}
		}
	}

	if content.MediaKind == provider.MediaImage {
		for i := range content.MediaURLs {
			ct := ""
			if i < len(content.MediaContentTypes) {
				ct = content.MediaContentTypes[i]
			}
			if ct != "" && !allowedPhotoContentTypes[ct] {
				add("media", "image %d is %s — TikTok only accepts JPEG or WebP for photo posts", i+1, ct)
			}

			w, h := 0, 0
			if i < len(content.MediaWidths) {
				w = content.MediaWidths[i]
			}
			if i < len(content.MediaHeights) {
				h = content.MediaHeights[i]
			}
			if w > maxPhotoDimensionPx || h > maxPhotoDimensionPx {
				add("media", "image %d is %dx%d — TikTok caps photo dimensions at %dp", i+1, w, h, maxPhotoDimensionPx)
			}
		}
	}

	discloseContent, _ := settings["disclose_content"].(bool)
	yourBrand, _ := settings["your_brand"].(bool)
	brandedContent, _ := settings["branded_content"].(bool)

	if discloseContent && !yourBrand && !brandedContent {
		add("disclose_content", "select at least one of \"your brand\" or \"branded content\"")
	}
	if brandedContent && privacyLevel == "SELF_ONLY" {
		add("privacy_level", "branded content cannot be private — choose a public or friends-only visibility")
	}

	return errs
}
