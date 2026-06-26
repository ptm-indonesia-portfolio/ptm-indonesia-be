ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user';

UPDATE users
SET role = CASE
    WHEN status = 1 THEN 'super_admin'
    ELSE 'user'
END;

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
