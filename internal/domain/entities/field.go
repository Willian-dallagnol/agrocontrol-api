package entities

import "time"

type Field struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"not null"`
	Area      float64 `gorm:"not null"`
	SoilType  string
	FarmID    uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
