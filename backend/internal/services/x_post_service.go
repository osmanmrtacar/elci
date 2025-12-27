package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// XPostService handles creating posts (tweets) on X (Twitter)
type XPostService struct {
	httpClient *http.Client
}

// NewXPostService creates a new X post service
func NewXPostService() *XPostService {
	return &XPostService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// XPostRequest represents a request to create a tweet
type XPostRequest struct {
	Text     string   `json:"text"`
	MediaIDs []string `json:"media_ids,omitempty"`
}

// XPostResponse represents the response from creating a tweet
type XPostResponse struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

// CreatePost creates a tweet with optional media
func (s *XPostService) CreatePost(accessToken, text string, mediaIDs []string) (*XPostResponse, error) {
	requestBody := map[string]interface{}{
		"text": text,
	}

	// Add media if provided
	if len(mediaIDs) > 0 {
		requestBody["media"] = map[string]interface{}{
			"media_ids": mediaIDs,
		}
	}

	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Creating tweet: %s (media: %d files)\n", text[:min(50, len(text))], len(mediaIDs))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("X API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var postResp XPostResponse
	if err := json.Unmarshal(respBody, &postResp); err != nil {
		return nil, fmt.Errorf("failed to parse post response: %w", err)
	}

	fmt.Printf("Tweet created: %s\n", postResp.Data.ID)
	return &postResp, nil
}

// GetUserTweets retrieves tweets for a user
func (s *XPostService) GetUserTweets(accessToken, userID string, maxResults int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets?max_results=%d&tweet.fields=created_at",
		userID, maxResults)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("X API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
