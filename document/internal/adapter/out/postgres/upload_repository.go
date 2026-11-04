package postgres

import (
	"database/sql"
	"fmt"
	"strings"

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
		INSERT INTO uploads (
			id,
			created_by,
			file_name,
			content_type,
			file_size,
			status,
			storage_upload_id,
			storage_key,
			total_chunks,
			created_at,
			updated_at
        )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(
		query,
		upload.ID,
		upload.CreatedBy,
		upload.FileName,
		upload.ContentType,
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
		status,
		uploadID,
	)
	return err
}

func (r *UploadRepository) SaveExtractedText(uploadID, text string) error {
	defer observeDB("uploads.save_extracted_text")()
	_, err := r.db.Exec(`UPDATE uploads SET extracted_text = $1, updated_at = NOW() WHERE id = $2`, text, uploadID)
	return err
}

func (r *UploadRepository) FindByID(uploadID string) (domain.Upload, error) {
	defer observeDB("uploads.find_by_id")()
	query := `
		SELECT 
			id,
			created_by,
			file_name,
			content_type,
			file_size,
			status,
			storage_upload_id,
			storage_key,
			total_chunks,
			extracted_text,
			created_at,
			updated_at
		FROM uploads
		WHERE id = $1	
	`
	row := r.db.QueryRow(query, uploadID)

	var upload domain.Upload
	err := row.Scan(
		&upload.ID,
		&upload.CreatedBy,
		&upload.FileName,
		&upload.ContentType,
		&upload.FileSize,
		&upload.Status,
		&upload.StorageUploadID,
		&upload.StorageKey,
		&upload.TotalChunks,
		&upload.ExtractedText,
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
		DELETE FROM chunks
		WHERE upload_id = $1
	`
	_, err := r.db.Exec(query, uploadID)
	return err
}

func (r *UploadRepository) List(filters out.UploadListFilters) ([]domain.Upload, int, error) {
	defer observeDB("uploads.list")()

	where := []string{"1 = 1"}
	args := make([]any, 0, 4)
	argIndex := 1
	if filters.CreatedBy != nil {
		where = append(where, fmt.Sprintf("created_by = $%d", argIndex))
		args = append(args, *filters.CreatedBy)
		argIndex++
	}
	if filters.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}
	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		where = append(where, fmt.Sprintf("file_name ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(*filters.Search)+"%")
		argIndex++
	}

	countQuery := "SELECT COUNT(*) FROM uploads WHERE " + strings.Join(where, " AND ")
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, created_by, file_name, content_type, file_size, status,
		storage_upload_id, storage_key, total_chunks, created_at, updated_at
		FROM uploads WHERE ` + strings.Join(where, " AND ") + " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Upload, 0)
	for rows.Next() {
		var upload domain.Upload
		if err := rows.Scan(&upload.ID, &upload.CreatedBy, &upload.FileName, &upload.ContentType, &upload.FileSize, &upload.Status, &upload.StorageUploadID, &upload.StorageKey, &upload.TotalChunks, &upload.CreatedAt, &upload.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, upload)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *UploadRepository) FindByIDs(uploadIDs []string) ([]domain.Upload, error) {
	defer observeDB("uploads.find_by_ids")()
	if len(uploadIDs) == 0 {
		return []domain.Upload{}, nil
	}
	placeholders := make([]string, len(uploadIDs))
	args := make([]any, len(uploadIDs))
	for i, id := range uploadIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	rows, err := r.db.Query(`SELECT id, created_by, file_name, content_type, file_size, status,
		storage_upload_id, storage_key, total_chunks, created_at, updated_at
		FROM uploads WHERE id IN (`+strings.Join(placeholders, ",")+") ORDER BY id", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Upload, 0, len(uploadIDs))
	for rows.Next() {
		var upload domain.Upload
		if err := rows.Scan(&upload.ID, &upload.CreatedBy, &upload.FileName, &upload.ContentType, &upload.FileSize, &upload.Status, &upload.StorageUploadID, &upload.StorageKey, &upload.TotalChunks, &upload.CreatedAt, &upload.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, upload)
	}
	return items, rows.Err()
}
