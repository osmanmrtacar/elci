package services

import (
	"fmt"
	"log"
)

// InstagramPostService handles Instagram post creation
type InstagramPostService struct {
	authService  *InstagramAuthService
	mediaService *InstagramMediaService
}

// NewInstagramPostService creates a new Instagram post service
func NewInstagramPostService(authService *InstagramAuthService, mediaService *InstagramMediaService) *InstagramPostService {
	return &InstagramPostService{
		authService:  authService,
		mediaService: mediaService,
	}
}

// CreatePost creates and publishes a post to Instagram
func (s *InstagramPostService) CreatePost(accessToken string, videoURL string, caption string) (string, string, error) {
	// Get Instagram user info (need IG user ID for posting)
	userInfo, err := s.authService.GetInstagramUserInfo(accessToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to get Instagram user info: %w", err)
	}

	log.Printf("Posting to Instagram account: @%s (ID: %s)", userInfo.Username, userInfo.ID)

	// Upload and publish the reel
	mediaID, permalink, err := s.mediaService.UploadAndPublishReel(
		accessToken,
		userInfo.ID,
		videoURL,
		caption,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to upload and publish: %w", err)
	}

	return mediaID, permalink, nil
}
