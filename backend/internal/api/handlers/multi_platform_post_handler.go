package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/osmanmertacar/sosyal/backend/internal/api/middleware"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
)

type MultiPlatformPostHandler struct {
	postService *services.MultiPlatformPostService
}

// NewMultiPlatformPostHandler creates a new multi-platform post handler
func NewMultiPlatformPostHandler(postService *services.MultiPlatformPostService) *MultiPlatformPostHandler {
	return &MultiPlatformPostHandler{
		postService: postService,
	}
}

// CreatePostRequest represents the request to create a post on multiple platforms
type CreateMultiPlatformPostRequest struct {
	Platforms []string `json:"platforms" binding:"required"` // ["tiktok", "x"]
	MediaURL  string   `json:"media_url" binding:"required"` // Video/image URL
	Caption   string   `json:"caption"`                      // Post text/caption
}

// CreatePost creates a new post on one or more platforms
func (h *MultiPlatformPostHandler) CreatePost(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse request body
	var req CreateMultiPlatformPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate request
	if len(req.Platforms) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one platform must be specified"})
		return
	}

	if req.MediaURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "media_url is required"})
		return
	}

	// Convert platform strings to Platform type
	platforms := make([]models.Platform, 0, len(req.Platforms))
	for _, p := range req.Platforms {
		platforms = append(platforms, models.Platform(p))
	}

	// Create post service request
	serviceReq := services.CreateMultiPlatformPostRequest{
		Platforms: platforms,
		MediaURL:  req.MediaURL,
		Caption:   req.Caption,
	}

	// Create posts
	resp, err := h.postService.CreateMultiPlatformPost(userID, serviceReq)
	if err != nil {
		log.Printf("Failed to create posts for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format response
	postList := make([]gin.H, 0, len(resp.Posts))
	for _, post := range resp.Posts {
		postData := gin.H{
			"id":         post.ID,
			"platform":   post.Platform,
			"media_url":  post.VideoURL,
			"caption":    post.Caption,
			"status":     post.Status,
			"created_at": post.CreatedAt,
		}
		postList = append(postList, postData)
	}

	log.Printf("Created %d posts for user %d across %d platforms", len(resp.Posts), userID, len(req.Platforms))

	response := gin.H{
		"posts":   postList,
		"message": "Posts created and are being processed",
	}

	if len(resp.Errors) > 0 {
		response["errors"] = resp.Errors
	}

	c.JSON(http.StatusCreated, response)
}

// GetPosts retrieves all posts for the authenticated user
func (h *MultiPlatformPostHandler) GetPosts(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Get query parameters
	limit := 20
	offset := 0
	platformFilter := models.Platform("")

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	if platform := c.Query("platform"); platform != "" {
		platformFilter = models.Platform(platform)
	}

	// Get posts
	posts, err := h.postService.GetUserPosts(userID, platformFilter, limit, offset)
	if err != nil {
		log.Printf("Failed to get posts for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	// Format response
	postList := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		postData := gin.H{
			"id":         post.ID,
			"platform":   post.Platform,
			"media_url":  post.VideoURL,
			"caption":    post.Caption,
			"status":     post.Status,
			"media_type": post.MediaType,
			"created_at": post.CreatedAt,
		}

		// Add platform-specific post ID and URL
		if post.PlatformPostID != "" {
			postData["platform_post_id"] = post.PlatformPostID

			// Generate platform-specific URLs
			switch post.Platform {
			case models.PlatformTikTok:
				postData["share_url"] = "https://www.tiktok.com/@user/video/" + post.PlatformPostID
			case models.PlatformX:
				postData["share_url"] = "https://twitter.com/i/web/status/" + post.PlatformPostID
			case models.PlatformInstagram:
				// Instagram post ID is the actual media ID, no URL construction needed
				// The ShareURL from the API response should be stored separately if needed
				postData["share_url"] = "https://www.instagram.com/p/" + post.PlatformPostID
			}
		}

		// Legacy TikTok field for backward compatibility
		if post.TikTokPostID != "" && post.Platform == models.PlatformTikTok {
			postData["tiktok_post_id"] = post.TikTokPostID
		}

		if post.PublishedAt != nil {
			postData["published_at"] = post.PublishedAt
		}

		if post.ErrorMessage != "" {
			postData["error_message"] = post.ErrorMessage
		}

		postList = append(postList, postData)
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  postList,
		"count":  len(postList),
		"limit":  limit,
		"offset": offset,
	})
}

// GetPost retrieves a specific post by ID
func (h *MultiPlatformPostHandler) GetPost(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Get post ID from URL
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Get post
	post, err := h.postService.GetPostByID(postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Format response
	postData := gin.H{
		"id":         post.ID,
		"platform":   post.Platform,
		"media_url":  post.VideoURL,
		"caption":    post.Caption,
		"status":     post.Status,
		"media_type": post.MediaType,
		"created_at": post.CreatedAt,
	}

	if post.PlatformPostID != "" {
		postData["platform_post_id"] = post.PlatformPostID

		// Generate platform-specific URLs
		switch post.Platform {
		case models.PlatformTikTok:
			postData["share_url"] = "https://www.tiktok.com/@user/video/" + post.PlatformPostID
		case models.PlatformX:
			postData["share_url"] = "https://twitter.com/i/web/status/" + post.PlatformPostID
		case models.PlatformInstagram:
			postData["share_url"] = "https://www.instagram.com/p/" + post.PlatformPostID
		}
	}

	if post.PublishedAt != nil {
		postData["published_at"] = post.PublishedAt
	}

	if post.ErrorMessage != "" {
		postData["error_message"] = post.ErrorMessage
	}

	c.JSON(http.StatusOK, gin.H{
		"post": postData,
	})
}

// GetPostStatus retrieves the status of a specific post
func (h *MultiPlatformPostHandler) GetPostStatus(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Get post ID from URL
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Get post status
	post, err := h.postService.GetPostStatus(postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Return status information
	response := gin.H{
		"id":       post.ID,
		"platform": post.Platform,
		"status":   post.Status,
	}

	if post.ErrorMessage != "" {
		response["error_message"] = post.ErrorMessage
	}

	if post.PlatformPostID != "" {
		response["platform_post_id"] = post.PlatformPostID

		// Generate platform-specific URLs
		switch post.Platform {
		case models.PlatformTikTok:
			response["share_url"] = "https://www.tiktok.com/@user/video/" + post.PlatformPostID
		case models.PlatformX:
			response["share_url"] = "https://twitter.com/i/web/status/" + post.PlatformPostID
		case models.PlatformInstagram:
			response["share_url"] = "https://www.instagram.com/p/" + post.PlatformPostID
		}
	}

	if post.PublishedAt != nil {
		response["published_at"] = post.PublishedAt
	}

	c.JSON(http.StatusOK, response)
}
