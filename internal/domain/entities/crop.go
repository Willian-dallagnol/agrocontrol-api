package entities

import "time"

type Crop struct {
	ID              uint    `gorm:"primaryKey"`
	Name            string  `gorm:"not null;index"`
	Variety         string  // cultivar/variedade
	Type            string  // grão, leguminosa, cereal
	CycleDays       int     // ciclo em dias
	SpacingCm       float64 // espaçamento entre linhas (cm)
	PlantPopulation int     // plantas por hectare
	CreatedBy       uint    `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
