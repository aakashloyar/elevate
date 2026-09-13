CREATE TABLE IF NOT EXISTS submission_answer_drafts (
    submission_id TEXT NOT NULL,
    problem_id TEXT NOT NULL,
    problem_type TEXT NOT NULL,
    option_ids TEXT[] NOT NULL DEFAULT '{}',
    option_texts TEXT[] NOT NULL DEFAULT '{}',
    answer TEXT[] NOT NULL DEFAULT '{}',
    answer_updated_at TIMESTAMP NULL,
    PRIMARY KEY (submission_id, problem_id),
    FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE
);
