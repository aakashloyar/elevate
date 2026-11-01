CREATE TABLE IF NOT EXISTS assessment_documents (
    assessment_id TEXT NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    document_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (assessment_id, document_id)
);

CREATE INDEX IF NOT EXISTS assessment_documents_document_id_idx
    ON assessment_documents (document_id);
