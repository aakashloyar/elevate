CREATE TABLE IF NOT EXISTS assessment_marking_schemes (
    assessment_id TEXT PRIMARY KEY REFERENCES assessments(id) ON DELETE CASCADE,
    single_correct_marks DOUBLE PRECISION NOT NULL,
    single_incorrect_marks DOUBLE PRECISION NOT NULL,
    single_skipped_marks DOUBLE PRECISION NOT NULL,
    multiple_correct_marks DOUBLE PRECISION NOT NULL,
    multiple_incorrect_marks DOUBLE PRECISION NOT NULL,
    multiple_skipped_marks DOUBLE PRECISION NOT NULL,
    numerical_correct_marks DOUBLE PRECISION NOT NULL,
    numerical_incorrect_marks DOUBLE PRECISION NOT NULL,
    numerical_skipped_marks DOUBLE PRECISION NOT NULL
);
