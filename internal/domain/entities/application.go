package entities

import "time"

type ApplicationType string

const (
	AppTypeFertilizer  ApplicationType = "fertilizer"
	AppTypeHerbicide   ApplicationType = "herbicide"
	AppTypeFungicide   ApplicationType = "fungicide"
	AppTypeInsecticide ApplicationType = "insecticide"
	AppTypeCorrectant  ApplicationType = "correctant"
	AppTypeBiological  ApplicationType = "biological"
)

type Application struct {
	ID              uint            `gorm:"primaryKey"`
	FieldID         uint            `gorm:"not null;index"`
	PlantingID      *uint           `gorm:"index"` // opcional: vincula ao plantio
	InputID         uint            `gorm:"not null;index"`
	ApplicationType ApplicationType `gorm:"not null;index"`
	ApplicationDate time.Time       `gorm:"not null;index"`
	DosePerHa       float64         `gorm:"not null"` // dose por hectare
	TotalUsed       float64         `gorm:"not null"` // total aplicado = dose * área
	SprayVolume     float64         // volume de calda (L/ha)
	Target          string          // alvo (ex: lagarta, ferrugem)
	Equipment       string          // equipamento usado
	Operator        string          // operador responsável
	WindSpeed       float64         // velocidade do vento (km/h)
	Temperature     float64         // temperatura (°C)
	Humidity        float64         // umidade relativa (%)
	Notes           string
	CreatedBy       uint `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Field    Field    `gorm:"foreignKey:FieldID"`
	Input    Input    `gorm:"foreignKey:InputID"`
}
