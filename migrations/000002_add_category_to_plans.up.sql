DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='plans' AND column_name='category') THEN
        ALTER TABLE plans ADD COLUMN category TEXT DEFAULT 'general';
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_plans_category ON plans(category);