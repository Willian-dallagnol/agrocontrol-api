package entities

import "time"

const BagWeightKg = 60.0 // peso padrão de uma saca em kg

type Harvest struct {
	ID                uint      `gorm:"primaryKey"`
	PlantingID        uint      `gorm:"not null;index;uniqueIndex"`
	FieldID           uint      `gorm:"not null;index"`
	HarvestDate       time.Time `gorm:"not null"`
	ProductivityBagHa float64
	ProductivityKgHa  float64
	TotalBags         float64
	GrainMoisture     float64
	Impurity          float64
	FieldLoss         float64
	Notes             string
	CreatedBy         uint `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Planting          Planting `gorm:"foreignKey:PlantingID"`
	Field             Field    `gorm:"foreignKey:FieldID"`
}

// CalculateProductivity calcula produtividade em sc/ha e kg/ha a partir do total
// de sacas e da área do talhão. Retorna zero se área for inválida.
func (h *Harvest) CalculateProductivity(fieldAreaHa float64) {
	if fieldAreaHa <= 0 {
		return
	}
	h.ProductivityBagHa = h.TotalBags / fieldAreaHa
	h.ProductivityKgHa = h.ProductivityBagHa * BagWeightKg
}

// IsHighYield verifica se a produtividade está acima do limiar considerado alta
// para a cultura (padrão: 60 sc/ha para soja, referência de mercado brasileiro).
func (h *Harvest) IsHighYield(thresholdBagHa float64) bool {
	return h.ProductivityBagHa >= thresholdBagHa
}

// AdjustedTotalBags retorna o total de sacas descontando perdas no campo.
func (h *Harvest) AdjustedTotalBags() float64 {
	if h.FieldLoss <= 0 {
		return h.TotalBags
	}
	return h.TotalBags * (1 - h.FieldLoss/100)
}
