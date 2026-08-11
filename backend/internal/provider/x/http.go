package x

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiBase = "https://api.x.com"

type apiErrorResponse struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (e apiErrorResponse) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("x: %s", e.Detail)
	}
	if len(e.Errors) > 0 {
		return fmt.Sprintf("x: %s", e.Errors[0].Message)
	}
	return "x: request failed"
}

func bearerJSON(ctx context.Context, method, accessToken, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return err
		}
		reqBody = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, apiBase+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		var apiErr apiErrorResponse
		_ = json.Unmarshal(respBody, &apiErr)
		if apiErr.Detail == "" && len(apiErr.Errors) == 0 {
			return fmt.Errorf("x: unexpected response (status %d): %s", resp.StatusCode, respBody)
		}
		return apiErr
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(respBody, out)
}
