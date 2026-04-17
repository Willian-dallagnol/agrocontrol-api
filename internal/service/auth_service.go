package service

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	JWTSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email ou senha inválidos")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("email ou senha inválidos")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
