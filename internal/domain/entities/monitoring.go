package entities

import "time"

type MonitoringType string

const (
	MonTypeInsect  MonitoringType = "insect"  // praga inseto
	MonTypeDisease MonitoringType = "disease" // doença
	MonTypeWeed    MonitoringType = "weed"    // planta daninha
	MonTypeNutri   MonitoringType = "nutritional" // deficiência nutricional
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityModerate Severity = "moderate"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type Monitoring struct {
	ID                 uint           `gorm:"primaryKey"`
	FieldID            uint           `gorm:"not null;index"`
	PlantingID         *uint          `gorm:"index"`
	InspectionDate     time.Time      `gorm:"not null;index"`
	Type               MonitoringType `gorm:"not null;index"`
	ProblemName        string         `gorm:"not null"` // nome da praga/doença
	InfestationLevel   float64        // nível de infestação (%)
	Severity           Severity       `gorm:"not null;default:'low';index"`
	CropStage          string         // estágio da cultura (V2, R3, etc)
	TechnicalRec       string         // recomendação técnica
	Urgent             bool           `gorm:"default:false"`
	Inspector          string         // quem fez a vistoria
	Notes              string
	CreatedBy          uint `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Field Field `gorm:"foreignKey:FieldID"`
}
