-- Migration: Add recurrence fields to transactions table
-- Created: 2025-12-27
-- Description: Adds fields to support recurring transactions

-- Add recurrence fields to transactions table
ALTER TABLE transactions
ADD COLUMN IF NOT EXISTS is_recurring BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS recurrence_frequency VARCHAR(20) NULL,
ADD COLUMN IF NOT EXISTS recurrence_end_date DATE NULL,
ADD COLUMN IF NOT EXISTS parent_transaction_id UUID NULL;

-- Add check constraint for recurrence_frequency
ALTER TABLE transactions
ADD CONSTRAINT chk_transactions_recurrence_frequency 
CHECK (recurrence_frequency IS NULL OR recurrence_frequency IN ('DAILY', 'WEEKLY', 'MONTHLY', 'YEARLY'));

-- Add foreign key constraint for parent_transaction_id
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_parent_transaction_id 
FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id) ON DELETE SET NULL;

-- Add index for parent_transaction_id
CREATE INDEX IF NOT EXISTS idx_transactions_parent_transaction_id ON transactions(parent_transaction_id);

-- Add index for recurring transactions
CREATE INDEX IF NOT EXISTS idx_transactions_is_recurring ON transactions(is_recurring) WHERE is_recurring = true;

-- Add composite index for finding active recurring transactions
CREATE INDEX IF NOT EXISTS idx_transactions_recurring_active ON transactions(is_recurring, recurrence_end_date) WHERE is_recurring = true;
