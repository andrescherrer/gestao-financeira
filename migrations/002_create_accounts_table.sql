-- Migration: Create accounts table
-- Created: 2025-12-23
-- Description: Creates the accounts table for the Account Management context

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- BANK, WALLET, INVESTMENT, CREDIT_CARD
    balance BIGINT NOT NULL DEFAULT 0, -- Amount in cents
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL', -- Currency code (ISO 4217)
    context VARCHAR(20) NOT NULL, -- PERSONAL, BUSINESS
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign key constraint
    CONSTRAINT fk_accounts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Check constraints
    CONSTRAINT chk_accounts_type CHECK (type IN ('BANK', 'WALLET', 'INVESTMENT', 'CREDIT_CARD')),
    CONSTRAINT chk_accounts_context CHECK (context IN ('PERSONAL', 'BUSINESS')),
    CONSTRAINT chk_accounts_currency CHECK (currency IN ('BRL', 'USD', 'EUR'))
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id_context ON accounts(user_id, context);
CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounts_is_active ON accounts(is_active);
CREATE INDEX IF NOT EXISTS idx_accounts_type ON accounts(type);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to table
COMMENT ON TABLE accounts IS 'Stores account information for the Account Management context';
COMMENT ON COLUMN accounts.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN accounts.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN accounts.name IS 'Account name (3-100 characters)';
COMMENT ON COLUMN accounts.type IS 'Account type: BANK, WALLET, INVESTMENT, or CREDIT_CARD';
COMMENT ON COLUMN accounts.balance IS 'Account balance in cents (BIGINT to support large amounts)';
COMMENT ON COLUMN accounts.currency IS 'Currency code (ISO 4217): BRL, USD, or EUR';
COMMENT ON COLUMN accounts.context IS 'Account context: PERSONAL or BUSINESS';
COMMENT ON COLUMN accounts.is_active IS 'Whether the account is active';
COMMENT ON COLUMN accounts.created_at IS 'Timestamp when the account was created';
COMMENT ON COLUMN accounts.updated_at IS 'Timestamp when the account was last updated';
COMMENT ON COLUMN accounts.deleted_at IS 'Timestamp when the account was soft deleted (NULL if not deleted)';

