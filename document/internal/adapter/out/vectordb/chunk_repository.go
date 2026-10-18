package vectordb

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type ChunkRepository struct {
	db *sql.DB
}

func NewChunkRepository(db *sql.DB) out.VectorRepository {
	return &ChunkRepository{db: db}
}

func (r *ChunkRepository) SaveChunk(chunk domain.DocumentChunk) error {
	if chunk.ID == "" || chunk.DocumentID == "" || strings.TrimSpace(chunk.Content) == "" {
		return fmt.Errorf("chunk id, document id, and content are required")
	}
	if len(chunk.Embedding) == 0 {
		return fmt.Errorf("chunk embedding is required")
	}

	metadata := chunk.Metadata
	if len(metadata) == 0 {
		metadata = []byte(`{}`)
	}
	createdAt := interface{}(chunk.CreatedAt)
	if chunk.CreatedAt.IsZero() {
		createdAt = nil
	}
	_, err := r.db.Exec(`
		INSERT INTO document_chunks (
			id, document_id, chunk_index, content, embedding, metadata, created_at
		)
		VALUES ($1, $2, $3, $4, $5::vector, $6::jsonb, COALESCE($7, NOW()))
		ON CONFLICT (document_id, chunk_index)
		DO UPDATE SET
			content = EXCLUDED.content,
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata
	`, chunk.ID, chunk.DocumentID, chunk.ChunkIndex, chunk.Content, vectorLiteral(chunk.Embedding), metadata, createdAt)
	return err
}

func (r *ChunkRepository) SearchSimilar(embedding []float32, limit int) ([]domain.DocumentChunk, error) {
	if len(embedding) == 0 {
		return nil, fmt.Errorf("query embedding is required")
	}
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.Query(`
		SELECT id, document_id, chunk_index, content, embedding::text, metadata, created_at
		FROM document_chunks
		ORDER BY embedding <=> $1::vector
		LIMIT $2
	`, vectorLiteral(embedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chunks := make([]domain.DocumentChunk, 0)
	for rows.Next() {
		var chunk domain.DocumentChunk
		var embeddingText string
		if err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.ChunkIndex, &chunk.Content, &embeddingText, &chunk.Metadata, &chunk.CreatedAt); err != nil {
			return nil, err
		}
		chunk.Embedding, err = parseVector(embeddingText)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return chunks, nil
}

func vectorLiteral(values []float32) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.FormatFloat(float64(value), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func parseVector(value string) ([]float32, error) {
	value = strings.TrimSpace(strings.Trim(value, "[]"))
	if value == "" {
		return []float32{}, nil
	}
	parts := strings.Split(value, ",")
	result := make([]float32, len(parts))
	for i, part := range parts {
		parsed, err := strconv.ParseFloat(strings.TrimSpace(part), 32)
		if err != nil {
			return nil, fmt.Errorf("parse vector value %q: %w", part, err)
		}
		result[i] = float32(parsed)
	}
	return result, nil
}
