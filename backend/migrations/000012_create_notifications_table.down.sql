-- Drop trigger
DROP TRIGGER IF EXISTS trigger_update_notifications_updated_at ON notifications;

-- Drop function
DROP FUNCTION IF EXISTS update_notifications_updated_at();

-- Drop table
DROP TABLE IF EXISTS notifications;

