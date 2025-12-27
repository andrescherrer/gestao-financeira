-- Migration: Create budgets table
-- Date: 2025-12-27

CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    amount BIGINT NOT NULL, -- Amount in cents
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    period_type VARCHAR(10) NOT NULL CHECK (period_type IN ('MONTHLY', 'YEARLY')),
    year INTEGER NOT NULL CHECK (year >= 1900 AND year <= 3000),
    month INTEGER CHECK (month IS NULL OR (month >= 1 AND month <= 12)),
    context VARCHAR(20) NOT NULL CHECK (context IN ('PERSONAL', 'BUSINESS')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_budgets_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_budgets_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    
    -- Unique constraint: one budget per user, category, period, year, and month (if monthly)
    CONSTRAINT uq_budgets_user_category_period UNIQUE (user_id, category_id, period_type, year, month)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);
CREATE INDEX IF NOT EXISTS idx_budgets_category_id ON budgets(category_id);
CREATE INDEX IF NOT EXISTS idx_budgets_user_period ON budgets(user_id, year, month);
CREATE INDEX IF NOT EXISTS idx_budgets_active ON budgets(user_id, is_active) WHERE deleted_at IS NULL;

-- Add comment to table
COMMENT ON TABLE budgets IS 'Stores budget information for users';
COMMENT ON COLUMN budgets.amount IS 'Budget amount in cents';
COMMENT ON COLUMN budgets.period_type IS 'Budget period type: MONTHLY or YEARLY';
COMMENT ON COLUMN budgets.month IS 'Month number (1-12) for monthly budgets, NULL for yearly budgets';
COMMENT ON COLUMN budgets.year IS 'Year for the budget period';

