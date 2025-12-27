package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxVideoSize        = 300 * 1024 * 1024 // 300 MB in bytes (TikTok limit is 287.6 MB)
	minVideoDuration    = 3                  // seconds
	maxVideoDuration    = 600                // 10 minutes in seconds
	tiktokChunkSize     = 5 * 1024 * 1024    // 5 MB chunks for upload
)

type VideoService struct {
	httpClient *http.Client
	tempDir    string
}

// VideoInfo contains information about a downloaded video
type VideoInfo struct {
	Path      string
	Size      int64
	MimeType  string
	Extension string
}

// NewVideoService creates a new video service
func NewVideoService() *VideoService {
	return &VideoService{
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Allow time for large video downloads
		},
		tempDir: os.TempDir(),
	}
}

// DownloadVideo downloads a video from a URL to a temporary file
func (s *VideoService) DownloadVideo(url string) (*VideoInfo, error) {
	// Validate URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("invalid URL: must start with http:// or https://")
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download video: HTTP %d", resp.StatusCode)
	}

	// Get content type
	contentType := resp.Header.Get("Content-Type")

	// Determine file extension from content type
	extension := s.getExtensionFromContentType(contentType)
	if extension == "" {
		// Try to get extension from URL
		extension = s.getExtensionFromURL(url)
	}
	if extension == "" {
		extension = ".mp4" // Default to mp4
	}

	// Create temporary file
	tempFile, err := os.CreateTemp(s.tempDir, "tiktok_video_*"+extension)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Download video with size limit
	limitedReader := io.LimitReader(resp.Body, maxVideoSize+1)
	written, err := io.Copy(tempFile, limitedReader)
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("failed to write video to file: %w", err)
	}

	// Check if video exceeds size limit
	if written > maxVideoSize {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("video file too large: %.2f MB (max: 287.6 MB)", float64(written)/(1024*1024))
	}

	videoInfo := &VideoInfo{
		Path:      tempFile.Name(),
		Size:      written,
		MimeType:  contentType,
		Extension: extension,
	}

	return videoInfo, nil
}

// ValidateVideo validates video format and size
func (s *VideoService) ValidateVideo(videoInfo *VideoInfo) error {
	// Check file size
	if videoInfo.Size == 0 {
		return fmt.Errorf("video file is empty")
	}

	if videoInfo.Size > maxVideoSize {
		return fmt.Errorf("video file too large: %.2f MB (max: 287.6 MB)", float64(videoInfo.Size)/(1024*1024))
	}

	// Validate file extension
	ext := strings.ToLower(videoInfo.Extension)
	validExtensions := []string{".mp4", ".mov", ".webm"}
	valid := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("unsupported video format: %s (supported: MP4, MOV, WEBM)", ext)
	}

	// Note: Duration validation would require ffmpeg or similar
	// For now, we'll rely on TikTok's API to reject invalid durations

	return nil
}

// ReadVideoFile reads the entire video file into memory
func (s *VideoService) ReadVideoFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read video file: %w", err)
	}
	return data, nil
}

// SplitIntoChunks splits video data into chunks for upload
func (s *VideoService) SplitIntoChunks(data []byte) [][]byte {
	var chunks [][]byte
	dataLen := len(data)

	for i := 0; i < dataLen; i += tiktokChunkSize {
		end := i + tiktokChunkSize
		if end > dataLen {
			end = dataLen
		}
		chunks = append(chunks, data[i:end])
	}

	return chunks
}

// CleanupVideo removes a temporary video file
func (s *VideoService) CleanupVideo(path string) error {
	if path == "" {
		return nil
	}

	// Only delete files in temp directory for safety
	if !strings.HasPrefix(path, s.tempDir) {
		return fmt.Errorf("refusing to delete file outside temp directory: %s", path)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove video file: %w", err)
	}
	return nil
}

// getExtensionFromContentType returns file extension based on MIME type
func (s *VideoService) getExtensionFromContentType(contentType string) string {
	contentType = strings.ToLower(contentType)
	switch {
	case strings.Contains(contentType, "mp4"):
		return ".mp4"
	case strings.Contains(contentType, "quicktime"):
		return ".mov"
	case strings.Contains(contentType, "webm"):
		return ".webm"
	default:
		return ""
	}
}

// getExtensionFromURL extracts file extension from URL
func (s *VideoService) getExtensionFromURL(url string) string {
	ext := filepath.Ext(url)
	// Remove query parameters if present
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}
	return strings.ToLower(ext)
}
