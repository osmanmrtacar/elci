package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/osmanmertacar/sosyal/backend/internal/api/middleware"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
)

type PostHandler struct {
	postService *services.PostService
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePostRequest represents the request to create a post
type CreatePostRequest struct {
	VideoURL string `json:"video_url" binding:"required"`
	Caption  string `json:"caption"`
}

// CreatePost creates a new post
func (h *PostHandler) CreatePost(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse request body
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate video URL
	if req.VideoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video_url is required"})
		return
	}

	// Create post
	post, err := h.postService.CreatePost(userID, req.VideoURL, req.Caption)
	if err != nil {
		log.Printf("Failed to create post for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	log.Printf("Post %d created for user %d", post.ID, userID)

	// Return post
	c.JSON(http.StatusCreated, gin.H{
		"post": gin.H{
			"id":         post.ID,
			"video_url":  post.VideoURL,
			"caption":    post.Caption,
			"status":     post.Status,
			"created_at": post.CreatedAt,
		},
		"message": "Post created and is being processed",
	})
}

// GetPosts retrieves all posts for the authenticated user
func (h *PostHandler) GetPosts(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Get pagination parameters
	limit := 20
	offset := 0

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

	// Get posts
	posts, err := h.postService.GetUserPosts(userID, limit, offset)
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
			"video_url":  post.VideoURL,
			"caption":    post.Caption,
			"status":     post.Status,
			"created_at": post.CreatedAt,
		}

		if post.TikTokPostID != "" {
			postData["tiktok_post_id"] = post.TikTokPostID
			postData["tiktok_url"] = "https://www.tiktok.com/@user/video/" + post.TikTokPostID
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
func (h *PostHandler) GetPost(c *gin.Context) {
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
		"video_url":  post.VideoURL,
		"caption":    post.Caption,
		"status":     post.Status,
		"created_at": post.CreatedAt,
	}

	if post.TikTokPostID != "" {
		postData["tiktok_post_id"] = post.TikTokPostID
		postData["tiktok_url"] = "https://www.tiktok.com/@user/video/" + post.TikTokPostID
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
func (h *PostHandler) GetPostStatus(c *gin.Context) {
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

	// Return just the status information
	response := gin.H{
		"id":     post.ID,
		"status": post.Status,
	}

	if post.ErrorMessage != "" {
		response["error_message"] = post.ErrorMessage
	}

	if post.TikTokPostID != "" {
		response["tiktok_post_id"] = post.TikTokPostID
		response["tiktok_url"] = "https://www.tiktok.com/@user/video/" + post.TikTokPostID
	}

	if post.PublishedAt != nil {
		response["published_at"] = post.PublishedAt
	}

	c.JSON(http.StatusOK, response)
}
