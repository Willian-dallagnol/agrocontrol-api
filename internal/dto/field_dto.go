package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreateFieldRequest struct {
	Name     string                  `json:"name"      binding:"required,min=2,max=100"`
	Area     float64                 `json:"area"      binding:"required,gt=0"`
	SoilType string                  `json:"soil_type"`
	Status   entities.FieldStatus    `json:"status"    binding:"omitempty,oneof=active inactive fallow"`
	FarmID   uint                    `json:"farm_id"   binding:"required"`
}

type UpdateFieldRequest struct {
	Name     string                  `json:"name"      binding:"required,min=2,max=100"`
	Area     float64                 `json:"area"      binding:"required,gt=0"`
	SoilType string                  `json:"soil_type"`
	Status   entities.FieldStatus    `json:"status"    binding:"required,oneof=active inactive fallow"`
	FarmID   uint                    `json:"farm_id"   binding:"required"`
}

type FieldResponse struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	Area      float64              `json:"area"`
	SoilType  string               `json:"soil_type"`
	Status    entities.FieldStatus `json:"status"`
	FarmID    uint                 `json:"farm_id"`
	CreatedBy uint                 `json:"created_by"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
