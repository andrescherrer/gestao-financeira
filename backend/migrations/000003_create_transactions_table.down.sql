-- Rollback: Drop transactions table
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP INDEX IF EXISTS idx_transactions_deleted_at;
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_user_id_type;
DROP INDEX IF EXISTS idx_transactions_user_id_account_id;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP TABLE IF EXISTS transactions;
