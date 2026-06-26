ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE users
SET is_active = CASE
    WHEN status = 0 THEN FALSE
    ELSE TRUE
END;

DROP INDEX IF EXISTS idx_users_status;

ALTER TABLE users
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS telp,
    DROP COLUMN IF EXISTS address;

CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
