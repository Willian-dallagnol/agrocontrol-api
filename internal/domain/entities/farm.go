package entities

import "time"

// Farm representa uma propriedade agrícola
type Farm struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"not null;index"`
	OwnerName string  `gorm:"not null"`
	Location  string
	TotalArea float64 `gorm:"not null;check:total_area > 0"`
	City      string  `gorm:"not null"`
	State     string  `gorm:"not null;size:2"`  // sigla do estado: PR, SP, etc.
	CreatedBy uint    `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Fields []Field `gorm:"foreignKey:FarmID;constraint:OnDelete:CASCADE"`
}
