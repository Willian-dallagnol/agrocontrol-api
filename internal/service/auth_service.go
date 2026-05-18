package service

import (
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/utils"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo         ports.UserRepository
	JWTSecret        string
	JWTExpHours      int
	RefreshExpHours  int
}

func NewAuthService(userRepo ports.UserRepository, secret string, expHours int) *AuthService {
	return &AuthService{
		UserRepo:        userRepo,
		JWTSecret:       secret,
		JWTExpHours:     expHours,
		RefreshExpHours: 168, // 7 dias
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email ou senha inválidos")
		}
		slog.Error("auth: erro ao buscar usuário por email", "error", err)
		return nil, errors.New("erro interno — tente novamente")
	}

	if !user.Active {
		return nil, errors.New("usuário inativo — contate o administrador")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("email ou senha inválidos")
	}

	// Access token: curta duração (configurável, padrão 15min em produção)
	accessToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), s.JWTSecret, s.JWTExpHours)
	if err != nil {
		slog.Error("auth: falha ao gerar access token", "user_id", user.ID, "error", err)
		return nil, errors.New("erro ao gerar token — tente novamente")
	}

	// Refresh token: longa duração (7 dias)
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, string(user.Role), s.JWTSecret, s.RefreshExpHours)
	if err != nil {
		slog.Error("auth: falha ao gerar refresh token", "user_id", user.ID, "error", err)
		return nil, errors.New("erro ao gerar refresh token — tente novamente")
	}

	slog.Info("auth: login bem-sucedido", "user_id", user.ID, "role", user.Role)

	return &dto.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.JWTExpHours * 3600,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.JWTSecret)
	if err != nil {
		return nil, errors.New("refresh token inválido ou expirado")
	}

	// Verifica se é realmente um refresh token
	if !claims.IsRefresh {
		return nil, errors.New("token fornecido não é um refresh token")
	}

	user, err := s.UserRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	if !user.Active {
		return nil, errors.New("usuário inativo")
	}

	// Gera novo access token
	newAccessToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), s.JWTSecret, s.JWTExpHours)
	if err != nil {
		return nil, errors.New("erro ao renovar token")
	}

	slog.Info("auth: token renovado", "user_id", user.ID)

	return &dto.RefreshTokenResponse{
		Token:     newAccessToken,
		ExpiresIn: s.JWTExpHours * 3600,
	}, nil
}
