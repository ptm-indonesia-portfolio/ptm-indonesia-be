UPDATE users
SET status_row = NULL
WHERE deleted_at IS NOT NULL
  AND status_row IS NOT NULL;
