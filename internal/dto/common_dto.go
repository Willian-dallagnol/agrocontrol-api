package dto

import "math"

// ── Paginação base ────────────────────────────────────────────────────────────

type PaginationQuery struct {
	Page   int    `form:"page,default=1"   binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

// ── Query com filtros por módulo ──────────────────────────────────────────────

type FarmQuery struct {
	PaginationQuery
}

type FieldQuery struct {
	PaginationQuery
}

type CropQuery struct {
	PaginationQuery
}

type SeasonQuery struct {
	PaginationQuery
}

type InputQuery struct {
	PaginationQuery
	Category string `form:"category"`
}

type ApplicationQuery struct {
	PaginationQuery
	FieldID uint   `form:"field_id"`
	Type    string `form:"type"`
}

type HarvestQuery struct {
	PaginationQuery
	FieldID uint `form:"field_id"`
}

// ── Resposta paginada genérica ────────────────────────────────────────────────

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func NewPaginatedResponse(data interface{}, total int64, page, limit int) PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// ── Dashboard ────────────────────────────────────────────────────────────────

type DashboardResponse struct {
	TotalFarms            int64           `json:"total_farms"`
	TotalFields           int64           `json:"total_fields"`
	ActivePlantings       int64           `json:"active_plantings"`
	OpenAlerts            int64           `json:"open_alerts"`
	LowStockInputs        int64           `json:"low_stock_inputs"`
	ApplicationsThisMonth int64           `json:"applications_this_month"`
	LastAlerts            []AlertResponse `json:"last_alerts"`
}
