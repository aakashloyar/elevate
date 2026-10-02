CREATE TABLE IF NOT EXISTS email_otps (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES auth_credentials(user_id) ON DELETE CASCADE,
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    consumed_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS email_otps_active_idx
    ON email_otps(user_id, expires_at)
    WHERE consumed_at IS NULL;
