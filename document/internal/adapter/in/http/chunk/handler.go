package http

import (
	"encoding/json"
	"github.com/aakashloyar/elevate/document/internal/application/ports/in/chunk"
	"net/http"
)

type Handler struct {
	markChunkService in.MarkChunkService
}

type markChunkRequest struct {
	UploadID    string `json:"artist_id"`
	ChunkNumber int    `json:"file_name"`
	ETag        string `json:"file_size"`
}

type markChunkResponse struct {
}

func NewHandler(markChunkService in.MarkChunkService) *Handler {
	return &Handler{
		markChunkService: markChunkService,
	}
}

func (h *Handler) MarkChunk(w http.ResponseWriter, r *http.Request) {
	var req markChunkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	input := in.MarkChunkInput{
		UploadID:    req.UploadID,
		ChunkNumber: req.ChunkNumber,
		ETag:        req.ETag,
	}

	_, err := h.markChunkService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := markChunkResponse{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
