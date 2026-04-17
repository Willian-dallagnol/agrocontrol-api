package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

// 👤 Repository responsável por acessar a tabela de User (usuários)
type UserRepository struct {
	DB *gorm.DB
	// 👉 conexão com o banco de dados
}

// 🏗️ Construtor do repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// 🚀 Criar novo usuário
func (r *UserRepository) Create(user *entities.User) error {
	// 👉 insere um novo registro na tabela users
	return r.DB.Create(user).Error
}

// 🔍 Buscar usuário pelo email
func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User

	// 👉 SELECT * FROM users WHERE email = ?
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		// ❌ retorna erro se não encontrar
		return nil, err
	}

	return &user, nil
}
