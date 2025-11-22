DROP INDEX IF EXISTS idx_holdings_deleted_at;
ALTER TABLE holdings DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE holdings DROP COLUMN IF EXISTS updated_at;
ALTER TABLE holdings DROP COLUMN IF EXISTS created_at;
