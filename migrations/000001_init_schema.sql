-- 1. Enable Geospatial capabilities
CREATE EXTENSION IF NOT EXISTS postgis;

-- 2. Create the Plans table
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    location GEOGRAPHY(Point, 4326),
    base_score FLOAT DEFAULT 0,
    check_ins INT DEFAULT 0,
    saves INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create a Spatial Index (Resume Flex: This makes searches 100x faster)
CREATE INDEX IF NOT EXISTS idx_plans_location ON plans USING GIST (location);

-- 4. Add some "Seed Data" for testing
INSERT INTO plans (title, description, location, base_score)
VALUES 
('Arts District Pop-up', 'Secret gallery opening', ST_SetSRID(ST_MakePoint(-118.2352, 34.0415), 4326), 50.0),
('Erewhon Pulse', 'High activity at the juice bar', ST_SetSRID(ST_MakePoint(-118.3617, 34.0768), 4326), 20.0),
('Santa Monica Sunset Yoga', 'Beachfront session', ST_SetSRID(ST_MakePoint(-118.4912, 34.0115), 4326), 35.0);