CREATE TABLE IF NOT EXISTS evaluations (
    submission_id TEXT PRIMARY KEY,
    assessment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    duration_seconds INTEGER NOT NULL,
    submitted_at TIMESTAMPTZ,
    score DOUBLE PRECISION NOT NULL,
    questions JSONB NOT NULL,
    evaluated_at TIMESTAMPTZ NOT NULL
);
