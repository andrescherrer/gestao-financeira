-- Rollback: Drop budgets table
DROP TRIGGER IF EXISTS update_budgets_updated_at ON budgets;
DROP INDEX IF EXISTS idx_budgets_active;
DROP INDEX IF EXISTS idx_budgets_user_period;
DROP INDEX IF EXISTS idx_budgets_category_id;
DROP INDEX IF EXISTS idx_budgets_user_id;
DROP TABLE IF EXISTS budgets;
