package domain

type Chunk struct {
	UploadID   string
	ChunkNumber int
	ETag       string
	CreatedAt  int64
	UpdatedAt  int64
}
