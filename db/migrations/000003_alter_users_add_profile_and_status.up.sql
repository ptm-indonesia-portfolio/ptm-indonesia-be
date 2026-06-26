ALTER TABLE users
    ADD COLUMN IF NOT EXISTS address TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS telp VARCHAR(30) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS status SMALLINT NOT NULL DEFAULT 0;

UPDATE users
SET status = CASE
    WHEN role = 'super_admin' THEN 1
    WHEN is_active = TRUE THEN 2
    ELSE 0
END
WHERE status = 0;

DROP INDEX IF EXISTS idx_users_is_active;

ALTER TABLE users
    DROP COLUMN IF EXISTS is_active;

CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
