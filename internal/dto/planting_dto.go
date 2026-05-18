package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreatePlantingRequest struct {
	FieldID         uint                   `json:"field_id"          binding:"required"`
	SeasonID        uint                   `json:"season_id"         binding:"required"`
	CropID          uint                   `json:"crop_id"           binding:"required"`
	PlantingDate    time.Time              `json:"planting_date"     binding:"required"`
	ExpectedHarvest time.Time              `json:"expected_harvest"`
	SeedsUsedKg     float64                `json:"seeds_used_kg"     binding:"omitempty,gt=0"`
	DensityKgHa     float64                `json:"density_kg_ha"     binding:"omitempty,gt=0"`
	DepthCm         float64                `json:"depth_cm"          binding:"omitempty,gt=0"`
	Spacing         float64                `json:"spacing"           binding:"omitempty,gt=0"`
	Responsible     string                 `json:"responsible"`
	Notes           string                 `json:"notes"`
}

type UpdatePlantingRequest struct {
	ExpectedHarvest time.Time              `json:"expected_harvest"`
	Status          entities.PlantingStatus `json:"status" binding:"omitempty,oneof=active harvested lost"`
	Responsible     string                 `json:"responsible"`
	Notes           string                 `json:"notes"`
}

type PlantingResponse struct {
	ID              uint                   `json:"id"`
	FieldID         uint                   `json:"field_id"`
	FieldName       string                 `json:"field_name,omitempty"`
	SeasonID        uint                   `json:"season_id"`
	SeasonName      string                 `json:"season_name,omitempty"`
	CropID          uint                   `json:"crop_id"`
	CropName        string                 `json:"crop_name,omitempty"`
	PlantingDate    time.Time              `json:"planting_date"`
	ExpectedHarvest time.Time              `json:"expected_harvest"`
	SeedsUsedKg     float64                `json:"seeds_used_kg"`
	DensityKgHa     float64                `json:"density_kg_ha"`
	DepthCm         float64                `json:"depth_cm"`
	Spacing         float64                `json:"spacing"`
	Responsible     string                 `json:"responsible"`
	Status          entities.PlantingStatus `json:"status"`
	Notes           string                 `json:"notes"`
	CreatedBy       uint                   `json:"created_by"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
