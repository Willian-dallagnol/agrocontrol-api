package dto

// ── Queries (parâmetros de entrada) ─────────────────────────────────────────

type ProductivityReportQuery struct {
	SeasonID uint   `form:"season_id"`
	FarmID   uint   `form:"farm_id"`
	CropID   uint   `form:"crop_id"`
	Page     int    `form:"page,default=1"  binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=50" binding:"omitempty,min=1,max=200"`
}

func (q *ProductivityReportQuery) Offset() int { return (q.Page - 1) * q.Limit }

type CostPerFieldQuery struct {
	FieldID  uint   `form:"field_id"`
	SeasonID uint   `form:"season_id"`
	FarmID   uint   `form:"farm_id"`
}

// ── Relatório 1: Produtividade por talhão/safra ──────────────────────────────

type ProductivityReportItem struct {
	// Identificadores
	PlantingID uint   `json:"planting_id"`
	FieldID    uint   `json:"field_id"`
	FieldName  string `json:"field_name"`
	FarmID     uint   `json:"farm_id"`
	FarmName   string `json:"farm_name"`
	SeasonID   uint   `json:"season_id"`
	SeasonName string `json:"season_name"`
	CropID     uint   `json:"crop_id"`
	CropName   string `json:"crop_name"`
	Variety    string `json:"variety"`

	// Área e plantio
	PlantedAreaHa float64 `json:"planted_area_ha"` // área do talhão em ha
	PlantingDate  string  `json:"planting_date"`
	HarvestDate   string  `json:"harvest_date"`

	// Produção
	TotalBags          float64 `json:"total_bags"`           // sacas colhidas
	ProductivityBagHa  float64 `json:"productivity_bag_ha"`  // sc/ha
	ProductivityKgHa   float64 `json:"productivity_kg_ha"`   // kg/ha
	GrainMoisture      float64 `json:"grain_moisture_pct"`   // umidade %
	FieldLossPct       float64 `json:"field_loss_pct"`       // perdas %

	// Indicadores calculados
	EstimatedTotalKg   float64 `json:"estimated_total_kg"`   // total kg = kg/ha × área
	AboveAverageArea   bool    `json:"above_average"`        // acima da média da safra
}

type ProductivityReportResponse struct {
	Items       []ProductivityReportItem `json:"items"`
	Total       int64                   `json:"total"`
	Page        int                     `json:"page"`
	Limit       int                     `json:"limit"`
	TotalPages  int                     `json:"total_pages"`
	// Totalizadores da safra/filtro
	Summary ProductivitySummary `json:"summary"`
}

type ProductivitySummary struct {
	SeasonName          string  `json:"season_name,omitempty"`
	TotalFields         int     `json:"total_fields"`
	TotalPlantedAreaHa  float64 `json:"total_planted_area_ha"`
	TotalBags           float64 `json:"total_bags"`
	AvgProductivityBagHa float64 `json:"avg_productivity_bag_ha"`
	AvgProductivityKgHa  float64 `json:"avg_productivity_kg_ha"`
	BestField           string  `json:"best_field"`
	BestFieldBagHa      float64 `json:"best_field_bag_ha"`
}

// ── Relatório 2: Custo por talhão ────────────────────────────────────────────

type CostPerFieldItem struct {
	FieldID   uint   `json:"field_id"`
	FieldName string `json:"field_name"`
	FarmName  string `json:"farm_name"`
	AreaHa    float64 `json:"area_ha"`

	// Custo por categoria de insumo
	CostFertilizer  float64 `json:"cost_fertilizer"`
	CostHerbicide   float64 `json:"cost_herbicide"`
	CostFungicide   float64 `json:"cost_fungicide"`
	CostInsecticide float64 `json:"cost_insecticide"`
	CostCorrectant  float64 `json:"cost_correctant"`
	CostBiological  float64 `json:"cost_biological"`
	CostOther       float64 `json:"cost_other"`

	// Totais
	TotalCost       float64 `json:"total_cost"`
	CostPerHa       float64 `json:"cost_per_ha"`
	TotalApplications int   `json:"total_applications"`

	// Se tiver colheita na safra: custo por saca
	CostPerBag      float64 `json:"cost_per_bag,omitempty"`
	TotalBags       float64 `json:"total_bags,omitempty"`
}

type CostPerFieldResponse struct {
	Items   []CostPerFieldItem `json:"items"`
	Summary CostSummary        `json:"summary"`
}

type CostSummary struct {
	TotalFields     int     `json:"total_fields"`
	TotalAreaHa     float64 `json:"total_area_ha"`
	TotalCost       float64 `json:"total_cost"`
	AvgCostPerHa    float64 `json:"avg_cost_per_ha"`
	MostExpensiveField string `json:"most_expensive_field"`
	MostExpensiveCost  float64 `json:"most_expensive_cost_per_ha"`
}

// ── Dashboard consolidado ────────────────────────────────────────────────────

type DashboardOverviewResponse struct {
	// Fazendas e área
	TotalFarms       int64   `json:"total_farms"`
	TotalFields      int64   `json:"total_fields"`
	TotalAreaHa      float64 `json:"total_area_ha"`
	PlantedAreaHa    float64 `json:"planted_area_ha"`

	// Safras e plantios
	TotalSeasons     int64   `json:"total_seasons"`
	ActivePlantings  int64   `json:"active_plantings"`
	HarvestedThisYear int64  `json:"harvested_this_year"`

	// Produtividade
	AvgProductivityBagHa float64 `json:"avg_productivity_bag_ha"`
	TotalBagsThisYear    float64 `json:"total_bags_this_year"`

	// Insumos e estoque
	TotalInputTypes  int64   `json:"total_input_types"`
	LowStockInputs   int64   `json:"low_stock_inputs"`
	ExpiringInputs   int64   `json:"expiring_inputs_30d"`

	// Aplicações
	TotalApplications     int64   `json:"total_applications"`
	ApplicationsThisMonth int64   `json:"applications_this_month"`
	EstimatedTotalCost    float64 `json:"estimated_total_cost"`

	// Alertas
	OpenAlerts       int64   `json:"open_alerts"`
	CriticalAlerts   int64   `json:"critical_alerts"`

	// Top performers
	TopFields []TopFieldItem `json:"top_fields"`
}

type TopFieldItem struct {
	FieldName         string  `json:"field_name"`
	FarmName          string  `json:"farm_name"`
	ProductivityBagHa float64 `json:"productivity_bag_ha"`
}
