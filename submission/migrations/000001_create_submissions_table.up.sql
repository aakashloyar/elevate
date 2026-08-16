CREATE TABLE IF NOT EXISTS submissions (
    id TEXT PRIMARY KEY,
    assessment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    duration_seconds INTEGER NOT NULL CHECK (duration_seconds > 0),
    expires_at TIMESTAMP,
    submitted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_submissions_expiration
    ON submissions (expires_at)
    WHERE status = 'IN_PROGRESS';
