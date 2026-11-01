package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type Handler struct {
	initUploadService     in.InitUploadService
	completeUploadService in.CompleteUploadService
	listDocumentsService  in.ListDocumentsService
	getDocumentsService   in.GetDocumentsService
}

type initUploadRequest struct {
	CreatedBy   string `json:"created_by"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
}

type initUploadResponse struct {
	UploadID     string      `json:"upload_id"`
	MaxChunkSize int64       `json:"max_chunk_size"`
	UploadURLs   []UploadURL `json:"chunks"`
}
type UploadURL struct {
	ChunkNumber int    `json:"chunk_number"`
	URL         string `json:"url"`
}

type completeUploadRequest struct {
	UploadID string `json:"upload_id"`
}

type completUploadResponse struct{}

func NewHandler(initUploadService in.InitUploadService, completeUploadService in.CompleteUploadService, listDocumentsService in.ListDocumentsService, getDocumentsService in.GetDocumentsService) *Handler {
	return &Handler{
		initUploadService:     initUploadService,
		completeUploadService: completeUploadService,
		listDocumentsService:  listDocumentsService,
		getDocumentsService:   getDocumentsService,
	}
}

func (h *Handler) InitUpload(w http.ResponseWriter, r *http.Request) {
	var req initUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input := in.InitUploadInput{
		CreatedBy: req.CreatedBy, FileName: req.FileName,
		ContentType: req.ContentType, FileSize: req.FileSize,
	}

	out, err := h.initUploadService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploadURLs := make([]UploadURL, 0, len(out.UploadURLs))
	for _, u := range out.UploadURLs {
		uploadURLs = append(uploadURLs, UploadURL{
			ChunkNumber: u.ChunkNumber,
			URL:         u.URL,
		})
	}
	res := initUploadResponse{
		UploadID:     out.UploadID,
		MaxChunkSize: out.MaxChunkSize,
		UploadURLs:   uploadURLs,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	var req completeUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input := in.CompleteUploadInput{
		UploadID: req.UploadID,
	}

	_, err := h.completeUploadService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := completUploadResponse{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	readyStatus := domain.StatusReady
	input := in.ListDocumentsInput{
		CreatedBy: r.URL.Query().Get("created_by"),
		Search:    r.URL.Query().Get("search"),
		Status:    &readyStatus,
		Limit:     10,
	}
	if value := strings.TrimSpace(r.URL.Query().Get("status")); value != "" {
		status := domain.UploadStatus(value)
		input.Status = &status
	}
	if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 100 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		input.Limit = parsed
	}
	if value := strings.TrimSpace(r.URL.Query().Get("offset")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		input.Offset = parsed
	}
	result, err := h.listDocumentsService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetDocumentsBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DocumentIDs []string `json:"document_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result, err := h.getDocumentsService.Execute(r.Context(), in.GetDocumentsInput{DocumentIDs: req.DocumentIDs})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
