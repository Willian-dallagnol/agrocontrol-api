package entities

import "time"

// 🌱 Estrutura que representa a tabela "fields" (talhões) no banco
type Field struct {
	ID uint `gorm:"primaryKey"`
	// 👉 identificador único do talhão

	Name string `gorm:"not null"`
	// 👉 nome do talhão (ex: Talhão 1, Área Norte)

	Area float64 `gorm:"not null"`
	// 👉 área do talhão
	// 🔥 usada em regra de negócio (não pode ser <= 0)

	SoilType string
	// 👉 tipo de solo (ex: argiloso, arenoso)
	// (campo opcional)

	FarmID uint `gorm:"not null"`
	// 👉 referência à fazenda
	// 🔗 relacionamento:
	// cada Field pertence a uma Farm

	CreatedAt time.Time
	// 👉 data de criação do registro

	UpdatedAt time.Time
	// 👉 data da última atualização
}
