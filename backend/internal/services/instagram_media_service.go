package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// InstagramMediaService handles Instagram media upload and publishing
type InstagramMediaService struct {
	httpClient *http.Client
}

// NewInstagramMediaService creates a new Instagram media service
func NewInstagramMediaService() *InstagramMediaService {
	return &InstagramMediaService{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// CreateMediaContainerResponse represents the response from creating a media container
type CreateMediaContainerResponse struct {
	ID string `json:"id"`
}

// MediaStatusResponse represents the status of a media container
type MediaStatusResponse struct {
	ID         string `json:"id"`
	StatusCode string `json:"status_code"`
}

// PublishMediaResponse represents the response from publishing media
type PublishMediaResponse struct {
	ID string `json:"id"`
}

// PermalinkResponse represents the permalink response
type PermalinkResponse struct {
	ID        string `json:"id"`
	Permalink string `json:"permalink"`
}

// CreateMediaContainer creates a container for Instagram Reels/Stories
// This is step 1 of the publishing process
func (s *InstagramMediaService) CreateMediaContainer(
	accessToken string,
	igUserID string,
	videoURL string,
	caption string,
	mediaType string, // "REELS" or "STORIES"
) (string, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/%s/media", igUserID)

	params := url.Values{}
	params.Set("media_type", mediaType)
	params.Set("video_url", videoURL)
	if caption != "" {
		params.Set("caption", caption)
	}
	params.Set("access_token", accessToken)

	resp, err := s.httpClient.PostForm(apiURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to create media container: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("create media container failed (%d): %s", resp.StatusCode, string(body))
	}

	var containerResp CreateMediaContainerResponse
	if err := json.Unmarshal(body, &containerResp); err != nil {
		return "", fmt.Errorf("failed to parse container response: %w", err)
	}

	return containerResp.ID, nil
}

// CreatePhotoContainer creates a container for Instagram photo posts
// Photos are processed immediately and don't require status polling
func (s *InstagramMediaService) CreatePhotoContainer(
	accessToken string,
	igUserID string,
	imageURL string,
	caption string,
) (string, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/%s/media", igUserID)

	params := url.Values{}
	params.Set("image_url", imageURL)
	if caption != "" {
		params.Set("caption", caption)
	}
	params.Set("access_token", accessToken)

	resp, err := s.httpClient.PostForm(apiURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to create photo container: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("create photo container failed (%d): %s", resp.StatusCode, string(body))
	}

	var containerResp CreateMediaContainerResponse
	if err := json.Unmarshal(body, &containerResp); err != nil {
		return "", fmt.Errorf("failed to parse container response: %w", err)
	}

	return containerResp.ID, nil
}

// CheckMediaStatus checks the upload status of a media container
// Returns the status code (FINISHED, IN_PROGRESS, ERROR, etc.)
func (s *InstagramMediaService) CheckMediaStatus(accessToken string, containerID string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/%s?fields=status_code&access_token=%s",
		containerID, accessToken)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to check media status: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("check media status failed (%d): %s", resp.StatusCode, string(body))
	}

	var statusResp MediaStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return "", fmt.Errorf("failed to parse status response: %w", err)
	}

	return statusResp.StatusCode, nil
}

// WaitForMediaProcessing polls the media status until it's finished or times out
func (s *InstagramMediaService) WaitForMediaProcessing(
	accessToken string,
	containerID string,
	maxWaitSeconds int,
) (bool, error) {
	startTime := time.Now()
	checkInterval := 5 * time.Second

	for {
		// Check if we've exceeded max wait time
		if time.Since(startTime).Seconds() > float64(maxWaitSeconds) {
			return false, fmt.Errorf("media processing timeout after %d seconds", maxWaitSeconds)
		}

		status, err := s.CheckMediaStatus(accessToken, containerID)
		if err != nil {
			return false, err
		}

		switch status {
		case "FINISHED":
			return true, nil
		case "ERROR":
			return false, fmt.Errorf("media processing failed with status: ERROR")
		case "IN_PROGRESS":
			// Continue waiting
			time.Sleep(checkInterval)
		default:
			// Unknown status, continue waiting
			time.Sleep(checkInterval)
		}
	}
}

