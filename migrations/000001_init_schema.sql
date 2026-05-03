-- 1. Enable Geospatial capabilities
CREATE EXTENSION IF NOT EXISTS postgis;

-- 2. Create the Plans table
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    category TEXT DEFAULT 'general',
    description TEXT,
    location GEOGRAPHY(Point, 4326),
    base_score FLOAT DEFAULT 0,
    check_ins INT DEFAULT 0,
    saves INT DEFAULT 0,
    last_checkin_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- 3. Create a Spatial Index (Resume Flex: This makes searches 100x faster)
CREATE INDEX IF NOT EXISTS idx_plans_location ON plans USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_plans_active ON plans (deleted_at) WHERE deleted_at IS NULL;

-- 4. Add some "Seed Data" for testing
INSERT INTO plans (title, description, location, base_score)
VALUES 
('Coffman Union Pulse', 'Heavy student activity', ST_SetSRID(ST_MakePoint(-93.2354, 44.9744), 4326), 50.0),
('Dinkytown Coffee', 'Finals week energy', ST_SetSRID(ST_MakePoint(-93.2363, 44.9808), 4326), 25.0),
('Huntington Bank Stadium', 'Game day momentum', ST_SetSRID(ST_MakePoint(-93.2246, 44.9760), 4326), 40.0);