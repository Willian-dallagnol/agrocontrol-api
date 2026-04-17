package entities

import "gorm.io/gorm"

// 🌾 Estrutura que representa a tabela "crops" no banco de dados
type Crop struct {
	gorm.Model
	// 👉 inclui automaticamente:
	// ID, CreatedAt, UpdatedAt, DeletedAt

	Name string
	// 👉 nome da cultura (ex: Soja, Milho, Trigo)

	Type string
	// 👉 tipo da cultura (ex: grão, leguminosa, cereal)

	FieldID uint
	// 👉 referência ao talhão (Field)
	// 🔗 define o relacionamento:
	// cada Crop pertence a um Field
}
