package tests

import (
	"errors"
	"testing"
	"time"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/mocks"
	"agrocontrol-api/internal/service"

	"gorm.io/gorm"
)

// ── FarmService com mocks ─────────────────────────────────────────────────────

func TestFarmService_Mock_Create_InvalidArea(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{}
	svc := service.NewFarmService(repo)

	_, err := svc.CreateFarm(dto.CreateFarmRequest{Name: "X", TotalArea: 0}, 1)
	if err == nil {
		t.Fatal("esperava erro para área zero")
	}
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput, got: %v", err)
	}
}

func TestFarmService_Mock_Create_Success(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		CreateFn: func(farm *entities.Farm) error {
			farm.ID = 42 // simula auto-increment do banco
			return nil
		},
	}
	svc := service.NewFarmService(repo)

	resp, err := svc.CreateFarm(dto.CreateFarmRequest{
		Name: "Fazenda Mock", OwnerName: "Dono", TotalArea: 100, City: "Londrina", State: "PR",
	}, 1)
	if err != nil {
		t.Fatalf("esperava sucesso, got: %v", err)
	}
	if resp.ID != 42 {
		t.Errorf("esperava ID=42, got: %d", resp.ID)
	}
}

func TestFarmService_Mock_GetByID_NotFound(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewFarmService(repo)

	_, err := svc.GetFarmByID(99, 1, "operator")
	if !apperrors.IsNotFound(err) {
		t.Errorf("esperava ErrNotFound, got: %v", err)
	}
}

func TestFarmService_Mock_GetByID_Forbidden(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil // pertence ao user 1
		},
	}
	svc := service.NewFarmService(repo)

	// user 2 tentando acessar fazenda do user 1
	_, err := svc.GetFarmByID(1, 2, "operator")
	if !apperrors.IsForbidden(err) {
		t.Errorf("esperava ErrForbidden, got: %v", err)
	}
}

func TestFarmService_Mock_GetByID_AdminBypass(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil
		},
	}
	svc := service.NewFarmService(repo)

	resp, err := svc.GetFarmByID(1, 99, "admin") // admin vê qualquer fazenda
	if err != nil {
		t.Fatalf("admin deveria acessar: %v", err)
	}
	if resp.ID != 1 {
		t.Error("ID incorreto")
	}
}

func TestFarmService_Mock_Delete_Forbidden(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil
		},
	}
	svc := service.NewFarmService(repo)

	err := svc.DeleteFarm(1, 2, "manager") // user 2 tentando deletar farm do user 1
	if !apperrors.IsForbidden(err) {
		t.Errorf("esperava ErrForbidden, got: %v", err)
	}
}

func TestFarmService_Mock_Update_InvalidArea(t *testing.T) {
	repo := &mocks.FarmRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Farm, error) {
			return &entities.Farm{ID: id, CreatedBy: 1}, nil
		},
	}
	svc := service.NewFarmService(repo)

	_, err := svc.UpdateFarm(1, dto.UpdateFarmRequest{Name: "X", TotalArea: -5}, 1, "manager")
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput, got: %v", err)
	}
}

// ── InputService com mocks ────────────────────────────────────────────────────

func TestInputService_Mock_Create_NegativeStock(t *testing.T) {
	repo := &mocks.InputRepositoryMock{}
	alertRepo := &mocks.AlertRepositoryMock{}
	tx := &mocks.TxRunnerMock{}
	svc := service.NewInputService(repo, alertRepo, tx)

	_, err := svc.CreateInput(dto.CreateInputRequest{Name: "X", StockQty: -1, Unit: "L"}, 1)
	if err == nil {
		t.Fatal("esperava erro para estoque negativo")
	}
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput, got: %v", err)
	}
}

func TestInputService_Mock_Create_NegativeCost(t *testing.T) {
	repo := &mocks.InputRepositoryMock{}
	alertRepo := &mocks.AlertRepositoryMock{}
	tx := &mocks.TxRunnerMock{}
	svc := service.NewInputService(repo, alertRepo, tx)

	_, err := svc.CreateInput(dto.CreateInputRequest{Name: "X", CostPerUnit: -10, Unit: "L"}, 1)
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput, got: %v", err)
	}
}

