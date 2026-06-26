UPDATE users
SET avatar_url = ''
WHERE avatar_url IS NULL;

UPDATE users
SET address = ''
WHERE address IS NULL;

UPDATE users
SET telp = ''
WHERE telp IS NULL;

ALTER TABLE users
    ALTER COLUMN avatar_url SET DEFAULT '',
    ALTER COLUMN avatar_url SET NOT NULL,
    ALTER COLUMN address SET DEFAULT '',
    ALTER COLUMN address SET NOT NULL,
    ALTER COLUMN telp SET DEFAULT '',
    ALTER COLUMN telp SET NOT NULL;
