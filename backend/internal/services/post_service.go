package services

import (
	"fmt"
	"log"
	"time"

	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
)

type PostService struct {
	postRepo      *models.PostRepository
	tokenService  *TokenService
	tiktokService *TikTokService
}

// NewPostService creates a new post service
func NewPostService(postRepo *models.PostRepository, tokenService *TokenService, tiktokService *TikTokService) *PostService {
	return &PostService{
		postRepo:      postRepo,
		tokenService:  tokenService,
		tiktokService: tiktokService,
	}
}

// CreatePost creates a new post and publishes it to TikTok
func (s *PostService) CreatePost(userID int64, videoURL string, caption string) (*models.Post, error) {
	// Create post record in database
	post := &models.Post{
		UserID:   userID,
		VideoURL: videoURL,
		Caption:  caption,
		Status:   models.PostStatusPending,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Process post asynchronously
	go s.processPost(post.ID, userID)

	return post, nil
}

// processPost handles the posting workflow asynchronously
func (s *PostService) processPost(postID int64, userID int64) {
	// Update status to processing
	if err := s.postRepo.UpdateStatus(postID, models.PostStatusProcessing, ""); err != nil {
		log.Printf("Failed to update post %d status to processing: %v", postID, err)
		return
	}

	// Get post details
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		log.Printf("Failed to get post %d: %v", postID, err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, "Failed to retrieve post details")
		return
	}

	// Get valid access token
	accessToken, err := s.tokenService.GetValidToken(userID)
	if err != nil {
		log.Printf("Failed to get access token for user %d: %v", userID, err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, "Failed to get access token")
		return
	}

	// Publish video to TikTok
	publishResponse, err := s.tiktokService.PublishVideoFromURL(accessToken, post.VideoURL, post.Caption)
	if err != nil {
		log.Printf("Failed to publish video to TikTok: %v", err)
		s.postRepo.UpdateStatus(postID, models.PostStatusFailed, fmt.Sprintf("TikTok error: %v", err))
		return
	}

	log.Printf("Video published to TikTok with publish_id: %s", publishResponse.Data.PublishID)

	// Poll status until complete
	publishID := publishResponse.Data.PublishID
	maxAttempts := 60 // Poll for up to 5 minutes (60 * 5 seconds)
	attempt := 0

	for attempt < maxAttempts {
		time.Sleep(5 * time.Second)
		attempt++

		statusResponse, err := s.tiktokService.GetPublishStatus(accessToken, publishID)
		if err != nil {
			log.Printf("Failed to get publish status: %v", err)
			continue
		}

		log.Printf("Publish status for %s: %s", publishID, statusResponse.Data.Status)

		switch statusResponse.Data.Status {
		case "PUBLISH_COMPLETE":
			// Video successfully published
			shareID := statusResponse.Data.ShareID
			if err := s.postRepo.MarkPublished(postID, shareID); err != nil {
				log.Printf("Failed to mark post %d as published: %v", postID, err)
			} else {
				log.Printf("Post %d successfully published with share_id: %s", postID, shareID)
			}
			return

		case "FAILED":
			// Publishing failed
			failReason := statusResponse.Data.FailReason
			if failReason == "" {
				failReason = "Unknown error"
			}
			s.postRepo.UpdateStatus(postID, models.PostStatusFailed, fmt.Sprintf("TikTok publish failed: %s", failReason))
			log.Printf("Post %d failed: %s", postID, failReason)
			return

		case "PROCESSING_DOWNLOAD":
			log.Printf("Post %d: TikTok is downloading video from URL", postID)

		case "PROCESSING_UPLOAD":
			log.Printf("Post %d: TikTok is processing the video", postID)

		default:
			log.Printf("Post %d: Unknown status %s", postID, statusResponse.Data.Status)
		}
	}

	// Timeout reached
	s.postRepo.UpdateStatus(postID, models.PostStatusFailed, "Publishing timeout - took too long")
	log.Printf("Post %d timed out after %d attempts", postID, maxAttempts)
}

// GetPostByID retrieves a post by ID
func (s *PostService) GetPostByID(postID int64, userID int64) (*models.Post, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// Verify post belongs to user
	if post.UserID != userID {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

// GetUserPosts retrieves all posts for a user
func (s *PostService) GetUserPosts(userID int64, limit, offset int) ([]*models.Post, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.postRepo.GetByUserID(userID, limit, offset)
}

// GetPostStatus retrieves the current status of a post
func (s *PostService) GetPostStatus(postID int64, userID int64) (*models.Post, error) {
	return s.GetPostByID(postID, userID)
}
