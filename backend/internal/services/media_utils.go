package services

import (
	"net/url"
	"path"
	"strings"
)

// MediaType represents the type of media (video or image)
type MediaType string

const (
	MediaTypeVideo   MediaType = "video"
	MediaTypeImage   MediaType = "image"
	MediaTypeUnknown MediaType = "unknown"
)

// imageExtensions contains common image file extensions
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".tiff": true,
	".heic": true,
	".heif": true,
}

// videoExtensions contains common video file extensions
var videoExtensions = map[string]bool{
	".mp4":  true,
	".mov":  true,
	".webm": true,
	".avi":  true,
	".mkv":  true,
	".m4v":  true,
	".ts":   true,
	".3gp":  true,
}

// DetectMediaTypeFromURL determines if the URL points to a video or image
// based on the file extension in the URL path
func DetectMediaTypeFromURL(mediaURL string) MediaType {
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		return MediaTypeUnknown
	}

	// Get the path and extract extension
	urlPath := parsedURL.Path
	ext := strings.ToLower(path.Ext(urlPath))

	// Check if it's an image
	if imageExtensions[ext] {
		return MediaTypeImage
	}

	// Check if it's a video
	if videoExtensions[ext] {
		return MediaTypeVideo
	}

	// Default to video for backwards compatibility
	return MediaTypeVideo
}

// IsImageURL checks if the URL points to an image
func IsImageURL(mediaURL string) bool {
	return DetectMediaTypeFromURL(mediaURL) == MediaTypeImage
}

// IsVideoURL checks if the URL points to a video
func IsVideoURL(mediaURL string) bool {
	return DetectMediaTypeFromURL(mediaURL) == MediaTypeVideo
}
