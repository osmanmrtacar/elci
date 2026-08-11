// Package instagram implements provider.Provider against the newer
// "Instagram API with Instagram Login" — no linked Facebook Page required,
// authenticating and publishing entirely through graph.instagram.com.
package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const apiVersion = "v26.0"
const graphBase = "https://graph.instagram.com/" + apiVersion

type apiErrorResponse struct {
	ErrorDetail struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (e apiErrorResponse) Error() string {
	return fmt.Sprintf("instagram: %s (type=%s, code=%d)", e.ErrorDetail.Message, e.ErrorDetail.Type, e.ErrorDetail.Code)
}

func getGraph(ctx context.Context, path string, query url.Values, out any) error {
	return doGraph(ctx, http.MethodGet, graphBase+path, query, out)
}

func postGraph(ctx context.Context, path string, form url.Values, out any) error {
	return doGraph(ctx, http.MethodPost, graphBase+path, form, out)
}

func doGraph(ctx context.Context, method, fullURL string, params url.Values, out any) error {
	reqURL := fullURL
	if method == http.MethodGet {
		reqURL += "?" + params.Encode()
	}

	var body io.Reader
	if method != http.MethodGet {
		body = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return err
	}
	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.ErrorDetail.Message != "" {
			return apiErr
		}
		return fmt.Errorf("instagram: unexpected response (status %d): %s", resp.StatusCode, respBody)
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(respBody, out)
}
