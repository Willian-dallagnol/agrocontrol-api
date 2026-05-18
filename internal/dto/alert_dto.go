package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreateAlertRequest struct {
	Title       string                `json:"title"       binding:"required,min=3"`
	Type        entities.AlertType    `json:"type"        binding:"required,oneof=low_stock expired_input pest weather pending_application harvest_soon"`
	Description string                `json:"description"`
	Priority    entities.AlertPriority `json:"priority"    binding:"omitempty,oneof=low medium high"`
	RefID       *uint                 `json:"ref_id"`
	RefType     string                `json:"ref_type"`
}

type UpdateAlertStatusRequest struct {
	Status entities.AlertStatus `json:"status" binding:"required,oneof=open resolved ignored"`
}

type AlertResponse struct {
	ID          uint                  `json:"id"`
	Title       string                `json:"title"`
	Type        entities.AlertType    `json:"type"`
	Description string                `json:"description"`
	Priority    entities.AlertPriority `json:"priority"`
	Status      entities.AlertStatus  `json:"status"`
	RefID       *uint                 `json:"ref_id"`
	RefType     string                `json:"ref_type"`
	CreatedBy   uint                  `json:"created_by"`
	ResolvedAt  *time.Time            `json:"resolved_at"`
	CreatedAt   time.Time             `json:"created_at"`
}
