CREATE TABLE IF NOT EXISTS fields (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    area DECIMAL(10,2) NOT NULL CHECK (area > 0),
    soil_type VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    farm_id INTEGER NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fields_farm_id ON fields(farm_id);
CREATE INDEX IF NOT EXISTS idx_fields_status ON fields(status);
CREATE INDEX IF NOT EXISTS idx_fields_created_by ON fields(created_by);
