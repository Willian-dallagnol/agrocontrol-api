package entities

import "time"

type Farm struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	OwnerName string `gorm:"not null"`
	Location  string
	TotalArea float64 `gorm:"not null"`
	City      string  `gorm:"not null"`
	State     string  `gorm:"not null"`
	CreatedBy uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
