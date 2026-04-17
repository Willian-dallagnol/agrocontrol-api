package dto

type CreateFarmRequest struct {
	Name      string  `json:"name" binding:"required"`
	OwnerName string  `json:"owner_name" binding:"required"`
	Location  string  `json:"location"`
	TotalArea float64 `json:"total_area"`
	City      string  `json:"city" binding:"required"`
	State     string  `json:"state" binding:"required"`
}

type UpdateFarmRequest struct {
	Name      string  `json:"name" binding:"required"`
	OwnerName string  `json:"owner_name" binding:"required"`
	Location  string  `json:"location"`
	TotalArea float64 `json:"total_area"`
	City      string  `json:"city" binding:"required"`
	State     string  `json:"state" binding:"required"`
}

type FarmResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	OwnerName string  `json:"owner_name"`
	Location  string  `json:"location"`
	TotalArea float64 `json:"total_area"`
	City      string  `json:"city"`
	State     string  `json:"state"`
}
