package domain

import "time"

type UploadStatus string

const (
	StatusInitiated  UploadStatus = "INITIATED"
	StatusProcessing UploadStatus = "PROCESSING"
	StatusReady      UploadStatus = "READY"
	StatusFailed     UploadStatus = "FAILED"
)

type Upload struct {
	ID              string
	CreatedBy       string
	FileName        string
	ContentType     string
	FileSize        int64
	Status          UploadStatus
	StorageUploadID string
	StorageKey      string
	TotalChunks     int
	ExtractedText   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type DocumentUploadedEvent struct {
	DocumentID string `json:"document_id"`
}
