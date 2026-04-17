package dto

type CreateFieldRequest struct {
	Name     string  `json:"name" binding:"required"`
	Area     float64 `json:"area"`
	SoilType string  `json:"soil_type"`
	FarmID   uint    `json:"farm_id" binding:"required"`
}

type UpdateFieldRequest struct {
	Name     string  `json:"name" binding:"required"`
	Area     float64 `json:"area"`
	SoilType string  `json:"soil_type"`
	FarmID   uint    `json:"farm_id" binding:"required"`
}

type FieldResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Area     float64 `json:"area"`
	SoilType string  `json:"soil_type"`
	FarmID   uint    `json:"farm_id"`
}
