-- ============================================================
-- AgroControl API — Índices para otimização dos relatórios
-- Execute no banco agro_control após o AutoMigrate
-- ============================================================

-- Relatório de Produtividade
-- JOIN: plantings → fields → farms → seasons → crops → harvests
CREATE INDEX IF NOT EXISTS idx_plantings_field_season   ON plantings(field_id, season_id);
CREATE INDEX IF NOT EXISTS idx_plantings_season_status  ON plantings(season_id, status);
CREATE INDEX IF NOT EXISTS idx_harvests_planting        ON harvests(planting_id);
CREATE INDEX IF NOT EXISTS idx_harvests_field_date      ON harvests(field_id, harvest_date DESC);
CREATE INDEX IF NOT EXISTS idx_harvests_productivity    ON harvests(productivity_bag_ha DESC);

-- Relatório de Custo por Talhão
-- JOIN: applications → fields → farms → inputs
CREATE INDEX IF NOT EXISTS idx_applications_field_date  ON applications(field_id, application_date DESC);
CREATE INDEX IF NOT EXISTS idx_applications_input       ON applications(input_id);
CREATE INDEX IF NOT EXISTS idx_applications_type        ON applications(application_type);

-- Segurança multi-usuário (JOIN Field → Farm → created_by)
CREATE INDEX IF NOT EXISTS idx_farms_created_by         ON farms(created_by);
CREATE INDEX IF NOT EXISTS idx_fields_farm              ON fields(farm_id);
CREATE INDEX IF NOT EXISTS idx_plantings_field          ON plantings(field_id);
CREATE INDEX IF NOT EXISTS idx_applications_field       ON applications(field_id);
CREATE INDEX IF NOT EXISTS idx_monitorings_field        ON monitorings(field_id);
CREATE INDEX IF NOT EXISTS idx_harvests_field           ON harvests(field_id);

-- Insumos — alertas de estoque e vencimento
CREATE INDEX IF NOT EXISTS idx_inputs_stock             ON inputs(stock_qty, min_stock_qty) WHERE active = true;
CREATE INDEX IF NOT EXISTS idx_inputs_expiration        ON inputs(expiration_date)          WHERE active = true AND expiration_date IS NOT NULL;

-- Alertas
CREATE INDEX IF NOT EXISTS idx_alerts_status_priority   ON alerts(status, priority);
CREATE INDEX IF NOT EXISTS idx_alerts_created_by        ON alerts(created_by, status);
