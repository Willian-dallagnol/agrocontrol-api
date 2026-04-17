package entities

import "time"

// 🚜 Estrutura que representa a tabela "farms" no banco de dados
type Farm struct {
	ID uint `gorm:"primaryKey"`
	// 👉 identificador único da fazenda (chave primária)

	Name string `gorm:"not null"`
	// 👉 nome da fazenda (obrigatório)

	OwnerName string `gorm:"not null"`
	// 👉 nome do proprietário (obrigatório)

	Location string
	// 👉 localização geral (ex: região, estrada, etc)

	TotalArea float64 `gorm:"not null"`
	// 👉 área total da fazenda (obrigatório)
	// 🔥 usado em regras de negócio (ex: não pode ser <= 0)

	City string `gorm:"not null"`
	// 👉 cidade da fazenda

	State string `gorm:"not null"`
	// 👉 estado da fazenda

	CreatedBy uint
	// 👉 ID do usuário que criou a fazenda
	// 🔗 relacionamento indireto com User

	CreatedAt time.Time
	// 👉 data de criação do registro

	UpdatedAt time.Time
	// 👉 data da última atualização
}
