-- Rollback: Drop accounts table
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP INDEX IF EXISTS idx_accounts_type;
DROP INDEX IF EXISTS idx_accounts_is_active;
DROP INDEX IF EXISTS idx_accounts_deleted_at;
DROP INDEX IF EXISTS idx_accounts_user_id_context;
DROP INDEX IF EXISTS idx_accounts_user_id;
DROP TABLE IF EXISTS accounts;
