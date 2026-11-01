package postgres

import (
	"database/sql"

	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type ChunkRepository struct {
	db *sql.DB
}

func NewChunkRepository(db *sql.DB) out.ChunkRespository {
	return &ChunkRepository{db: db}
}

func (r *ChunkRepository) Mark(chunk domain.Chunk) error {
	defer observeDB("chunks.mark")()
	query := `
		INSERT INTO chunks (
			upload_id,
			chunk_number,
			etag,
			created_at,
			updated_at
		}
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		query,
		chunk.UploadID,
		chunk.ChunkNumber,
		chunk.ETag,
		chunk.CreatedAt,
		chunk.UpdatedAt,
	)
	return err
}
