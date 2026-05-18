CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    field_id INTEGER NOT NULL REFERENCES fields(id) ON DELETE CASCADE,
    planting_id INTEGER REFERENCES plantings(id),
    input_id INTEGER NOT NULL REFERENCES inputs(id),
    application_type VARCHAR(100) NOT NULL,
    application_date TIMESTAMP WITH TIME ZONE NOT NULL,
    dose_per_ha DECIMAL(10,2) NOT NULL,
    total_used DECIMAL(10,2) NOT NULL,
    spray_volume DECIMAL(10,2),
    target VARCHAR(255),
    equipment VARCHAR(255),
    operator VARCHAR(255),
    wind_speed DECIMAL(5,2),
    temperature DECIMAL(5,2),
    humidity DECIMAL(5,2),
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_field_id ON applications(field_id);
CREATE INDEX IF NOT EXISTS idx_applications_input_id ON applications(input_id);
CREATE INDEX IF NOT EXISTS idx_applications_date ON applications(application_date);
CREATE INDEX IF NOT EXISTS idx_applications_created_by ON applications(created_by);