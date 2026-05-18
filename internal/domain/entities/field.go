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
	ID        uint        `gorm:"primaryKey"`
	Name      string      `gorm:"not null;index"`
	Area      float64     `gorm:"not null"`
	SoilType  string
	Status    FieldStatus `gorm:"not null;default:'active';index"`
	FarmID    uint        `gorm:"not null;index"`
	CreatedBy uint        `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Farm Farm `gorm:"foreignKey:FarmID"`
}
