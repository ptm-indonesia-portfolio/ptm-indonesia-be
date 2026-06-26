ALTER TABLE users
    ALTER COLUMN avatar_url DROP NOT NULL,
    ALTER COLUMN avatar_url DROP DEFAULT,
    ALTER COLUMN address DROP NOT NULL,
    ALTER COLUMN address DROP DEFAULT,
    ALTER COLUMN telp DROP NOT NULL,
    ALTER COLUMN telp DROP DEFAULT;

UPDATE users
SET avatar_url = NULL
WHERE avatar_url = '';

UPDATE users
SET address = NULL
WHERE address = '';

UPDATE users
SET telp = NULL
WHERE telp = '';
