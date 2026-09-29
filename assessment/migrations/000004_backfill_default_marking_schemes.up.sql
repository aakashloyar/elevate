INSERT INTO assessment_marking_schemes (
    assessment_id,
    single_correct_marks, single_incorrect_marks, single_skipped_marks,
    multiple_correct_marks, multiple_incorrect_marks, multiple_skipped_marks,
    numerical_correct_marks, numerical_incorrect_marks, numerical_skipped_marks
)
SELECT
    assessments.id,
    4, -1, 0,
    4, -2, 0,
    4, 0, 0
FROM assessments
WHERE NOT EXISTS (
    SELECT 1
    FROM assessment_marking_schemes
    WHERE assessment_marking_schemes.assessment_id = assessments.id
);
