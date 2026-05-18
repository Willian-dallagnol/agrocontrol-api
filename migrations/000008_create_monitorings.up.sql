CREATE TABLE IF NOT EXISTS monitorings (
    id SERIAL PRIMARY KEY,
    field_id INTEGER NOT NULL REFERENCES fields(id) ON DELETE CASCADE,
    planting_id INTEGER REFERENCES plantings(id),
    inspection_date TIMESTAMP WITH TIME ZONE NOT NULL,
    type VARCHAR(100) NOT NULL,
    problem_name VARCHAR(255),
    infestation_level DECIMAL(5,2),
    severity VARCHAR(50),
    crop_stage VARCHAR(100),
    technical_rec TEXT,
    urgent BOOLEAN NOT NULL DEFAULT false,
    inspector VARCHAR(255),
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_monitorings_field_id ON monitorings(field_id);
CREATE INDEX IF NOT EXISTS idx_monitorings_date ON monitorings(inspection_date);
CREATE INDEX IF NOT EXISTS idx_monitorings_urgent ON monitorings(urgent);