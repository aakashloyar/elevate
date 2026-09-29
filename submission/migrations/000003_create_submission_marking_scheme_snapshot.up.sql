ALTER TABLE submissions
    ADD COLUMN IF NOT EXISTS marking_scheme JSONB NOT NULL DEFAULT '{"assessment_id":"","single":{"correct":4,"incorrect":-1,"skipped":0},"multiple":{"correct":4,"incorrect":-2,"skipped":0},"numerical":{"correct":4,"incorrect":0,"skipped":0}}'::jsonb;
