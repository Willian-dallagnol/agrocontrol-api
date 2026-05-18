package entities

import "time"

type PlantingStatus string

const (
	PlantingStatusActive    PlantingStatus = "active"
	PlantingStatusHarvested PlantingStatus = "harvested"
	PlantingStatusLost      PlantingStatus = "lost"
)

type Planting struct {
	ID              uint           `gorm:"primaryKey"`
	FieldID         uint           `gorm:"not null;index"`
	SeasonID        uint           `gorm:"not null;index"`
	CropID          uint           `gorm:"not null;index"`
	PlantingDate    time.Time      `gorm:"not null"`
	ExpectedHarvest time.Time
	SeedsUsedKg     float64 // kg de sementes usadas
	DensityKgHa     float64 // kg por hectare
	DepthCm         float64 // profundidade de plantio
	Spacing         float64 // espaçamento entre linhas
	Responsible     string  // operador/técnico responsável
	Status          PlantingStatus `gorm:"not null;default:'active';index"`
	Notes           string
	CreatedBy       uint `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Field  Field  `gorm:"foreignKey:FieldID"`
	Season Season `gorm:"foreignKey:SeasonID"`
	Crop   Crop   `gorm:"foreignKey:CropID"`
}
