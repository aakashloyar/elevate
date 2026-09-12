UPDATE submissions
SET started_at = created_at
WHERE started_at IS NULL;

ALTER TABLE submissions
    ALTER COLUMN started_at SET NOT NULL;
