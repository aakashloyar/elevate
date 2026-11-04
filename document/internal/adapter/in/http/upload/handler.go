package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type Handler struct {
	initUploadService     in.InitUploadService
	completeUploadService in.CompleteUploadService
	listDocumentsService  in.ListDocumentsService
	getDocumentService    in.GetDocumentService
	getContentService     in.GetDocumentContentService
}

type InitUploadRequest struct {
	CreatedBy   string `json:"created_by"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
}

type InitUploadResponse struct {
	UploadID     string      `json:"upload_id"`
	MaxChunkSize int64       `json:"max_chunk_size"`
	UploadURLs   []UploadURL `json:"chunks"`
}
type UploadURL struct {
	ChunkNumber int    `json:"chunk_number"`
	URL         string `json:"url"`
}

type CompleteUploadRequest struct {
	UploadID string `json:"upload_id"`
}

type CompleteUploadResponse struct{}

type ListDocumentsResponse struct {
	Documents []DocumentResponse `json:"documents"`
	Total     int                `json:"total"`
}

type DocumentResponse struct {
	ID          string              `json:"id"`
	CreatedBy   string              `json:"created_by"`
	FileName    string              `json:"file_name"`
	ContentType string              `json:"content_type"`
	FileSize    int64               `json:"file_size"`
	Status      domain.UploadStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type GetDocumentContentResponse struct {
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
}

func NewHandler(initUploadService in.InitUploadService, completeUploadService in.CompleteUploadService, listDocumentsService in.ListDocumentsService, getDocumentService in.GetDocumentService, getContentService in.GetDocumentContentService) *Handler {
	return &Handler{
		initUploadService:     initUploadService,
		completeUploadService: completeUploadService,
		listDocumentsService:  listDocumentsService,
		getDocumentService:    getDocumentService,
		getContentService:     getContentService,
	}
}

func (h *Handler) GetDocumentContent(w http.ResponseWriter, r *http.Request, documentID string) {
	result, err := h.getContentService.Execute(r.Context(), in.GetDocumentContentInput{DocumentID: documentID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response := GetDocumentContentResponse{DocumentID: result.DocumentID, Content: result.Content}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request, documentID string) {
	result, err := h.getDocumentService.Execute(r.Context(), in.GetDocumentInput{DocumentID: documentID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	document := result.Document
	response := DocumentResponse{
		ID: document.ID, CreatedBy: document.CreatedBy, FileName: document.FileName,
		ContentType: document.ContentType, FileSize: document.FileSize, Status: document.Status,
		CreatedAt: document.CreatedAt, UpdatedAt: document.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) InitUpload(w http.ResponseWriter, r *http.Request) {
	var req InitUploadRequest
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
	res := InitUploadResponse{
		UploadID:     out.UploadID,
		MaxChunkSize: out.MaxChunkSize,
		UploadURLs:   uploadURLs,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	var req CompleteUploadRequest
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

	res := CompleteUploadResponse{}
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
	}
	if value := strings.TrimSpace(r.URL.Query().Get("status")); value != "" {
		status := domain.UploadStatus(value)
		input.Status = &status
	}
	result, err := h.listDocumentsService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := ListDocumentsResponse{
		Documents: make([]DocumentResponse, 0, len(result.Documents)),
		Total:     result.Total,
	}
	for _, document := range result.Documents {
		response.Documents = append(response.Documents, DocumentResponse{
			ID: document.ID, CreatedBy: document.CreatedBy, FileName: document.FileName,
			ContentType: document.ContentType, FileSize: document.FileSize, Status: document.Status,
			CreatedAt: document.CreatedAt, UpdatedAt: document.UpdatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
