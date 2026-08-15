CREATE TABLE IF NOT EXISTS problems (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    statement TEXT NOT NULL,
    type TEXT NOT NULL
        CHECK (type IN ('single', 'multiple', 'numerical')),
    difficulty TEXT NOT NULL
        CHECK (difficulty IN ('easy', 'medium', 'hard')),
    source_type TEXT NOT NULL
        CHECK (source_type IN ('manual', 'ai')),
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);