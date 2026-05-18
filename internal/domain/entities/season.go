package entities

import "time"

type SeasonStatus string

const (
	SeasonStatusPlanning SeasonStatus = "planning"
	SeasonStatusActive   SeasonStatus = "active"
	SeasonStatusFinished SeasonStatus = "finished"
)

type Season struct {
	ID        uint         `gorm:"primaryKey"`
	Name      string       `gorm:"not null;index"`
	StartDate time.Time    `gorm:"not null"`
	EndDate   time.Time    `gorm:"not null"`
	Status    SeasonStatus `gorm:"not null;default:'planning';index"`
	CreatedBy uint         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
