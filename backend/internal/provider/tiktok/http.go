package tiktok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const apiBase = "https://open.tiktokapis.com/v2"

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LogID   string `json:"log_id"`
}

func (e apiError) isOK() bool {
	return e.Code == "" || e.Code == "ok"
}

func (e apiError) Error() string {
	return fmt.Sprintf("tiktok: %s (%s, log_id=%s)", e.Message, e.Code, e.LogID)
}

// postJSON sends a bearer-authenticated JSON POST and decodes the "data"
// envelope into out, returning the API's own error if it reports one.
func postJSON(ctx context.Context, accessToken, path string, body, out any) error {
	buf := &bytes.Buffer{}
	if body == nil {
		buf.WriteString("{}")
	} else if err := json.NewEncoder(buf).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+path, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	return doAndDecode(req, out)
}

func getJSON(ctx context.Context, accessToken, path string, query url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+path+"?"+query.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return doAndDecode(req, out)
}

// postOAuthForm calls the OAuth token endpoint, which — unlike the rest of
// the v2 API — returns fields at the top level and errors as
// {"error": "...", "error_description": "..."} rather than the usual
// {"data": ..., "error": {...}} envelope.
func postOAuthForm(ctx context.Context, fullURL string, form url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var maybeErr struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := json.Unmarshal(respBody, &maybeErr); err == nil && maybeErr.Error != "" {
		return fmt.Errorf("tiktok oauth: %s: %s", maybeErr.Error, maybeErr.Description)
	}

	return json.Unmarshal(respBody, out)
}

func doAndDecode(req *http.Request, out any) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var envelope struct {
		Data  json.RawMessage `json:"data"`
		Error apiError        `json:"error"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return fmt.Errorf("tiktok: unexpected response (status %d): %s", resp.StatusCode, respBody)
	}
	if !envelope.Error.isOK() {
		return envelope.Error
	}
	if out == nil || len(envelope.Data) == 0 {
		return nil
	}
	return json.Unmarshal(envelope.Data, out)
}
