-- Migration: Add performance indexes
-- Created: 2025-12-27
-- Description: Adds additional indexes to improve query performance for common operations

-- ============================================
-- TRANSACTIONS TABLE INDEXES
-- ============================================

-- Index for date range queries (used in reports)
-- Composite index for user_id + date (for date range filtering)
CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date DESC);

-- Index for user_id + type + date (common filter combination)
CREATE INDEX IF NOT EXISTS idx_transactions_user_type_date ON transactions(user_id, type, date DESC);

-- Index for account_id + date (for account-specific reports)
CREATE INDEX IF NOT EXISTS idx_transactions_account_date ON transactions(account_id, date DESC);

-- Index for recurring transactions queries
CREATE INDEX IF NOT EXISTS idx_transactions_is_recurring ON transactions(is_recurring) WHERE is_recurring = true;
CREATE INDEX IF NOT EXISTS idx_transactions_parent_id ON transactions(parent_transaction_id) WHERE parent_transaction_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_recurrence_end ON transactions(recurrence_end_date) WHERE recurrence_end_date IS NOT NULL;

-- Index for active recurring transactions (used by recurring transaction processor)
CREATE INDEX IF NOT EXISTS idx_transactions_active_recurring ON transactions(user_id, is_recurring, recurrence_end_date) 
    WHERE is_recurring = true AND (recurrence_end_date IS NULL OR recurrence_end_date >= CURRENT_DATE);

-- ============================================
-- ACCOUNTS TABLE INDEXES
-- ============================================

-- Index for user_id + is_active (common filter)
CREATE INDEX IF NOT EXISTS idx_accounts_user_active ON accounts(user_id, is_active) WHERE deleted_at IS NULL;

-- Index for user_id + type (filter by account type)
CREATE INDEX IF NOT EXISTS idx_accounts_user_type ON accounts(user_id, type) WHERE deleted_at IS NULL;

-- ============================================
-- CATEGORIES TABLE INDEXES
-- ============================================

-- Index for user_id + slug (unique lookup)
CREATE INDEX IF NOT EXISTS idx_categories_user_slug ON categories(user_id, slug) WHERE deleted_at IS NULL;

-- Index for user_id + is_active (common filter)
-- Note: idx_categories_user_active already exists, but adding WHERE clause for better performance
CREATE INDEX IF NOT EXISTS idx_categories_user_active_optimized ON categories(user_id, is_active) WHERE deleted_at IS NULL;

-- ============================================
-- BUDGETS TABLE INDEXES
-- ============================================

-- Index for user_id + period_type + year + month (common filter)
CREATE INDEX IF NOT EXISTS idx_budgets_user_period_detail ON budgets(user_id, period_type, year, month) WHERE deleted_at IS NULL;

-- Index for category_id + period (for budget progress calculations)
CREATE INDEX IF NOT EXISTS idx_budgets_category_period ON budgets(category_id, period_type, year, month) WHERE deleted_at IS NULL AND is_active = true;

-- ============================================
-- USERS TABLE INDEXES
-- ============================================

-- Index for email lookup (already exists, but ensuring it's there)
-- idx_users_email already exists in 001_create_users_table.sql

-- Index for active users lookup
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active) WHERE deleted_at IS NULL;

-- ============================================
-- COMMENTS
-- ============================================

COMMENT ON INDEX idx_transactions_user_date IS 'Optimizes queries filtering transactions by user and date range';
COMMENT ON INDEX idx_transactions_user_type_date IS 'Optimizes queries filtering transactions by user, type, and date';
COMMENT ON INDEX idx_transactions_account_date IS 'Optimizes queries filtering transactions by account and date';
COMMENT ON INDEX idx_transactions_active_recurring IS 'Optimizes queries for active recurring transactions';
COMMENT ON INDEX idx_accounts_user_active IS 'Optimizes queries for active accounts by user';
COMMENT ON INDEX idx_categories_user_slug IS 'Optimizes category lookup by user and slug';
COMMENT ON INDEX idx_budgets_user_period_detail IS 'Optimizes budget queries by user and period details';

