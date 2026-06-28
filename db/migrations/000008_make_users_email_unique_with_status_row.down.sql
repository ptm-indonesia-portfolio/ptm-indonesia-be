DROP INDEX IF EXISTS idx_users_email_status_row_unique;

ALTER TABLE users
    DROP COLUMN IF EXISTS status_row;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active_unique
    ON users (LOWER(email))
    WHERE deleted_at IS NULL;
