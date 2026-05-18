package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreateInputRequest struct {
	Name           string                `json:"name"            binding:"required,min=2,max=150"`
	Category       entities.InputCategory `json:"category"        binding:"required,oneof=fertilizer herbicide fungicide insecticide correctant biological seed other"`
	Manufacturer   string                `json:"manufacturer"`
	BatchNumber    string                `json:"batch_number"`
	ExpirationDate *time.Time            `json:"expiration_date"`
	Unit           string                `json:"unit"            binding:"required"`
	StockQty       float64               `json:"stock_qty"       binding:"omitempty,gte=0"`
	MinStockQty    float64               `json:"min_stock_qty"   binding:"omitempty,gte=0"`
	CostPerUnit    float64               `json:"cost_per_unit"   binding:"omitempty,gte=0"`
}

type UpdateInputRequest struct {
	Name           string                `json:"name"            binding:"required,min=2,max=150"`
	Manufacturer   string                `json:"manufacturer"`
	BatchNumber    string                `json:"batch_number"`
	ExpirationDate *time.Time            `json:"expiration_date"`
	Unit           string                `json:"unit"            binding:"required"`
	MinStockQty    float64               `json:"min_stock_qty"   binding:"omitempty,gte=0"`
	CostPerUnit    float64               `json:"cost_per_unit"   binding:"omitempty,gte=0"`
	Active         bool                  `json:"active"`
}

type AdjustStockRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Reason   string  `json:"reason"   binding:"required"`
}

type InputResponse struct {
	ID             uint                  `json:"id"`
	Name           string                `json:"name"`
	Category       entities.InputCategory `json:"category"`
	Manufacturer   string                `json:"manufacturer"`
	BatchNumber    string                `json:"batch_number"`
	ExpirationDate *time.Time            `json:"expiration_date"`
	Unit           string                `json:"unit"`
	StockQty       float64               `json:"stock_qty"`
	MinStockQty    float64               `json:"min_stock_qty"`
	CostPerUnit    float64               `json:"cost_per_unit"`
	Active         bool                  `json:"active"`
	LowStock       bool                  `json:"low_stock"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}
