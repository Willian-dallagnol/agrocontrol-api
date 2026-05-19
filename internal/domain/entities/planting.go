package entities

import "time"

type PlantingStatus string

const (
	PlantingStatusActive    PlantingStatus = "active"
	PlantingStatusHarvested PlantingStatus = "harvested"
	PlantingStatusLost      PlantingStatus = "lost"
)

type Planting struct {
	ID              uint      `gorm:"primaryKey"`
	FieldID         uint      `gorm:"not null;index"`
	SeasonID        uint      `gorm:"not null;index"`
	CropID          uint      `gorm:"not null;index"`
	PlantingDate    time.Time `gorm:"not null"`
	ExpectedHarvest time.Time
	SeedsUsedKg     float64
	DensityKgHa     float64
	DepthCm         float64
	Spacing         float64
	Responsible     string
	Status          PlantingStatus `gorm:"not null;default:'active';index"`
	Notes           string
	CreatedBy       uint `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Field           Field  `gorm:"foreignKey:FieldID"`
	Season          Season `gorm:"foreignKey:SeasonID"`
	Crop            Crop   `gorm:"foreignKey:CropID"`
}

// IsActive retorna true se o plantio está em andamento.
func (p *Planting) IsActive() bool {
	return p.Status == PlantingStatusActive
}

// IsHarvested retorna true se o plantio já foi colhido.
func (p *Planting) IsHarvested() bool {
	return p.Status == PlantingStatusHarvested
}

// IsLate retorna true se a data esperada de colheita já passou e o plantio
// ainda está ativo — sinal de atraso na colheita.
func (p *Planting) IsLate() bool {
	if p.ExpectedHarvest.IsZero() || !p.IsActive() {
		return false
	}
	return time.Now().After(p.ExpectedHarvest)
}

// DaysUntilHarvest retorna quantos dias faltam para a colheita esperada.
// Retorna valor negativo se já passou da data.
func (p *Planting) DaysUntilHarvest() int {
	if p.ExpectedHarvest.IsZero() {
		return 0
	}
	return int(time.Until(p.ExpectedHarvest).Hours() / 24)
}

// MarkHarvested encerra o plantio como colhido.
func (p *Planting) MarkHarvested() {
	p.Status = PlantingStatusHarvested
}

// MarkLost registra o plantio como perdido.
func (p *Planting) MarkLost() {
	p.Status = PlantingStatusLost
}

// TotalSeedsForArea calcula o total de sementes necessárias para uma área em ha.
func (p *Planting) TotalSeedsForArea(areaHa float64) float64 {
	if p.DensityKgHa <= 0 || areaHa <= 0 {
		return 0
	}
	return p.DensityKgHa * areaHa
}
