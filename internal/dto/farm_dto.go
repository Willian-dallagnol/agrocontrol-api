package dto

import "time"

type CreateFarmRequest struct {
	Name      string  `json:"name"       binding:"required,min=2,max=150"`
	OwnerName string  `json:"owner_name" binding:"required,min=2,max=150"`
	Location  string  `json:"location"`
	TotalArea float64 `json:"total_area" binding:"required,gt=0"`
	City      string  `json:"city"       binding:"required"`
	State     string  `json:"state"      binding:"required,len=2"`
}

type UpdateFarmRequest struct {
	Name      string  `json:"name"       binding:"required,min=2,max=150"`
	OwnerName string  `json:"owner_name" binding:"required,min=2,max=150"`
	Location  string  `json:"location"`
	TotalArea float64 `json:"total_area" binding:"required,gt=0"`
	City      string  `json:"city"       binding:"required"`
	State     string  `json:"state"      binding:"required,len=2"`
}

type FarmResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	OwnerName string    `json:"owner_name"`
	Location  string    `json:"location"`
	TotalArea float64   `json:"total_area"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
