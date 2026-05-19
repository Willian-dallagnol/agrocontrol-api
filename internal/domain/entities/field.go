package entities

import "time"

// FieldStatus define o estado operacional do talhão
type FieldStatus string

const (
	FieldStatusActive   FieldStatus = "active"
	FieldStatusInactive FieldStatus = "inactive"
	FieldStatusFallow   FieldStatus = "fallow" // pousio
)

type Field struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"not null;index"`
	Area      float64 `gorm:"not null"`
	SoilType  string
	Status    FieldStatus `gorm:"not null;default:'active';index"`
	FarmID    uint        `gorm:"not null;index"`
	CreatedBy uint        `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Farm      Farm `gorm:"foreignKey:FarmID"`
}

// IsActive retorna true se o talhão está operacional para plantio.
func (f *Field) IsActive() bool {
	return f.Status == FieldStatusActive
}

// IsAvailableForPlanting verifica se o talhão pode receber um novo plantio.
// Talhões em pousio (fallow) podem ser plantados; inativos não.
func (f *Field) IsAvailableForPlanting() bool {
	return f.Status == FieldStatusActive || f.Status == FieldStatusFallow
}

// AreaInSquareMeters converte a área de hectares para metros quadrados.
func (f *Field) AreaInSquareMeters() float64 {
	return f.Area * 10_000
}

// Deactivate marca o talhão como inativo.
func (f *Field) Deactivate() {
	f.Status = FieldStatusInactive
}

// SetFallow coloca o talhão em pousio.
func (f *Field) SetFallow() {
	f.Status = FieldStatusFallow
}