// PublishMedia publishes a media container to Instagram
// This is step 3 of the publishing process (after creation and processing)
func (s *InstagramMediaService) PublishMedia(
	accessToken string,
	igUserID string,
	containerID string,
) (string, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/%s/media_publish", igUserID)

	params := url.Values{}
	params.Set("creation_id", containerID)
	params.Set("access_token", accessToken)

	resp, err := s.httpClient.PostForm(apiURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to publish media: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("publish media failed (%d): %s", resp.StatusCode, string(body))
	}

	var publishResp PublishMediaResponse
	if err := json.Unmarshal(body, &publishResp); err != nil {
		return "", fmt.Errorf("failed to parse publish response: %w", err)
	}

	return publishResp.ID, nil
}

// GetPermalink retrieves the permalink (share URL) for a published media
func (s *InstagramMediaService) GetPermalink(accessToken string, mediaID string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/%s?fields=permalink&access_token=%s",
		mediaID, accessToken)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to get permalink: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get permalink failed (%d): %s", resp.StatusCode, string(body))
	}

	var permalinkResp PermalinkResponse
	if err := json.Unmarshal(body, &permalinkResp); err != nil {
		return "", fmt.Errorf("failed to parse permalink response: %w", err)
	}

	return permalinkResp.Permalink, nil
}

// UploadAndPublishReel is a complete flow for uploading and publishing a reel
func (s *InstagramMediaService) UploadAndPublishReel(
	accessToken string,
	igUserID string,
	videoURL string,
	caption string,
) (string, string, error) {
	// Step 1: Create media container
	containerID, err := s.CreateMediaContainer(accessToken, igUserID, videoURL, caption, "REELS")
	if err != nil {
		return "", "", fmt.Errorf("create container failed: %w", err)
	}

	// Step 2: Wait for processing (max 5 minutes)
	success, err := s.WaitForMediaProcessing(accessToken, containerID, 300)
	if err != nil {
		return "", "", fmt.Errorf("processing failed: %w", err)
	}
	if !success {
		return "", "", fmt.Errorf("media processing did not complete successfully")
	}

	// Step 3: Publish
	mediaID, err := s.PublishMedia(accessToken, igUserID, containerID)
	if err != nil {
		return "", "", fmt.Errorf("publish failed: %w", err)
	}

	// Step 4: Get permalink
	permalink, err := s.GetPermalink(accessToken, mediaID)
	if err != nil {
		// Don't fail if we can't get permalink, just return empty string
		permalink = ""
	}

	return mediaID, permalink, nil
}

// UploadAndPublishPhoto is a complete flow for uploading and publishing a photo
// Photos don't require async processing like videos do
func (s *InstagramMediaService) UploadAndPublishPhoto(
	accessToken string,
	igUserID string,
	imageURL string,
	caption string,
) (string, string, error) {
	// Step 1: Create photo container
	containerID, err := s.CreatePhotoContainer(accessToken, igUserID, imageURL, caption)
	if err != nil {
		return "", "", fmt.Errorf("create photo container failed: %w", err)
	}

	// Step 2: Wait briefly for container to be ready (photos process quickly)
	// Instagram recommends checking status even for photos
	success, err := s.WaitForMediaProcessing(accessToken, containerID, 60)
	if err != nil {
		return "", "", fmt.Errorf("photo processing failed: %w", err)
	}
	if !success {
		return "", "", fmt.Errorf("photo processing did not complete successfully")
	}

	// Step 3: Publish
	mediaID, err := s.PublishMedia(accessToken, igUserID, containerID)
	if err != nil {
		return "", "", fmt.Errorf("publish failed: %w", err)
	}

	// Step 4: Get permalink
	permalink, err := s.GetPermalink(accessToken, mediaID)
	if err != nil {
		// Don't fail if we can't get permalink, just return empty string
		permalink = ""
	}

	return mediaID, permalink, nil
}

// CreateCarouselItemContainer creates a container for a single item in a carousel
// This marks the media as a carousel item (is_carousel_item=true)
func (s *InstagramMediaService) CreateCarouselItemContainer(
	accessToken string,
	igUserID string,
	mediaURL string,
	isVideo bool,
) (string, error) {
	apiURL := fmt.Sprintf("https://graph.instagram.com/%s/media", igUserID)

	params := url.Values{}
	params.Set("is_carousel_item", "true")
	params.Set("access_token", accessToken)

	if isVideo {
		params.Set("media_type", "VIDEO")
		params.Set("video_url", mediaURL)
	} else {
		params.Set("image_url", mediaURL)
	}

	resp, err := s.httpClient.PostForm(apiURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to create carousel item container: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("create carousel item container failed (%d): %s", resp.StatusCode, string(body))
	}

	var containerResp CreateMediaContainerResponse
	if err := json.Unmarshal(body, &containerResp); err != nil {
		return "", fmt.Errorf("failed to parse container response: %w", err)
	}

	return containerResp.ID, nil
}

