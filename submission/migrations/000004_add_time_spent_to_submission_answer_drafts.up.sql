ALTER TABLE submission_answer_drafts
    ADD COLUMN IF NOT EXISTS time_spent_seconds INTEGER NOT NULL DEFAULT 0
    CHECK (time_spent_seconds >= 0);
