-- Migration: Insert default user
-- Date: 2025-01-27
-- Description: Inserts a default user (André Scherrer) for development/testing purposes
-- This user will be created automatically when the database is recreated

-- Insert default user if it doesn't already exist
-- Using a fixed UUID for consistency across database recreations
INSERT INTO users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active,
    created_at,
    updated_at
) VALUES (
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890'::uuid,
    'andrescherrer@gmail.com',
    '$2a$10$lLuJcNnNMGlX/nin/zM5sO02kUHuVk0Vzc77dt0IwCIi5xLGz1f7m', -- Hash for: @!Rafa2021
    'André',
    'Scherrer',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

-- Add comment
COMMENT ON TABLE users IS 'Stores user account information for the Identity context. Default user: andrescherrer@gmail.com';