// CreateCarouselContainer creates a carousel container with multiple children
// This is the main container that references all the individual item containers
func (s *InstagramMediaService) CreateCarouselContainer(
	accessToken string,
	igUserID string,
	childrenIDs []string,
	caption string,
) (string, error) {
	if len(childrenIDs) < 2 {
		return "", fmt.Errorf("carousel requires at least 2 items, got %d", len(childrenIDs))
	}
	if len(childrenIDs) > 10 {
		return "", fmt.Errorf("carousel supports maximum 10 items, got %d", len(childrenIDs))
	}

	apiURL := fmt.Sprintf("https://graph.instagram.com/%s/media", igUserID)

	params := url.Values{}
	params.Set("media_type", "CAROUSEL")
	params.Set("access_token", accessToken)

	if caption != "" {
		params.Set("caption", caption)
	}

	// Children is a comma-separated list of container IDs
	childrenStr := ""
	for i, id := range childrenIDs {
		if i > 0 {
			childrenStr += ","
		}
		childrenStr += id
	}
	params.Set("children", childrenStr)

	resp, err := s.httpClient.PostForm(apiURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to create carousel container: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("create carousel container failed (%d): %s", resp.StatusCode, string(body))
	}

	var containerResp CreateMediaContainerResponse
	if err := json.Unmarshal(body, &containerResp); err != nil {
		return "", fmt.Errorf("failed to parse container response: %w", err)
	}

	return containerResp.ID, nil
}

// MediaItem represents a single media item in a carousel
type MediaItem struct {
	URL     string
	IsVideo bool
}

// UploadAndPublishCarousel is a complete flow for uploading and publishing a carousel
func (s *InstagramMediaService) UploadAndPublishCarousel(
	accessToken string,
	igUserID string,
	mediaItems []MediaItem,
	caption string,
) (string, string, error) {
	if len(mediaItems) < 2 {
		return "", "", fmt.Errorf("carousel requires at least 2 items")
	}
	if len(mediaItems) > 10 {
		return "", "", fmt.Errorf("carousel supports maximum 10 items")
	}

	// Step 1: Create individual containers for each media item
	var childrenIDs []string
	for i, item := range mediaItems {
		containerID, err := s.CreateCarouselItemContainer(accessToken, igUserID, item.URL, item.IsVideo)
		if err != nil {
			return "", "", fmt.Errorf("failed to create container for item %d: %w", i, err)
		}
		childrenIDs = append(childrenIDs, containerID)

		// Wait for each video item to process before creating carousel container
		if item.IsVideo {
			success, err := s.WaitForMediaProcessing(accessToken, containerID, 300)
			if err != nil {
				return "", "", fmt.Errorf("processing failed for video item %d: %w", i, err)
			}
			if !success {
				return "", "", fmt.Errorf("video item %d processing did not complete", i)
			}
		}
	}

	// Step 2: Create the carousel container with all children
	carouselID, err := s.CreateCarouselContainer(accessToken, igUserID, childrenIDs, caption)
	if err != nil {
		return "", "", fmt.Errorf("failed to create carousel container: %w", err)
	}

	// Step 3: Wait for carousel to be ready
	success, err := s.WaitForMediaProcessing(accessToken, carouselID, 300)
	if err != nil {
		return "", "", fmt.Errorf("carousel processing failed: %w", err)
	}
	if !success {
		return "", "", fmt.Errorf("carousel processing did not complete successfully")
	}

	// Step 4: Publish
	mediaID, err := s.PublishMedia(accessToken, igUserID, carouselID)
	if err != nil {
		return "", "", fmt.Errorf("publish failed: %w", err)
	}

	// Step 5: Get permalink
	permalink, err := s.GetPermalink(accessToken, mediaID)
	if err != nil {
		// Don't fail if we can't get permalink, just return empty string
		permalink = ""
	}

	return mediaID, permalink, nil
}
