-- Rollback: Remove performance indexes

-- USERS TABLE
DROP INDEX IF EXISTS idx_users_active;

-- BUDGETS TABLE
DROP INDEX IF EXISTS idx_budgets_category_period;
DROP INDEX IF EXISTS idx_budgets_user_period_detail;

-- CATEGORIES TABLE
DROP INDEX IF EXISTS idx_categories_user_active_optimized;
DROP INDEX IF EXISTS idx_categories_user_slug;

-- ACCOUNTS TABLE
DROP INDEX IF EXISTS idx_accounts_user_type;
DROP INDEX IF EXISTS idx_accounts_user_active;

-- TRANSACTIONS TABLE
DROP INDEX IF EXISTS idx_transactions_active_recurring;
DROP INDEX IF EXISTS idx_transactions_recurrence_end;
DROP INDEX IF EXISTS idx_transactions_parent_id;
DROP INDEX IF EXISTS idx_transactions_is_recurring;
DROP INDEX IF EXISTS idx_transactions_account_date;
DROP INDEX IF EXISTS idx_transactions_user_type_date;
DROP INDEX IF EXISTS idx_transactions_user_date;
