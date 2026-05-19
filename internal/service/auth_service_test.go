package service

import (
	"errors"
	"testing"

	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/utils"
)

// ── Mock UserRepository ────────────────────────────────────────────────────

type mockUserRepo struct {
	user *entities.User
	err  error
}

func (m *mockUserRepo) Create(u *entities.User) error            { return m.err }
func (m *mockUserRepo) FindByID(id uint) (*entities.User, error) { return m.user, m.err }
func (m *mockUserRepo) FindByEmail(email string) (*entities.User, error) {
	return m.user, m.err
}
func (m *mockUserRepo) FindAll() ([]entities.User, error) { return nil, m.err }
func (m *mockUserRepo) Update(u *entities.User) error     { return m.err }

// ── Helpers ────────────────────────────────────────────────────────────────

const testJWTSecret = "test_secret_key_with_32_chars_ok!"

func newAuthSvc(repo *mockUserRepo) *AuthService {
	return NewAuthService(repo, testJWTSecret, 1)
}

func hashedPassword(plain string) string {
	h, _ := utils.HashPassword(plain)
	return h
}

// ── Auth tests ─────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	user := &entities.User{
		ID:           1,
		Email:        "willian@teste.com",
		PasswordHash: hashedPassword("senha123"),
		Role:         entities.RoleManager,
		Active:       true,
	}
	svc := newAuthSvc(&mockUserRepo{user: user})

	resp, err := svc.Login(dto.LoginRequest{
		Email:    "willian@teste.com",
		Password: "senha123",
	})

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Token == "" {
		t.Error("esperava token não vazio")
	}
	if resp.RefreshToken == "" {
		t.Error("esperava refresh token não vazio")
	}
	if resp.User.Email != "willian@teste.com" {
		t.Errorf("esperava email correto, got '%s'", resp.User.Email)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc := newAuthSvc(&mockUserRepo{err: errors.New("not found")})

	_, err := svc.Login(dto.LoginRequest{
		Email:    "naoexiste@teste.com",
		Password: "senha123",
	})

	if err == nil {
		t.Fatal("esperava erro para usuário não encontrado")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	user := &entities.User{
		ID:           1,
		Email:        "willian@teste.com",
		PasswordHash: hashedPassword("senha_correta"),
		Role:         entities.RoleManager,
		Active:       true,
	}
	svc := newAuthSvc(&mockUserRepo{user: user})

	_, err := svc.Login(dto.LoginRequest{
		Email:    "willian@teste.com",
		Password: "senha_errada",
	})

	if err == nil {
		t.Fatal("esperava erro para senha incorreta")
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	user := &entities.User{
		ID:           1,
		Email:        "inativo@teste.com",
		PasswordHash: hashedPassword("senha123"),
		Role:         entities.RoleOperator,
		Active:       false,
	}
	svc := newAuthSvc(&mockUserRepo{user: user})

	_, err := svc.Login(dto.LoginRequest{
		Email:    "inativo@teste.com",
		Password: "senha123",
	})

	if err == nil {
		t.Fatal("esperava erro para usuário inativo")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	user := &entities.User{
		ID:     1,
		Email:  "willian@teste.com",
		Role:   entities.RoleManager,
		Active: true,
	}
	svc := newAuthSvc(&mockUserRepo{user: user})

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, string(user.Role), testJWTSecret, 168)
	if err != nil {
		t.Fatalf("erro ao gerar refresh token: %v", err)
	}

	resp, err := svc.RefreshToken(dto.RefreshTokenRequest{RefreshToken: refreshToken})

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Token == "" {
		t.Error("esperava novo token não vazio")
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	svc := newAuthSvc(&mockUserRepo{})

	_, err := svc.RefreshToken(dto.RefreshTokenRequest{RefreshToken: "token_invalido"})

	if err == nil {
		t.Fatal("esperava erro para token inválido")
	}
}

func TestRefreshToken_UserInactive(t *testing.T) {
	user := &entities.User{
		ID:     1,
		Email:  "inativo@teste.com",
		Role:   entities.RoleOperator,
		Active: false,
	}
	svc := newAuthSvc(&mockUserRepo{user: user})

	refreshToken, _ := utils.GenerateRefreshToken(user.ID, user.Email, string(user.Role), testJWTSecret, 168)
	_, err := svc.RefreshToken(dto.RefreshTokenRequest{RefreshToken: refreshToken})

	if err == nil {
		t.Fatal("esperava erro para usuário inativo no refresh")
	}
}
