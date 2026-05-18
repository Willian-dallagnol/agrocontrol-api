CREATE TABLE IF NOT EXISTS farms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_name VARCHAR(255) NOT NULL,
    location TEXT,
    total_area DECIMAL(10,2) NOT NULL CHECK (total_area > 0),
    city VARCHAR(255) NOT NULL,
    state CHAR(2) NOT NULL,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_farms_name ON farms(name);
CREATE INDEX IF NOT EXISTS idx_farms_created_by ON farms(created_by);