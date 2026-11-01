CREATE TABLE IF NOT EXISTS chunks (
    upload_id TEXT NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    chunk_number INTEGER NOT NULL CHECK (chunk_number > 0),
    etag TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (upload_id, chunk_number)
);
