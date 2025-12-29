-- Rollback: Drop categories table
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP INDEX IF EXISTS idx_categories_user_active;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_deleted_at;
DROP INDEX IF EXISTS idx_categories_user_id;
DROP TABLE IF EXISTS categories;
