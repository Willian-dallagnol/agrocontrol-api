package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreateSeasonRequest struct {
	Name      string               `json:"name"       binding:"required,min=2,max=100"`
	StartDate time.Time            `json:"start_date" binding:"required"`
	EndDate   time.Time            `json:"end_date"   binding:"required"`
	Status    entities.SeasonStatus `json:"status"     binding:"omitempty,oneof=planning active finished"`
}

type UpdateSeasonRequest struct {
	Name      string               `json:"name"       binding:"required,min=2,max=100"`
	StartDate time.Time            `json:"start_date" binding:"required"`
	EndDate   time.Time            `json:"end_date"   binding:"required"`
	Status    entities.SeasonStatus `json:"status"     binding:"required,oneof=planning active finished"`
}

type SeasonResponse struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	StartDate time.Time            `json:"start_date"`
	EndDate   time.Time            `json:"end_date"`
	Status    entities.SeasonStatus `json:"status"`
	CreatedBy uint                 `json:"created_by"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