func TestInputService_Mock_Create_Success(t *testing.T) {
	repo := &mocks.InputRepositoryMock{
		CreateFn: func(i *entities.Input) error {
			i.ID = 1
			return nil
		},
	}
	alertRepo := &mocks.AlertRepositoryMock{}
	tx := &mocks.TxRunnerMock{}
	svc := service.NewInputService(repo, alertRepo, tx)

	resp, err := svc.CreateInput(dto.CreateInputRequest{Name: "Herbicida X", Unit: "L", StockQty: 100}, 1)
	if err != nil {
		t.Fatalf("esperava sucesso: %v", err)
	}
	if resp.Name != "Herbicida X" {
		t.Errorf("nome incorreto: %q", resp.Name)
	}
}

func TestInputService_Mock_GetByID_NotFound(t *testing.T) {
	repo := &mocks.InputRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Input, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	alertRepo := &mocks.AlertRepositoryMock{}
	tx := &mocks.TxRunnerMock{}
	svc := service.NewInputService(repo, alertRepo, tx)

	_, err := svc.GetInputByID(99, 1, "operator")
	if !apperrors.IsNotFound(err) {
		t.Errorf("esperava ErrNotFound, got: %v", err)
	}
}

func TestInputService_Mock_GetByID_Forbidden(t *testing.T) {
	repo := &mocks.InputRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Input, error) {
			return &entities.Input{CreatedBy: 1}, nil
		},
	}
	alertRepo := &mocks.AlertRepositoryMock{}
	tx := &mocks.TxRunnerMock{}
	svc := service.NewInputService(repo, alertRepo, tx)

	_, err := svc.GetInputByID(1, 2, "operator")
	if !apperrors.IsForbidden(err) {
		t.Errorf("esperava ErrForbidden, got: %v", err)
	}
}

func TestInputService_Mock_AdjustStock_GoNegative(t *testing.T) {
	repo := &mocks.InputRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Input, error) {
			return &entities.Input{ID: id, StockQty: 10, CreatedBy: 1}, nil
		},
	}
	alertRepo := &mocks.AlertRepositoryMock{}
	tx := &mocks.TxRunnerMock{}
	svc := service.NewInputService(repo, alertRepo, tx)

	_, err := svc.AdjustStock(1, dto.AdjustStockRequest{Quantity: -100}, 1, "manager")
	if !apperrors.IsInsufficientStock(err) {
		t.Errorf("esperava ErrInsufficientStock, got: %v", err)
	}
}

// ── AuthService com mocks ─────────────────────────────────────────────────────

func TestAuthService_Mock_Login_UserNotFound(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(email string) (*entities.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewAuthService(repo, "secret_32_chars_aqui_padded_ok!!", 1)

	_, err := svc.Login(dto.LoginRequest{Email: "nao@existe.com", Password: "qualquer"})
	if err == nil {
		t.Fatal("esperava erro para usuário não encontrado")
	}
}

func TestAuthService_Mock_Login_WrongPassword(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(email string) (*entities.User, error) {
			return &entities.User{
				Email:        email,
				PasswordHash: "$2a$10$invalido",
				Active:       true,
			}, nil
		},
	}
	svc := service.NewAuthService(repo, "secret_32_chars_aqui_padded_ok!!", 1)

	_, err := svc.Login(dto.LoginRequest{Email: "user@test.com", Password: "errada"})
	if err == nil {
		t.Fatal("esperava erro para senha errada")
	}
}

func TestAuthService_Mock_Login_InactiveUser(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(email string) (*entities.User, error) {
			return &entities.User{Email: email, Active: false}, nil
		},
	}
	svc := service.NewAuthService(repo, "secret_32_chars_aqui_padded_ok!!", 1)

	_, err := svc.Login(dto.LoginRequest{Email: "inativo@test.com", Password: "qualquer"})
	if err == nil {
		t.Fatal("esperava erro para usuário inativo")
	}
}

