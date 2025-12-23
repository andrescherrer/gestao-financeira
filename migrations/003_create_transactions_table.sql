-- Migration: Create transactions table
-- Created: 2025-12-23
-- Description: Creates the transactions table for the Transaction context

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- INCOME, EXPENSE
    amount BIGINT NOT NULL, -- Amount in cents
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL', -- Currency code (ISO 4217)
    description VARCHAR(500) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign key constraints
    CONSTRAINT fk_transactions_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transactions_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    
    -- Check constraints
    CONSTRAINT chk_transactions_type CHECK (type IN ('INCOME', 'EXPENSE')),
    CONSTRAINT chk_transactions_currency CHECK (currency IN ('BRL', 'USD', 'EUR')),
    CONSTRAINT chk_transactions_amount CHECK (amount > 0)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_account_id ON transactions(user_id, account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_type ON transactions(user_id, type);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to table
COMMENT ON TABLE transactions IS 'Stores transaction information for the Transaction context';
COMMENT ON COLUMN transactions.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN transactions.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN transactions.account_id IS 'Foreign key to accounts table';
COMMENT ON COLUMN transactions.type IS 'Transaction type: INCOME or EXPENSE';
COMMENT ON COLUMN transactions.amount IS 'Transaction amount in cents (BIGINT to support large amounts)';
COMMENT ON COLUMN transactions.currency IS 'Currency code (ISO 4217): BRL, USD, or EUR';
COMMENT ON COLUMN transactions.description IS 'Transaction description (3-500 characters)';
COMMENT ON COLUMN transactions.date IS 'Date when the transaction occurred';
COMMENT ON COLUMN transactions.created_at IS 'Timestamp when the transaction was created';
COMMENT ON COLUMN transactions.updated_at IS 'Timestamp when the transaction was last updated';
COMMENT ON COLUMN transactions.deleted_at IS 'Timestamp when the transaction was soft deleted (NULL if not deleted)';

