package handlers

import (
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/osmanmertacar/elci/backend/internal/api/middleware"
	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

type MediaHandler struct {
	media *service.MediaService
	repo  repository.MediaRepository
}

func NewMediaHandler(media *service.MediaService, repo repository.MediaRepository) *MediaHandler {
	return &MediaHandler{media: media, repo: repo}
}

type presignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

func (h *MediaHandler) Presign(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req presignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Filename == "" || req.ContentType == "" {
		writeError(w, http.StatusBadRequest, "filename and contentType are required")
		return
	}

	upload, err := h.media.PresignUpload(r.Context(), userID, req.Filename, req.ContentType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create upload URL")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"key":       upload.Key,
		"uploadUrl": upload.UploadURL,
		"publicUrl": upload.PublicURL,
	})
}

type confirmRequest struct {
	Key             string `json:"key"`
	PublicURL       string `json:"publicUrl"`
	Kind            string `json:"kind"`
	ContentType     string `json:"contentType"`
	SizeBytes       int64  `json:"sizeBytes"`
	DurationSeconds *int   `json:"durationSeconds"`
}

func (h *MediaHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req confirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	kind := domain.MediaKind(req.Kind)
	if kind != domain.MediaImage && kind != domain.MediaVideo {
		writeError(w, http.StatusBadRequest, "kind must be \"image\" or \"video\"")
		return
	}
	if req.Key == "" || req.PublicURL == "" {
		writeError(w, http.StatusBadRequest, "key and publicUrl are required")
		return
	}

	var width, height *int
	if kind == domain.MediaImage {
		if w, h, ok := probeImageDimensions(req.PublicURL); ok {
			width, height = &w, &h
		}
	}

	asset, err := h.repo.Create(r.Context(), domain.MediaAsset{
		UserID:          userID,
		Kind:            kind,
		StorageKey:      req.Key,
		PublicURL:       req.PublicURL,
		ContentType:     req.ContentType,
		SizeBytes:       req.SizeBytes,
		DurationSeconds: req.DurationSeconds,
		Width:           width,
		Height:          height,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not record upload")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":        asset.ID,
		"publicUrl": asset.PublicURL,
		"kind":      asset.Kind,
	})
}

// probeImageDimensions reads only as much of the uploaded object as
// image.DecodeConfig needs for its header, not the whole file, so this stays
// cheap even for multi-megabyte images. A failure here (unsupported format,
// network hiccup) just means dimensions stay unknown — it never blocks the
// upload, since ValidateSettings treats unknown dimensions as unverifiable
// rather than as a rejection.
func probeImageDimensions(publicURL string) (width, height int, ok bool) {
	resp, err := http.Get(publicURL)
	if err != nil {
		return 0, 0, false
	}
	defer resp.Body.Close()

	cfg, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0, false
	}
	return cfg.Width, cfg.Height, true
}
