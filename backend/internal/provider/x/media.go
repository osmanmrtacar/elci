package x

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

const chunkSize = 4 << 20 // 4MB

type mediaIDResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type processingInfo struct {
	State          string `json:"state"`
	CheckAfterSecs int    `json:"check_after_secs"`
}

type finalizeResponse struct {
	Data struct {
		ID             string          `json:"id"`
		ProcessingInfo *processingInfo `json:"processing_info"`
	} `json:"data"`
}

// uploadMedia downloads the file at mediaURL and re-uploads it to X via the
// chunked INIT/APPEND/FINALIZE flow — X's media API has no PULL_FROM_URL
// equivalent, so unlike TikTok this always round-trips the bytes.
func uploadMedia(ctx context.Context, accessToken, mediaURL string, kind provider.MediaKind) (string, error) {
	resp, err := http.Get(mediaURL)
	if err != nil {
		return "", fmt.Errorf("x: fetch media: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("x: read media: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	category := "tweet_image"
	if kind == provider.MediaVideo {
		category = "tweet_video"
	} else if contentType == "image/gif" {
		category = "tweet_gif"
	}

	var initRes mediaIDResponse
	err = bearerJSON(ctx, "POST", accessToken, "/2/media/upload/initialize", map[string]any{
		"media_type":     contentType,
		"total_bytes":    len(data),
		"media_category": category,
	}, &initRes)
	if err != nil {
		return "", fmt.Errorf("x: init media upload: %w", err)
	}
	mediaID := initRes.Data.ID

	for i, offset := 0, 0; offset < len(data); i, offset = i+1, offset+chunkSize {
		end := min(offset+chunkSize, len(data))
		if err := appendChunk(ctx, accessToken, mediaID, i, data[offset:end]); err != nil {
			return "", fmt.Errorf("x: append media chunk %d: %w", i, err)
		}
	}

	var finalizeRes finalizeResponse
	err = bearerJSON(ctx, "POST", accessToken, "/2/media/upload/"+mediaID+"/finalize", map[string]any{}, &finalizeRes)
	if err != nil {
		return "", fmt.Errorf("x: finalize media upload: %w", err)
	}

	if finalizeRes.Data.ProcessingInfo != nil {
		if err := waitForProcessing(ctx, accessToken, mediaID, *finalizeRes.Data.ProcessingInfo); err != nil {
			return "", err
		}
	}

	return mediaID, nil
}

func appendChunk(ctx context.Context, accessToken, mediaID string, segmentIndex int, chunk []byte) error {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	_ = w.WriteField("segment_index", strconv.Itoa(segmentIndex))
	part, err := w.CreateFormFile("media", "chunk")
	if err != nil {
		return err
	}
	if _, err := part.Write(chunk); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+"/2/media/upload/"+mediaID+"/append", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}
	return nil
}

func waitForProcessing(ctx context.Context, accessToken, mediaID string, info processingInfo) error {
	for range 30 {
		if info.State == "succeeded" {
			return nil
		}
		if info.State == "failed" {
			return fmt.Errorf("x: media processing failed for %s", mediaID)
		}

		wait := time.Duration(info.CheckAfterSecs) * time.Second
		if wait <= 0 {
			wait = time.Second
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}

		var statusRes struct {
			Data struct {
				ProcessingInfo processingInfo `json:"processing_info"`
			} `json:"data"`
		}
		path := fmt.Sprintf("/2/media/upload?media_id=%s&command=STATUS", mediaID)
		if err := bearerJSON(ctx, "GET", accessToken, path, nil, &statusRes); err != nil {
			return err
		}
		info = statusRes.Data.ProcessingInfo
	}
	return fmt.Errorf("x: media processing timed out for %s", mediaID)
}
