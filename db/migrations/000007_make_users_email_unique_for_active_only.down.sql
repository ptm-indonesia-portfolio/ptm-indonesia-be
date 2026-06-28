DROP INDEX IF EXISTS idx_users_email_active_unique;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'users_email_key'
          AND conrelid = 'users'::regclass
    ) THEN
        ALTER TABLE users
            ADD CONSTRAINT users_email_key UNIQUE (email);
    END IF;
END $$;
