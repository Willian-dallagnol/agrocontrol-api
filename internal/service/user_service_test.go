package service

import (
	"errors"
	"testing"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
)

func TestCreateUser_Success(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	resp, err := svc.CreateUser(dto.CreateUserRequest{
		Name:     "Willian",
		Email:    "willian@teste.com",
		Password: "senha123",
		Role:     "manager",
	})

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Willian" {
		t.Errorf("esperava 'Willian', got '%s'", resp.Name)
	}
	if resp.Role != "manager" {
		t.Errorf("esperava role 'manager', got '%s'", resp.Role)
	}
	if resp.Active != true {
		t.Error("esperava usuário ativo")
	}
}

func TestCreateUser_InvalidRole(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})

	_, err := svc.CreateUser(dto.CreateUserRequest{
		Name:     "Teste",
		Email:    "teste@teste.com",
		Password: "senha123",
		Role:     "superadmin",
	})

	if err == nil {
		t.Fatal("esperava erro para role inválida")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	existing := &entities.User{ID: 1, Email: "willian@teste.com"}
	svc := NewUserService(&mockUserRepo{user: existing})

	_, err := svc.CreateUser(dto.CreateUserRequest{
		Name:     "Outro",
		Email:    "willian@teste.com",
		Password: "senha123",
		Role:     "operator",
	})

	if err == nil {
		t.Fatal("esperava erro para email duplicado")
	}
}

func TestCreateUser_EmailLowercased(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	resp, err := svc.CreateUser(dto.CreateUserRequest{
		Name:     "Willian",
		Email:    "WILLIAN@TESTE.COM",
		Password: "senha123",
		Role:     "operator",
	})

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Email != "willian@teste.com" {
		t.Errorf("esperava email em lowercase, got '%s'", resp.Email)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	svc := NewUserService(&mockUserRepo{err: errors.New("not found")})
	_, err := svc.GetUserByID(99)

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetUserByID_Success(t *testing.T) {
	user := &entities.User{ID: 1, Name: "Willian", Email: "willian@teste.com", Role: entities.RoleManager, Active: true}
	svc := NewUserService(&mockUserRepo{user: user})

	resp, err := svc.GetUserByID(1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Willian" {
		t.Errorf("esperava 'Willian', got '%s'", resp.Name)
	}
}

func TestGetUsers_Success(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})
	result, err := svc.GetUsers()

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if result == nil {
		t.Fatal("esperava slice vazio, não nil")
	}
}
func TestCreateUser_InvalidEmail(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})

	invalidEmails := []string{
		"nao-e-email",
		"@semlocal.com",
		"semdominio@",
		"s@s",
		"",
	}

	for _, email := range invalidEmails {
		_, err := svc.CreateUser(dto.CreateUserRequest{
			Name:     "Teste",
			Email:    email,
			Password: "senha123",
			Role:     "operator",
		})
		if err == nil {
			t.Errorf("esperava erro para email inválido '%s'", email)
		}
		if !errors.Is(err, apperrors.ErrInvalidInput) {
			t.Errorf("esperava ErrInvalidInput para '%s', got %v", email, err)
		}
	}
}

func TestCreateUser_ValidEmail(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})

	_, err := svc.CreateUser(dto.CreateUserRequest{
		Name:     "Willian",
		Email:    "willian@fazenda.com.br",
		Password: "senha123",
		Role:     "manager",
	})

	if err != nil {
		t.Fatalf("esperava nil para email válido, got %v", err)
	}
}

func TestCreateUser_AllRoles(t *testing.T) {
	roles := []string{"admin", "manager", "operator"}
	for _, role := range roles {
		svc := NewUserService(&mockUserRepo{})
		_, err := svc.CreateUser(dto.CreateUserRequest{
			Name:     "Teste",
			Email:    role + "@teste.com",
			Password: "senha123",
			Role:     role,
		})
		if err != nil {
			t.Errorf("role '%s' deveria ser válida, got %v", role, err)
		}
	}

}
