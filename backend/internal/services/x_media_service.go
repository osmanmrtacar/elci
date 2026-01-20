package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	xMediaUploadURL = "https://api.x.com/2/media/upload"
	xChunkSize      = 512 * 1024 // 512KB chunks (X API limit)
)

// XMediaService handles media uploads to X (Twitter)
type XMediaService struct {
	httpClient *http.Client
}

// NewXMediaService creates a new X media service
func NewXMediaService() *XMediaService {
	return &XMediaService{
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for uploads
		},
	}
}

// InitResponse represents the response from initializing an upload
type InitResponse struct {
	Data struct {
		ID               string `json:"id"`
		MediaKey         string `json:"media_key"`
		ExpiresAfterSecs int    `json:"expires_after_secs"`
	} `json:"data"`
}

// FinalizeResponse represents the response from finalizing an upload
type FinalizeResponse struct {
	Data struct {
		ID             string          `json:"id"`
		MediaKey       string          `json:"media_key"`
		Size           int64           `json:"size"`
		ProcessingInfo *ProcessingInfo `json:"processing_info,omitempty"`
	} `json:"data"`
}

// ProcessingInfo represents video processing status
type ProcessingInfo struct {
	State           string `json:"state"` // pending, in_progress, succeeded, failed
	CheckAfterSecs  int    `json:"check_after_secs,omitempty"`
	ProgressPercent int    `json:"progress_percent,omitempty"`
}

// StatusResponse represents the response from checking upload status
type StatusResponse struct {
	Data struct {
		ID             string          `json:"id"`
		ProcessingInfo *ProcessingInfo `json:"processing_info,omitempty"`
	} `json:"data"`
}

// GetMediaType determines media type from file extension
func (s *XMediaService) GetMediaType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	typeMap := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".tiff": "image/tiff",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".webm": "video/webm",
		".ts":   "video/mp2t",
	}
	if mediaType, ok := typeMap[ext]; ok {
		return mediaType
	}
	return "video/mp4" // Default
}

// GetMediaCategory determines X media category from media type
func (s *XMediaService) GetMediaCategory(mediaType string) string {
	if strings.HasPrefix(mediaType, "video/") {
		return "amplify_video"
	}
	if mediaType == "image/gif" {
		return "tweet_gif"
	}
	return "tweet_image"
}

// InitUpload initializes a chunked upload
func (s *XMediaService) InitUpload(accessToken, filePath string) (*InitResponse, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	mediaType := s.GetMediaType(filePath)
	mediaCategory := s.GetMediaCategory(mediaType)

	requestBody := map[string]interface{}{
		"media_type":     mediaType,
		"total_bytes":    fileInfo.Size(),
		"media_category": mediaCategory,
	}

	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", xMediaUploadURL+"/initialize", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Initializing upload: %s (%d bytes, %s, %s)\n",
		filepath.Base(filePath), fileInfo.Size(), mediaType, mediaCategory)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("X API init error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var initResp InitResponse
	if err := json.Unmarshal(respBody, &initResp); err != nil {
		return nil, fmt.Errorf("failed to parse init response: %w", err)
	}

	fmt.Printf("Upload initialized: media_id=%s\n", initResp.Data.ID)
	return &initResp, nil
}

// AppendChunks uploads file in chunks
// CRITICAL: Read entire file into memory ONCE, then slice for each chunk
// This is what was broken in the TypeScript implementation before the fix
func (s *XMediaService) AppendChunks(accessToken, mediaID, filePath string) error {
	// CRITICAL: Read entire file into buffer ONCE (matching working TypeScript implementation)
	fileBuffer, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fileSize := len(fileBuffer)
	totalChunks := (fileSize + xChunkSize - 1) / xChunkSize

	fmt.Printf("File loaded: %d bytes, uploading %d chunks\n", fileSize, totalChunks)

	for segmentIndex := 0; segmentIndex < totalChunks; segmentIndex++ {
		start := segmentIndex * xChunkSize
		end := start + xChunkSize
		if end > fileSize {
			end = fileSize
		}

		// CRITICAL: Slice the buffer to get chunk (don't re-read file!)
		chunk := fileBuffer[start:end]

		fmt.Printf("Chunk %d: %d bytes (%d-%d)\n", segmentIndex, len(chunk), start, end)

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add media chunk
		part, err := writer.CreateFormFile("media", "blob")
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(chunk); err != nil {
			return fmt.Errorf("failed to write chunk: %w", err)
		}

		// Add segment index
		if err := writer.WriteField("segment_index", fmt.Sprintf("%d", segmentIndex)); err != nil {
			return fmt.Errorf("failed to write segment index: %w", err)
		}

		if err := writer.Close(); err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		// Upload chunk
		req, err := http.NewRequest("POST",
			fmt.Sprintf("%s/%s/append", xMediaUploadURL, mediaID), body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to upload chunk %d: %w", segmentIndex, err)
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("chunk %d upload failed (status %d): %s",
				segmentIndex, resp.StatusCode, string(respBody))
		}

		fmt.Printf("Uploaded chunk %d/%d (%.0f%%)\n",
			segmentIndex+1, totalChunks, float64(end)/float64(fileSize)*100)
	}

	return nil
}

