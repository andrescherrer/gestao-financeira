-- Migration: Create categories table
-- Created: 2025-12-26
-- Description: Creates the categories table for the Category context

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);
CREATE INDEX IF NOT EXISTS idx_categories_user_active ON categories(user_id, is_active);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comment to table
COMMENT ON TABLE categories IS 'Stores category information for transactions';
COMMENT ON COLUMN categories.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN categories.user_id IS 'User who owns this category (foreign key to users.id)';
COMMENT ON COLUMN categories.name IS 'Category name (e.g., Alimentação, Transporte)';
COMMENT ON COLUMN categories.description IS 'Optional description of the category';
COMMENT ON COLUMN categories.is_active IS 'Whether the category is active';
COMMENT ON COLUMN categories.created_at IS 'Timestamp when the category was created';
COMMENT ON COLUMN categories.updated_at IS 'Timestamp when the category was last updated';
COMMENT ON COLUMN categories.deleted_at IS 'Timestamp when the category was soft deleted (NULL if not deleted)';

