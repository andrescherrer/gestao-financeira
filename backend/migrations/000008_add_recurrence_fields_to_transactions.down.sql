-- Rollback: Remove recurrence fields from transactions table
DROP INDEX IF EXISTS idx_transactions_recurring_active;
DROP INDEX IF EXISTS idx_transactions_is_recurring;
DROP INDEX IF EXISTS idx_transactions_parent_transaction_id;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_parent_transaction_id;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS chk_transactions_recurrence_frequency;
ALTER TABLE transactions DROP COLUMN IF EXISTS parent_transaction_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS recurrence_end_date;
ALTER TABLE transactions DROP COLUMN IF EXISTS recurrence_frequency;
ALTER TABLE transactions DROP COLUMN IF EXISTS is_recurring;
