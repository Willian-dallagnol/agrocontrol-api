package dto

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

type CreateMonitoringRequest struct {
	FieldID          uint                   `json:"field_id"           binding:"required"`
	PlantingID       *uint                  `json:"planting_id"`
	InspectionDate   time.Time              `json:"inspection_date"    binding:"required"`
	Type             entities.MonitoringType `json:"type"               binding:"required,oneof=insect disease weed nutritional"`
	ProblemName      string                 `json:"problem_name"       binding:"required,min=2"`
	InfestationLevel float64                `json:"infestation_level"  binding:"omitempty,gte=0,lte=100"`
	Severity         entities.Severity      `json:"severity"           binding:"required,oneof=low moderate high critical"`
	CropStage        string                 `json:"crop_stage"`
	TechnicalRec     string                 `json:"technical_rec"`
	Urgent           bool                   `json:"urgent"`
	Inspector        string                 `json:"inspector"          binding:"required"`
	Notes            string                 `json:"notes"`
}

type MonitoringResponse struct {
	ID               uint                   `json:"id"`
	FieldID          uint                   `json:"field_id"`
	FieldName        string                 `json:"field_name,omitempty"`
	PlantingID       *uint                  `json:"planting_id"`
	InspectionDate   time.Time              `json:"inspection_date"`
	Type             entities.MonitoringType `json:"type"`
	ProblemName      string                 `json:"problem_name"`
	InfestationLevel float64                `json:"infestation_level"`
	Severity         entities.Severity      `json:"severity"`
	CropStage        string                 `json:"crop_stage"`
	TechnicalRec     string                 `json:"technical_rec"`
	Urgent           bool                   `json:"urgent"`
	Inspector        string                 `json:"inspector"`
	Notes            string                 `json:"notes"`
	CreatedBy        uint                   `json:"created_by"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}
