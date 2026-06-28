ALTER TABLE users
    ADD COLUMN IF NOT EXISTS status_row SMALLINT NULL;

UPDATE users
SET status_row = 1
WHERE deleted_at IS NULL
  AND status_row IS NULL;

UPDATE users
SET status_row = NULL
WHERE deleted_at IS NOT NULL;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_email_key;

DROP INDEX IF EXISTS idx_users_email_active_unique;
DROP INDEX IF EXISTS idx_users_email_status_row_unique;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_status_row_unique
    ON users (email, status_row);
