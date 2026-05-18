package entities

import "time"

// InputCategory define as categorias de insumos agrícolas
type InputCategory string

const (
	InputCategoryFertilizer  InputCategory = "fertilizer"
	InputCategoryHerbicide   InputCategory = "herbicide"
	InputCategoryFungicide   InputCategory = "fungicide"
	InputCategoryInsecticide InputCategory = "insecticide"
	InputCategoryCorrectant  InputCategory = "correctant"
	InputCategoryBiological  InputCategory = "biological"
	InputCategorySeed        InputCategory = "seed"
	InputCategoryOther       InputCategory = "other"
)

// Input representa um insumo agrícola com controle de estoque
type Input struct {
	ID             uint          `gorm:"primaryKey"`
	Name           string        `gorm:"not null;index"`
	Category       InputCategory `gorm:"not null;index"`
	Manufacturer   string
	BatchNumber    string
	ExpirationDate *time.Time
	Unit           string  `gorm:"not null"`
	StockQty       float64 `gorm:"not null;default:0;check:stock_qty >= 0"`
	MinStockQty    float64 `gorm:"not null;default:0;check:min_stock_qty >= 0"`
	CostPerUnit    float64 `gorm:"not null;default:0;check:cost_per_unit >= 0"`
	Active         bool    `gorm:"not null;default:true;index"`
	CreatedBy      uint    `gorm:"not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// IsLowStock retorna true se o estoque está abaixo do mínimo
func (i *Input) IsLowStock() bool {
	return i.StockQty <= i.MinStockQty
}

// IsExpired retorna true se o insumo já venceu
func (i *Input) IsExpired() bool {
	return i.ExpirationDate != nil && i.ExpirationDate.Before(time.Now())
}

// IsExpiringSoon retorna true se vence dentro de N dias
func (i *Input) IsExpiringSoon(days int) bool {
	if i.ExpirationDate == nil {
		return false
	}
	deadline := time.Now().AddDate(0, 0, days)
	return i.ExpirationDate.Before(deadline)
}
