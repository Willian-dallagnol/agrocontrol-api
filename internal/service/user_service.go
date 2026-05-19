package service

import (
	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/utils"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

type UserService struct {
	Repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	role := entities.Role(strings.ToLower(req.Role))
	if _, ok := entities.ValidRoles[role]; !ok {
		return nil, fmt.Errorf("role inválida '%s': use admin, manager ou operator: %w",
			req.Role, apperrors.ErrInvalidInput)
	}
	if !isValidEmail(strings.TrimSpace(req.Email)) {
		return nil, fmt.Errorf("email inválido: %w", apperrors.ErrInvalidInput)
	}
	if len(strings.TrimSpace(req.Password)) < 8 {
		return nil, fmt.Errorf("senha deve ter no mínimo 8 caracteres: %w", apperrors.ErrInvalidInput)
	}

	existing, err := s.Repo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("user: erro ao verificar email", "email", req.Email, "error", err)
		return nil, errors.New("erro interno ao verificar email")
	}
	if existing != nil {
		return nil, apperrors.ConflictError("email")
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		slog.Error("user: erro ao hashear senha", "error", err)
		return nil, errors.New("erro ao processar senha")
	}
	user := &entities.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
		Role:         role,
		Active:       true,
	}
	if err := s.Repo.Create(user); err != nil {
		slog.Error("user: erro ao criar usuário", "error", err)
		return nil, errors.New("erro ao criar usuário")
	}
	slog.Info("user: criado", "user_id", user.ID, "role", user.Role)
	return toUserResponse(user), nil
}

func (s *UserService) GetUsers() ([]dto.UserResponse, error) {
	users, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}
	resp := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		resp = append(resp, *toUserResponse(&users[i]))
	}
	return resp, nil
}

func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func toUserResponse(u *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID: u.ID, Name: u.Name, Email: u.Email,
		Role: string(u.Role), Active: u.Active, CreatedAt: u.CreatedAt,
	}
}

func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := strings.LastIndex(email, "@")
	if at < 1 {
		return false
	}
	domain := email[at+1:]
	if len(domain) < 3 || !strings.Contains(domain, ".") {
		return false
	}
	return true
}
