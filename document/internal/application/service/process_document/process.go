package process_document

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type Service struct {
	uploads  out.UploadRespository
	storage  out.Storage
	embedder out.Embedder
	vectors  out.VectorRepository
	idGen    domain.IDGenerator
}

func NewService(uploads out.UploadRespository, storage out.Storage, embedder out.Embedder, vectors out.VectorRepository, idGen domain.IDGenerator) in.ProcessDocumentService {
	return &Service{uploads: uploads, storage: storage, embedder: embedder, vectors: vectors, idGen: idGen}
}

func (s *Service) Execute(ctx context.Context, event domain.DocumentUploadedEvent) error {
	if strings.TrimSpace(event.DocumentID) == "" {
		return fmt.Errorf("document id is required")
	}

	upload, err := s.uploads.FindByID(event.DocumentID)
	if err != nil {
		return fmt.Errorf("find uploaded document %s: %w", event.DocumentID, err)
	}
	if strings.TrimSpace(upload.StorageKey) == "" {
		return fmt.Errorf("uploaded document %s has no storage key", event.DocumentID)
	}

	body, err := s.storage.Download(ctx, upload.StorageKey)
	if err != nil {
		return fmt.Errorf("download document %s: %w", event.DocumentID, err)
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read document %s: %w", event.DocumentID, err)
	}
	if !validText(data) {
		return fmt.Errorf("document %s is not UTF-8 text; add a file extractor before embedding it", event.DocumentID)
	}

	chunks := splitText(string(data))
	if len(chunks) == 0 {
		return fmt.Errorf("document %s has no text to embed", event.DocumentID)
	}

	// Gemini's batch endpoint is used in groups so one large document does not
	// create an oversized request. The repository upsert makes retries safe.
	const batchSize = 100
	for start := 0; start < len(chunks); start += batchSize {
		end := start + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		embeddings, err := s.embedder.EmbedDocuments(ctx, chunks[start:end])
		if err != nil {
			return fmt.Errorf("embed document %s chunks %d-%d: %w", event.DocumentID, start, end-1, err)
		}
		if len(embeddings) != end-start {
			return fmt.Errorf("embedder returned %d vectors for %d chunks", len(embeddings), end-start)
		}

		for i, embedding := range embeddings {
			metadata, err := json.Marshal(map[string]any{
				"file_name": upload.FileName,
				"source":    "document.uploaded",
			})
			if err != nil {
				return fmt.Errorf("marshal chunk metadata: %w", err)
			}
			if err := s.vectors.SaveChunk(domain.DocumentChunk{
				ID: s.idGen.NewID(), 
				DocumentID: event.DocumentID,
				ChunkIndex: start + i,
				Content: chunks[start+i],
				Embedding: embedding,
				Metadata: metadata,
			}); err != nil {
				return fmt.Errorf("save document %s chunk %d: %w", event.DocumentID, start+i, err)
			}
		}
	}
	return nil
}
