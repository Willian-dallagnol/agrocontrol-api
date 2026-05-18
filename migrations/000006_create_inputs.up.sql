CREATE TABLE IF NOT EXISTS inputs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    manufacturer VARCHAR(255),
    batch_number VARCHAR(100),
    expiration_date TIMESTAMP WITH TIME ZONE,
    unit VARCHAR(50) NOT NULL,
    stock_qty DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (stock_qty >= 0),
    min_stock_qty DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (min_stock_qty >= 0),
    cost_per_unit DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (cost_per_unit >= 0),
    active BOOLEAN NOT NULL DEFAULT true,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inputs_name ON inputs(name);
CREATE INDEX IF NOT EXISTS idx_inputs_category ON inputs(category);
CREATE INDEX IF NOT EXISTS idx_inputs_active ON inputs(active);
CREATE INDEX IF NOT EXISTS idx_inputs_created_by ON inputs(created_by);
