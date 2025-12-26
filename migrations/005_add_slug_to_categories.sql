-- Migration: Add slug column to categories table
-- Date: 2025-01-27

-- Add slug column (nullable first to allow population)
ALTER TABLE categories
ADD COLUMN IF NOT EXISTS slug VARCHAR(100);

-- Generate slugs for existing categories based on their names
-- This will create slugs from existing category names
UPDATE categories
SET slug = LOWER(
  REGEXP_REPLACE(
    REGEXP_REPLACE(
      REGEXP_REPLACE(
        TRANSLATE(name, 'áàâãäéèêëíìîïóòôõöúùûüçñÁÀÂÃÄÉÈÊËÍÌÎÏÓÒÔÕÖÚÙÛÜÇÑ', 'aaaaaeeeeiiiioooouuuucnAAAAAEEEEIIIIOOOOUUUUCN'),
        '[^a-z0-9 ]', '', 'g'
      ),
      ' +', '-', 'g'
    ),
    '^-|-$', '', 'g'
  )
)
WHERE slug IS NULL;

-- Set default for any remaining NULL values
UPDATE categories
SET slug = 'categoria-' || SUBSTRING(id::text, 1, 8)
WHERE slug IS NULL;

-- Make slug NOT NULL after populating existing records
ALTER TABLE categories
ALTER COLUMN slug SET NOT NULL;

-- Create unique index for user_id and slug combination (only for non-deleted records)
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_user_slug ON categories(user_id, slug) WHERE deleted_at IS NULL;

