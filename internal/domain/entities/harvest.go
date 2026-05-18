package entities

import "time"

type Harvest struct {
	ID               uint      `gorm:"primaryKey"`
	PlantingID       uint      `gorm:"not null;index;uniqueIndex"` // 1 colheita por plantio
	FieldID          uint      `gorm:"not null;index"`
	HarvestDate      time.Time `gorm:"not null"`
	ProductivityBagHa float64  // produtividade em sc/ha
	ProductivityKgHa  float64  // produtividade em kg/ha
	TotalBags        float64   // total de sacas colhidas
	GrainMoisture    float64   // umidade dos grãos (%)
	Impurity         float64   // impureza (%)
	FieldLoss        float64   // perdas no campo (%)
	Notes            string
	CreatedBy        uint `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Planting Planting `gorm:"foreignKey:PlantingID"`
	Field    Field    `gorm:"foreignKey:FieldID"`
}
