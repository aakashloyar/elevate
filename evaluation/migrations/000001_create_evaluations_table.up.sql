CREATE TABLE IF NOT EXISTS evaluations (
    submission_id TEXT PRIMARY KEY,
    assessment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    score DOUBLE PRECISION NOT NULL,
    questions JSONB NOT NULL,
    evaluated_at TIMESTAMPTZ NOT NULL
);
