package service

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/utils"
	"errors"

	"gorm.io/gorm"
)

// 🔐 Service responsável pela lógica de autenticação
type AuthService struct {
	UserRepo *repository.UserRepository
	// 👉 usado para buscar usuário no banco

	JWTSecret string
	// 👉 chave usada para gerar e validar o token JWT
}

// 🏗️ Construtor do service
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

// 🔑 Função de login
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {

	// 🔍 busca usuário pelo email
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {

		// ❌ se não encontrar, retorna erro genérico (segurança)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email ou senha inválidos")
		}

		// ❌ erro inesperado
		return nil, err
	}

	// 🔐 valida senha comparando com hash (bcrypt)
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		// ❌ senha inválida
		return nil, errors.New("email ou senha inválidos")
	}

	// 🎫 gera token JWT com dados do usuário
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.JWTSecret)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta de login
	return &dto.LoginResponse{
		Token: token, // 👉 token para autenticação

		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
			// 👉 retorna dados seguros (sem senha)
		},
	}, nil
}
