CREATE TABLE IF NOT EXISTS assessment_problems (
    assessment_id TEXT NOT NULL,
    problem_id TEXT NOT NULL,
    PRIMARY KEY (assessment_id, problem_id)
);