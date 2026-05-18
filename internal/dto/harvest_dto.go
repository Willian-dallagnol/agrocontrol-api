package dto

import "time"

type CreateHarvestRequest struct {
	PlantingID        uint      `json:"planting_id"          binding:"required"`
	HarvestDate       time.Time `json:"harvest_date"         binding:"required"`
	ProductivityBagHa float64   `json:"productivity_bag_ha"  binding:"omitempty,gte=0"`
	ProductivityKgHa  float64   `json:"productivity_kg_ha"   binding:"omitempty,gte=0"`
	TotalBags         float64   `json:"total_bags"           binding:"omitempty,gte=0"`
	GrainMoisture     float64   `json:"grain_moisture"       binding:"omitempty,gte=0,lte=100"`
	Impurity          float64   `json:"impurity"             binding:"omitempty,gte=0,lte=100"`
	FieldLoss         float64   `json:"field_loss"           binding:"omitempty,gte=0,lte=100"`
	Notes             string    `json:"notes"`
}

type HarvestResponse struct {
	ID                uint      `json:"id"`
	PlantingID        uint      `json:"planting_id"`
	FieldID           uint      `json:"field_id"`
	FieldName         string    `json:"field_name,omitempty"`
	HarvestDate       time.Time `json:"harvest_date"`
	ProductivityBagHa float64   `json:"productivity_bag_ha"`
	ProductivityKgHa  float64   `json:"productivity_kg_ha"`
	TotalBags         float64   `json:"total_bags"`
	GrainMoisture     float64   `json:"grain_moisture"`
	Impurity          float64   `json:"impurity"`
	FieldLoss         float64   `json:"field_loss"`
	Notes             string    `json:"notes"`
	CreatedBy         uint      `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
