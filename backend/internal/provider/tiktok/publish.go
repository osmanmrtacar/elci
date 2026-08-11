package tiktok

import (
	"context"
	"errors"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type initResponse struct {
	PublishID string `json:"publish_id"`
}

// Publish always uses PULL_FROM_URL: by the time a post reaches this point
// its media already lives in our own (ownership-verified) object storage,
// which is exactly the case TikTok's own guidelines say PULL_FROM_URL is
// for — re-uploading the bytes through FILE_UPLOAD would be redundant.
func (p *Provider) Publish(ctx context.Context, token provider.Token, content provider.Content, settings map[string]any) (provider.PublishRef, error) {
	if len(content.MediaURLs) == 0 {
		return provider.PublishRef{}, errors.New("tiktok: at least one media URL is required")
	}

	privacyLevel, _ := settings["privacy_level"].(string)
	comment, _ := settings["comment"].(bool)
	brandedContent, _ := settings["branded_content"].(bool)
	yourBrand, _ := settings["your_brand"].(bool)

	var res initResponse
	var err error

	if content.MediaKind == provider.MediaVideo {
		duet, _ := settings["duet"].(bool)
		stitch, _ := settings["stitch"].(bool)

		err = postJSON(ctx, token.AccessToken, "/post/publish/video/init/", map[string]any{
			"post_info": map[string]any{
				"title":                content.Caption,
				"privacy_level":        privacyLevel,
				"disable_duet":         !duet,
				"disable_stitch":       !stitch,
				"disable_comment":      !comment,
				"brand_content_toggle": brandedContent,
				"brand_organic_toggle": yourBrand,
			},
			"source_info": map[string]any{
				"source":    "PULL_FROM_URL",
				"video_url": content.MediaURLs[0],
			},
		}, &res)
	} else {
		autoAddMusic, _ := settings["auto_add_music"].(bool)

		err = postJSON(ctx, token.AccessToken, "/post/publish/content/init/", map[string]any{
			"media_type": "PHOTO",
			"post_mode":  "DIRECT_POST",
			"post_info": map[string]any{
				"description":          content.Caption,
				"privacy_level":        privacyLevel,
				"disable_comment":      !comment,
				"auto_add_music":       autoAddMusic,
				"brand_content_toggle": brandedContent,
				"brand_organic_toggle": yourBrand,
			},
			"source_info": map[string]any{
				"source":            "PULL_FROM_URL",
				"photo_images":      content.MediaURLs,
				"photo_cover_index": 0,
			},
		}, &res)
	}

	if err != nil {
		return provider.PublishRef{}, friendlyPublishError(err)
	}

	return provider.PublishRef{ID: res.PublishID}, nil
}
