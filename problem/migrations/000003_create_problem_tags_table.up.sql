CREATE TABLE IF NOT EXISTS problem_tags (
    problem_id UUID NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (problem_id, tag),
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);