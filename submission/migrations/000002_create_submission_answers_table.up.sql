CREATE TABLE IF NOT EXISTS submission_answers (
    id TEXT PRIMARY KEY,
    submission_id TEXT NOT NULL,
    problem_id TEXT NOT NULL,
    answer TEXT[] NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (submission_id, problem_id),
    FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
);
