package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreateApplicationRequest struct {
	FieldID         uint                    `json:"field_id"          binding:"required"`
	PlantingID      *uint                   `json:"planting_id"`
	InputID         uint                    `json:"input_id"          binding:"required"`
	ApplicationType entities.ApplicationType `json:"application_type"  binding:"required,oneof=fertilizer herbicide fungicide insecticide correctant biological"`
	ApplicationDate time.Time               `json:"application_date"  binding:"required"`
	DosePerHa       float64                 `json:"dose_per_ha"       binding:"required,gt=0"`
	SprayVolume     float64                 `json:"spray_volume"      binding:"omitempty,gt=0"`
	Target          string                  `json:"target"`
	Equipment       string                  `json:"equipment"`
	Operator        string                  `json:"operator"          binding:"required"`
	WindSpeed       float64                 `json:"wind_speed"        binding:"omitempty,gte=0"`
	Temperature     float64                 `json:"temperature"       binding:"omitempty"`
	Humidity        float64                 `json:"humidity"          binding:"omitempty,gte=0,lte=100"`
	Notes           string                  `json:"notes"`
}

type ApplicationResponse struct {
	ID              uint                    `json:"id"`
	FieldID         uint                    `json:"field_id"`
	FieldName       string                  `json:"field_name,omitempty"`
	PlantingID      *uint                   `json:"planting_id"`
	InputID         uint                    `json:"input_id"`
	InputName       string                  `json:"input_name,omitempty"`
	ApplicationType entities.ApplicationType `json:"application_type"`
	ApplicationDate time.Time               `json:"application_date"`
	DosePerHa       float64                 `json:"dose_per_ha"`
	TotalUsed       float64                 `json:"total_used"`
	SprayVolume     float64                 `json:"spray_volume"`
	Target          string                  `json:"target"`
	Equipment       string                  `json:"equipment"`
	Operator        string                  `json:"operator"`
	WindSpeed       float64                 `json:"wind_speed"`
	Temperature     float64                 `json:"temperature"`
	Humidity        float64                 `json:"humidity"`
	Notes           string                  `json:"notes"`
	CreatedBy       uint                    `json:"created_by"`
	CreatedAt       time.Time               `json:"created_at"`
}
