-- Migration: Drop goals table
-- Created: 2025-12-31
-- Description: Drops the goals table for the Goal context

-- Drop trigger
DROP TRIGGER IF EXISTS update_goals_updated_at ON goals;

-- Drop indexes
DROP INDEX IF EXISTS idx_goals_status;
DROP INDEX IF EXISTS idx_goals_deleted_at;
DROP INDEX IF EXISTS idx_goals_deadline;
DROP INDEX IF EXISTS idx_goals_user_id_status;
DROP INDEX IF EXISTS idx_goals_user_id_context;
DROP INDEX IF EXISTS idx_goals_user_id;

-- Drop table
DROP TABLE IF EXISTS goals;

