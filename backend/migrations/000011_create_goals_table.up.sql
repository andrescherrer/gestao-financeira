-- Migration: Create goals table
-- Created: 2025-12-31
-- Description: Creates the goals table for the Goal context

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create goals table
CREATE TABLE IF NOT EXISTS goals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    target_amount BIGINT NOT NULL,
    target_currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    current_amount BIGINT NOT NULL DEFAULT 0,
    current_currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    deadline DATE NOT NULL,
    context VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'IN_PROGRESS',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign key constraint
    CONSTRAINT fk_goals_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Check constraints
    CONSTRAINT chk_goals_context CHECK (context IN ('PERSONAL', 'BUSINESS')),
    CONSTRAINT chk_goals_status CHECK (status IN ('IN_PROGRESS', 'COMPLETED', 'OVERDUE', 'CANCELLED')),
    CONSTRAINT chk_goals_currency CHECK (target_currency IN ('BRL', 'USD', 'EUR')),
    CONSTRAINT chk_goals_target_amount CHECK (target_amount > 0)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_user_id_context ON goals(user_id, context);
CREATE INDEX IF NOT EXISTS idx_goals_user_id_status ON goals(user_id, status);
CREATE INDEX IF NOT EXISTS idx_goals_deadline ON goals(deadline);
CREATE INDEX IF NOT EXISTS idx_goals_deleted_at ON goals(deleted_at);
CREATE INDEX IF NOT EXISTS idx_goals_status ON goals(status);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_goals_updated_at
    BEFORE UPDATE ON goals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to table
COMMENT ON TABLE goals IS 'Stores goal information for the Goal context';
COMMENT ON COLUMN goals.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN goals.user_id IS 'User who owns this goal (foreign key to users.id)';
COMMENT ON COLUMN goals.name IS 'Goal name (e.g., Viagem para Europa, Reserva de Emergência)';
COMMENT ON COLUMN goals.target_amount IS 'Target amount in cents';
COMMENT ON COLUMN goals.target_currency IS 'Currency code (ISO 4217)';
COMMENT ON COLUMN goals.current_amount IS 'Current amount in cents';
COMMENT ON COLUMN goals.current_currency IS 'Currency code (ISO 4217)';
COMMENT ON COLUMN goals.deadline IS 'Goal deadline date';
COMMENT ON COLUMN goals.context IS 'Account context (PERSONAL, BUSINESS)';
COMMENT ON COLUMN goals.status IS 'Goal status (IN_PROGRESS, COMPLETED, OVERDUE, CANCELLED)';
COMMENT ON COLUMN goals.created_at IS 'Timestamp when the goal was created';
COMMENT ON COLUMN goals.updated_at IS 'Timestamp when the goal was last updated';
COMMENT ON COLUMN goals.deleted_at IS 'Timestamp when the goal was soft deleted (NULL if not deleted)';

