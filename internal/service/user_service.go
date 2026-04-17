package service

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/utils"
	"errors"

	"gorm.io/gorm"
)

// 👤 Service responsável pela regra de negócio de usuários
type UserService struct {
	Repo *repository.UserRepository
	// 👉 acesso ao banco (users)
}

// 🏗️ Construtor do service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// 🚀 Criar novo usuário
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {

	// 🔐 valida se a role é permitida
	if req.Role != "admin" && req.Role != "manager" && req.Role != "operator" {
		return nil, errors.New("role inválida")
	}

	// 🔍 verifica se já existe usuário com esse email
	existingUser, err := s.Repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		// ❌ email já cadastrado
		return nil, errors.New("email já cadastrado")
	}

	// ⚠️ trata erro inesperado (não relacionado a "não encontrado")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 🔐 gera hash da senha (bcrypt)
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 🧩 monta entidade User
	user := &entities.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword, // 🔒 nunca salvar senha pura
		Role:         req.Role,
	}

	// 💾 salva no banco
	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta segura (sem senha)
	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