// FinalizeUpload finalizes the upload
func (s *XMediaService) FinalizeUpload(accessToken, mediaID string) (*FinalizeResponse, error) {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/%s/finalize", xMediaUploadURL, mediaID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("finalize error (status %d): %s", resp.StatusCode, string(body))
	}

	var finalizeResp FinalizeResponse
	if err := json.Unmarshal(body, &finalizeResp); err != nil {
		return nil, fmt.Errorf("failed to parse finalize response: %w", err)
	}

	fmt.Printf("Upload finalized: media_id=%s, size=%d\n",
		finalizeResp.Data.ID, finalizeResp.Data.Size)

	return &finalizeResp, nil
}

// CheckStatus checks the processing status of uploaded media
func (s *XMediaService) CheckStatus(accessToken, mediaID string) (*StatusResponse, error) {
	req, err := http.NewRequest("GET", xMediaUploadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("command", "STATUS")
	q.Add("media_id", mediaID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status check error (status %d): %s", resp.StatusCode, string(body))
	}

	var statusResp StatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	return &statusResp, nil
}

// WaitForProcessing polls until video processing completes
func (s *XMediaService) WaitForProcessing(accessToken, mediaID string, maxWaitSec int) error {
	startTime := time.Now()

	for {
		statusResp, err := s.CheckStatus(accessToken, mediaID)
		if err != nil {
			return err
		}

		if statusResp.Data.ProcessingInfo == nil {
			// No processing needed
			fmt.Println("No processing required")
			return nil
		}

		state := statusResp.Data.ProcessingInfo.State

		if state == "succeeded" {
			fmt.Println("Media processing completed")
			return nil
		}

		if state == "failed" {
			return fmt.Errorf("media processing failed")
		}

		// Check timeout
		elapsed := time.Since(startTime).Seconds()
		if elapsed > float64(maxWaitSec) {
			return fmt.Errorf("media processing timeout after %d seconds", maxWaitSec)
		}

		// Wait before next check
		waitSec := statusResp.Data.ProcessingInfo.CheckAfterSecs
		if waitSec == 0 {
			waitSec = 5
		}

		fmt.Printf("Processing %s (%d%%), checking again in %ds...\n",
			state, statusResp.Data.ProcessingInfo.ProgressPercent, waitSec)

		time.Sleep(time.Duration(waitSec) * time.Second)
	}
}

// UploadFromURL downloads media from URL and uploads it to X
// This is the complete upload flow: download → init → append → finalize → wait
func (s *XMediaService) UploadFromURL(accessToken, mediaURL string) (string, error) {
	// Create temp file
	tempFile := fmt.Sprintf("/tmp/x-media-%d.tmp", time.Now().UnixNano())
	defer os.Remove(tempFile)

	// Download file
	fmt.Printf("Downloading media from: %s\n", mediaURL)
	resp, err := http.Get(mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download media: status %d", resp.StatusCode)
	}

	out, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save media: %w", err)
	}

	fmt.Println("Media downloaded successfully")

	// Initialize upload
	initResp, err := s.InitUpload(accessToken, tempFile)
	if err != nil {
		return "", err
	}

	mediaID := initResp.Data.ID

	// Upload chunks
	if err := s.AppendChunks(accessToken, mediaID, tempFile); err != nil {
		return "", err
	}

	// Finalize upload
	finalizeResp, err := s.FinalizeUpload(accessToken, mediaID)
	if err != nil {
		return "", err
	}

	// Wait for processing if needed
	if finalizeResp.Data.ProcessingInfo != nil {
		if err := s.WaitForProcessing(accessToken, mediaID, 300); err != nil {
			return "", err
		}
	}

	fmt.Printf("Media uploaded successfully: %s\n", mediaID)
	return mediaID, nil
}

// UploadMultipleFromURLs downloads and uploads multiple media files to X
// X allows maximum 4 photos OR 1 video per tweet
func (s *XMediaService) UploadMultipleFromURLs(accessToken string, mediaURLs []string) ([]string, error) {
	if len(mediaURLs) == 0 {
		return nil, fmt.Errorf("at least one media URL is required")
	}

	// Check if any URL is a video - videos can't be mixed with images
	hasVideo := false
	for _, url := range mediaURLs {
		lowerURL := strings.ToLower(url)
		if strings.HasSuffix(lowerURL, ".mp4") || strings.HasSuffix(lowerURL, ".mov") ||
			strings.HasSuffix(lowerURL, ".webm") || strings.HasSuffix(lowerURL, ".ts") {
			hasVideo = true
			break
		}
	}

	// Validate limits
	if hasVideo && len(mediaURLs) > 1 {
		return nil, fmt.Errorf("X only allows 1 video per tweet, got %d media items", len(mediaURLs))
	}
	if !hasVideo && len(mediaURLs) > 4 {
		return nil, fmt.Errorf("X only allows maximum 4 photos per tweet, got %d", len(mediaURLs))
	}

	// Upload each media item
	var mediaIDs []string
	for i, mediaURL := range mediaURLs {
		fmt.Printf("Uploading media %d/%d: %s\n", i+1, len(mediaURLs), mediaURL)
		mediaID, err := s.UploadFromURL(accessToken, mediaURL)
		if err != nil {
			return nil, fmt.Errorf("failed to upload media %d: %w", i+1, err)
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	return mediaIDs, nil
}
