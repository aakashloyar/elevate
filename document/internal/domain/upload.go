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
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	DocumentID string    `json:"document_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	StorageKey string    `json:"storage_key"`
	UploadedAt time.Time `json:"uploaded_at"`
}
