CREATE TABLE IF NOT EXISTS harvests (
    id SERIAL PRIMARY KEY,
    planting_id INTEGER NOT NULL REFERENCES plantings(id),
    field_id INTEGER NOT NULL REFERENCES fields(id),
    harvest_date TIMESTAMP WITH TIME ZONE NOT NULL,
    productivity_bag_ha DECIMAL(10,2),
    productivity_kg_ha DECIMAL(10,2),
    total_bags DECIMAL(10,2),
    grain_moisture DECIMAL(5,2),
    impurity DECIMAL(5,2),
    field_loss DECIMAL(5,2),
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_harvests_planting_id ON harvests(planting_id);
CREATE INDEX IF NOT EXISTS idx_harvests_field_id ON harvests(field_id);
CREATE INDEX IF NOT EXISTS idx_harvests_date ON harvests(harvest_date);
