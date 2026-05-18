package dto

import "time"

type CreateCropRequest struct {
	Name            string  `json:"name"             binding:"required,min=2,max=100"`
	Variety         string  `json:"variety"`
	Type            string  `json:"type"`
	CycleDays       int     `json:"cycle_days"       binding:"omitempty,gt=0"`
	SpacingCm       float64 `json:"spacing_cm"       binding:"omitempty,gt=0"`
	PlantPopulation int     `json:"plant_population" binding:"omitempty,gt=0"`
}

type UpdateCropRequest struct {
	Name            string  `json:"name"             binding:"required,min=2,max=100"`
	Variety         string  `json:"variety"`
	Type            string  `json:"type"`
	CycleDays       int     `json:"cycle_days"       binding:"omitempty,gt=0"`
	SpacingCm       float64 `json:"spacing_cm"       binding:"omitempty,gt=0"`
	PlantPopulation int     `json:"plant_population" binding:"omitempty,gt=0"`
}

type CropResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Variety         string    `json:"variety"`
	Type            string    `json:"type"`
	CycleDays       int       `json:"cycle_days"`
	SpacingCm       float64   `json:"spacing_cm"`
	PlantPopulation int       `json:"plant_population"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