func TestAuthService_Mock_RefreshToken_Invalid(t *testing.T) {
	repo := &mocks.UserRepositoryMock{}
	svc := service.NewAuthService(repo, "secret_32_chars_aqui_padded_ok!!", 1)

	_, err := svc.RefreshToken(dto.RefreshTokenRequest{RefreshToken: "token-invalido"})
	if err == nil {
		t.Fatal("esperava erro para refresh token inválido")
	}
}

// ── AlertService com mocks ────────────────────────────────────────────────────

func TestAlertService_Mock_GetByID_NotFound(t *testing.T) {
	repo := &mocks.AlertRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Alert, error) {
			return nil, errors.New("not found")
		},
	}
	svc := service.NewAlertService(repo)

	_, err := svc.GetAlertByID(99, 1, "operator")
	if !apperrors.IsNotFound(err) {
		t.Errorf("esperava ErrNotFound, got: %v", err)
	}
}

func TestAlertService_Mock_GetByID_Forbidden(t *testing.T) {
	repo := &mocks.AlertRepositoryMock{
		FindByIDFn: func(id uint) (*entities.Alert, error) {
			return &entities.Alert{ID: id, CreatedBy: 1}, nil
		},
	}
	svc := service.NewAlertService(repo)

	_, err := svc.GetAlertByID(1, 2, "operator") // user 2 tentando ver alerta do user 1
	if !apperrors.IsForbidden(err) {
		t.Errorf("esperava ErrForbidden, got: %v", err)
	}
}

func TestAlertService_Mock_Create_DefaultPriority(t *testing.T) {
	var savedAlert *entities.Alert
	repo := &mocks.AlertRepositoryMock{
		CreateFn: func(a *entities.Alert) error {
			savedAlert = a
			return nil
		},
	}
	svc := service.NewAlertService(repo)

	_, err := svc.CreateAlert(dto.CreateAlertRequest{
		Title: "Teste", Type: entities.AlertTypeLowStock,
	}, 1)
	if err != nil {
		t.Fatalf("esperava sucesso: %v", err)
	}
	if savedAlert.Priority != entities.AlertPriorityMedium {
		t.Errorf("esperava prioridade medium por padrão, got: %q", savedAlert.Priority)
	}
}

// ── UserService com mocks ─────────────────────────────────────────────────────

func TestUserService_Mock_Create_InvalidRole(t *testing.T) {
	repo := &mocks.UserRepositoryMock{}
	svc := service.NewUserService(repo)

	_, err := svc.CreateUser(dto.CreateUserRequest{
		Name: "X", Email: "x@x.com", Password: "senha123456", Role: "superadmin",
	})
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput para role inválida, got: %v", err)
	}
}

func TestUserService_Mock_Create_DuplicateEmail(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(email string) (*entities.User, error) {
			return &entities.User{Email: email}, nil // email já existe
		},
	}
	svc := service.NewUserService(repo)

	_, err := svc.CreateUser(dto.CreateUserRequest{
		Name: "João", Email: "jo@jo.com", Password: "senha123456", Role: "operator",
	})
	if !apperrors.IsConflict(err) {
		t.Errorf("esperava ErrConflict para email duplicado, got: %v", err)
	}
}

// ── SeasonService com mocks ───────────────────────────────────────────────────

func TestSeasonService_Mock_Create_InvalidDates(t *testing.T) {
	repo := &mocks.SeasonRepositoryMock{}
	svc := service.NewSeasonService(repo)

	start := time.Now()
	end := start.AddDate(0, 0, -1) // end antes de start

	_, err := svc.CreateSeason(dto.CreateSeasonRequest{
		Name:      "Safra X",
		StartDate: start,
		EndDate:   end,
	}, 1)
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput para datas inválidas, got: %v", err)
	}
}
