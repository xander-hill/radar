ALTER TABLE plans 
ADD COLUMN IF NOT EXISTS last_checkin_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Indexing the timestamp helps if we want to query "only active in last 24h" later
CREATE INDEX IF NOT EXISTS idx_plans_last_checkin ON plans(last_checkin_at);