CREATE TABLE IF NOT EXISTS plantings (
    id SERIAL PRIMARY KEY,
    field_id INTEGER NOT NULL REFERENCES fields(id) ON DELETE CASCADE,
    season_id INTEGER NOT NULL REFERENCES seasons(id),
    crop_id INTEGER NOT NULL REFERENCES crops(id),
    planting_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expected_harvest TIMESTAMP WITH TIME ZONE,
    seeds_used_kg DECIMAL(10,2),
    density_kg_ha DECIMAL(10,2),
    depth_cm DECIMAL(10,2),
    spacing VARCHAR(100),
    responsible VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_plantings_field_id ON plantings(field_id);
CREATE INDEX IF NOT EXISTS idx_plantings_season_id ON plantings(season_id);
CREATE INDEX IF NOT EXISTS idx_plantings_status ON plantings(status);
CREATE INDEX IF NOT EXISTS idx_plantings_field_season ON plantings(field_id, season_id);