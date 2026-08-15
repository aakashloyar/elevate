CREATE TABLE IF NOT EXISTS problem_options (
    id UUID PRIMARY KEY,
    problem_id UUID NOT NULL,
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);