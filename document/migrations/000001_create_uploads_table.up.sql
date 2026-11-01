CREATE TABLE IF NOT EXISTS uploads (
    id TEXT PRIMARY KEY,
    created_by TEXT NOT NULL DEFAULT '',
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'application/octet-stream',
    file_size BIGINT NOT NULL CHECK (file_size >= 0),
    status TEXT NOT NULL,
    storage_upload_id TEXT NOT NULL,
    storage_key TEXT NOT NULL UNIQUE,
    total_chunks INTEGER NOT NULL CHECK (total_chunks > 0),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS uploads_created_by_status_idx
    ON uploads (created_by, status, created_at DESC);
