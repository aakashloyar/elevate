CREATE TABLE IF NOT EXISTS generation_jobs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    single_correct_count INTEGER NOT NULL,
    multi_correct_count INTEGER NOT NULL,
    numerical_count INTEGER NOT NULL,
    document_id TEXT NULL,
    assessment_id TEXT NULL,
    level TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
