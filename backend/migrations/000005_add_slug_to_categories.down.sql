-- Rollback: Remove slug column from categories table
DROP INDEX IF EXISTS idx_categories_user_slug;
ALTER TABLE categories DROP COLUMN IF EXISTS slug;
