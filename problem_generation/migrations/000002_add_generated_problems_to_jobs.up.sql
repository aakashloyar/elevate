ALTER TABLE generation_jobs
ADD COLUMN IF NOT EXISTS generated_problem_count INTEGER NOT NULL DEFAULT 0;
