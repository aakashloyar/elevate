package postgres

import (
	"database/sql"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type UploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) out.UploadRespository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) Init(upload domain.Upload) error {
	defer observeDB("uploads.init")()
	query := `
		INSERT INTO uploads {
			id,
			artist_ids,
			album_id,
			cover_image_url,
			file_name,
			file_size,
			status,
			storage_upload_id,
			storage_key,
			total_chunks,
			created_at,
			updated_at
		}
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(
		query,
		upload.ID,
		upload.FileName,
		upload.FileSize,
		upload.Status,
		upload.StorageUploadID,
		upload.StorageKey,
		upload.TotalChunks,
		upload.CreatedAt,
		upload.UpdatedAt,
	)
	return err
}

func (r *UploadRepository) Complete(uploadID string, status domain.UploadStatus) error {
	defer observeDB("uploads.complete")()
	query := `
		Update uploads SET status = $1
		WHERE id = $2 
	`
	_, err := r.db.Exec(
		query,
		uploadID,
		status,
	)
	return err
}

func (r *UploadRepository) FindByID(uploadID string) (domain.Upload, error) {
	defer observeDB("uploads.find_by_id")()
	query := `
		SELECT 
			id,
			artist_id,
			file_name,
			file_size,
			status,
			storage_upload_id,
			storage_key,
			total_chunks,
			created_at,
			updated_at
		FROM uploads
		WHERE id = $1	
	`
	row := r.db.QueryRow(query, uploadID)

	var upload domain.Upload
	err := row.Scan(
		&upload.ID,
		&upload.FileName,
		&upload.FileSize,
		&upload.Status,
		&upload.StorageUploadID,
		&upload.StorageKey,
		&upload.TotalChunks,
		&upload.CreatedAt,
		&upload.UpdatedAt,
	)

	if err != nil {
		return domain.Upload{}, err
	}
	return upload, nil
}

func (r *UploadRepository) GetChunksByUpload(uploadID string) ([]domain.Chunk, error) {
	defer observeDB("chunks.find_by_upload")()
	query := `
		SELECT
			upload_id,
			chunk_number,
			etag
		FROM chunks
		WHERE upload_id = $1
	`

	rows, err := r.db.Query(query, uploadID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []domain.Chunk

	for rows.Next() {
		var chunk domain.Chunk
		err := rows.Scan(
			&chunk.UploadID,
			&chunk.ChunkNumber,
			&chunk.ETag,
		)
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

func (r *UploadRepository) DeleteUploadParts(uploadID string) error {
	defer observeDB("upload_parts.delete")()
	query := `
		DELETE FROM upload_parts
		WHERE upload_id = $1
	`
	_, err := r.db.Exec(query, uploadID)
	return err
}
