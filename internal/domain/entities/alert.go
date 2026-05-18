package entities

import "time"

type AlertType string

const (
	AlertTypeLowStock     AlertType = "low_stock"
	AlertTypeExpiredInput AlertType = "expired_input"
	AlertTypePest         AlertType = "pest"
	AlertTypeWeather      AlertType = "weather"
	AlertTypePending      AlertType = "pending_application"
	AlertTypeHarvestSoon  AlertType = "harvest_soon"
)

type AlertPriority string

const (
	AlertPriorityLow    AlertPriority = "low"
	AlertPriorityMedium AlertPriority = "medium"
	AlertPriorityHigh   AlertPriority = "high"
)

type AlertStatus string

const (
	AlertStatusOpen     AlertStatus = "open"
	AlertStatusResolved AlertStatus = "resolved"
	AlertStatusIgnored  AlertStatus = "ignored"
)

type Alert struct {
	ID          uint          `gorm:"primaryKey"`
	Title       string        `gorm:"not null"`
	Type        AlertType     `gorm:"not null;index"`
	Description string
	Priority    AlertPriority `gorm:"not null;default:'medium';index"`
	Status      AlertStatus   `gorm:"not null;default:'open';index"`
	RefID       *uint         // ID da entidade relacionada (opcional)
	RefType     string        // tipo da entidade relacionada (ex: "input", "field")
	CreatedBy   uint          `gorm:"not null"`
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
