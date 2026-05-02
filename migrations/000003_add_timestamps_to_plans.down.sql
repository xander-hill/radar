DROP INDEX IF EXISTS idx_plans_last_checkin;
ALTER TABLE plans DROP COLUMN IF EXISTS last_checkin_at;
ALTER TABLE plans DROP COLUMN IF EXISTS created_at;