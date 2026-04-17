package entities

import "time"

// 👤 Estrutura que representa a tabela "users" no banco de dados
type User struct {
	ID uint `gorm:"primaryKey"`
	// 👉 identificador único do usuário

	Name string `gorm:"not null"`
	// 👉 nome do usuário (obrigatório)

	Email string `gorm:"unique;not null"`
	// 👉 email do usuário (único no sistema)
	// 🔥 evita duplicidade de contas

	PasswordHash string `gorm:"not null"`
	// 👉 senha criptografada (bcrypt)
	// 🔒 nunca armazenar senha em texto puro

	Role string `gorm:"not null"` // admin, manager, operator
	// 👉 nível de acesso do usuário
	// 🔥 usado para controle de permissões

	CreatedAt time.Time
	// 👉 data de criação do usuário

	UpdatedAt time.Time
	// 👉 data da última atualização
}
