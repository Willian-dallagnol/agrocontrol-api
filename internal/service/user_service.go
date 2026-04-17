package service

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/utils"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if req.Role != "admin" && req.Role != "manager" && req.Role != "operator" {
		return nil, errors.New("role inválida")
	}

	existingUser, err := s.Repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email já cadastrado")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
