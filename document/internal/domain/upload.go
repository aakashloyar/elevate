package domain

import "time"

type UploadStatus string

const (
	StatusInitiated UploadStatus = "INITIATED"
	StatusMerging   UploadStatus = "MERGING"
	StatusStored    UploadStatus = "STORED"
	StatusFailed    UploadStatus = "FAILED"
)

type Upload struct {
	ID              string
	FileName        string
	FileSize        int64
	Status          UploadStatus
	StorageUploadID string
	StorageKey      string
	TotalChunks     int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type DocumentUploadedEvent struct {
	DocumentID string `json:"document_id"`
}
