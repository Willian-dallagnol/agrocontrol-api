package tests

import (
	"testing"
	"time"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"agrocontrol-api/internal/service"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB cria banco SQLite em memória com todas as entidades
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("setupTestDB: falha ao abrir banco: %v", err)
	}
	if err := db.AutoMigrate(
		&entities.User{}, &entities.Farm{}, &entities.Field{},
		&entities.Crop{}, &entities.Season{}, &entities.Planting{},
		&entities.Input{}, &entities.Application{}, &entities.Monitoring{},
		&entities.Harvest{}, &entities.Alert{},
	); err != nil {
		t.Fatalf("setupTestDB: falha no migrate: %v", err)
	}
	return db
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func mustCreateFarm(t *testing.T, db *gorm.DB, userID uint) *entities.Farm {
	t.Helper()
	farm := &entities.Farm{Name: "Fazenda Teste", OwnerName: "Dono", TotalArea: 100, City: "Londrina", State: "PR", CreatedBy: userID}
	if err := db.Create(farm).Error; err != nil {
		t.Fatalf("mustCreateFarm: %v", err)
	}
	return farm
}

func mustCreateField(t *testing.T, db *gorm.DB, farmID, userID uint, status entities.FieldStatus) *entities.Field {
	t.Helper()
	field := &entities.Field{Name: "Talhão 01", Area: 50, Status: status, FarmID: farmID, CreatedBy: userID}
	if err := db.Create(field).Error; err != nil {
		t.Fatalf("mustCreateField: %v", err)
	}
	return field
}

func mustCreateSeason(t *testing.T, db *gorm.DB) *entities.Season {
	t.Helper()
	season := &entities.Season{
		Name:      "Safra 24/25",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 6, 0),
		Status:    entities.SeasonStatusActive,
		CreatedBy: 1,
	}
	if err := db.Create(season).Error; err != nil {
		t.Fatalf("mustCreateSeason: %v", err)
	}
	return season
}

func mustCreateCrop(t *testing.T, db *gorm.DB) *entities.Crop {
	t.Helper()
	crop := &entities.Crop{Name: "Soja", CreatedBy: 1}
	if err := db.Create(crop).Error; err != nil {
		t.Fatalf("mustCreateCrop: %v", err)
	}
	return crop
}

func mustCreateInput(t *testing.T, db *gorm.DB, userID uint, stock, minStock float64) *entities.Input {
	t.Helper()
	input := &entities.Input{
		Name: "Herbicida X", Category: "herbicide", Unit: "L",
		StockQty: stock, MinStockQty: minStock, CostPerUnit: 50,
		Active: true, CreatedBy: userID,
	}
	if err := db.Create(input).Error; err != nil {
		t.Fatalf("mustCreateInput: %v", err)
	}
	return input
}

// ── FarmService ───────────────────────────────────────────────────────────────

func TestFarmService_Create_Success(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewFarmService(repository.NewFarmRepository(db))

	farm, err := svc.CreateFarm(dto.CreateFarmRequest{
		Name: "Fazenda A", OwnerName: "João", TotalArea: 100, City: "Londrina", State: "PR",
	}, 1)

	if err != nil {
		t.Fatalf("esperava sucesso, got err: %v", err)
	}
	if farm.ID == 0 {
		t.Error("esperava ID > 0")
	}
	if farm.Name != "Fazenda A" {
		t.Errorf("nome incorreto: got %q", farm.Name)
	}
}

func TestFarmService_Create_InvalidArea(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewFarmService(repository.NewFarmRepository(db))

	_, err := svc.CreateFarm(dto.CreateFarmRequest{Name: "X", TotalArea: -1}, 1)
	if err == nil {
		t.Error("esperava erro para área negativa")
	}
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("esperava ErrInvalidInput, got: %v", err)
	}
}

func TestFarmService_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewFarmService(repository.NewFarmRepository(db))

	_, err := svc.GetFarmByID(9999, 1, "operator")
	if err == nil {
		t.Error("esperava erro not found")
	}
	if !apperrors.IsNotFound(err) {
		t.Errorf("esperava ErrNotFound, got: %v", err)
	}
}

func TestFarmService_GetByID_Forbidden(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewFarmService(repository.NewFarmRepository(db))

	farm := mustCreateFarm(t, db, 1)

	// userID=2 tentando acessar farm do userID=1
	_, err := svc.GetFarmByID(farm.ID, 2, "operator")
	if err == nil {
		t.Error("esperava erro forbidden")
	}
	if !apperrors.IsForbidden(err) {
		t.Errorf("esperava ErrForbidden, got: %v", err)
	}
}

func TestFarmService_GetByID_AdminCanSeeAll(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewFarmService(repository.NewFarmRepository(db))

	farm := mustCreateFarm(t, db, 1)

	// admin pode ver fazenda de qualquer usuário
	result, err := svc.GetFarmByID(farm.ID, 99, "admin")
	if err != nil {
		t.Fatalf("admin deveria ver a fazenda: %v", err)
	}
	if result.ID != farm.ID {
		t.Error("ID inconsistente")
	}
}

// ── FieldService ──────────────────────────────────────────────────────────────

func TestFieldService_Create_DuplicateName(t *testing.T) {
	db := setupTestDB(t)
	farmRepo  := repository.NewFarmRepository(db)
	fieldRepo := repository.NewFieldRepository(db)
	svc := service.NewFieldService(fieldRepo, farmRepo)

	farm := mustCreateFarm(t, db, 1)

	req := dto.CreateFieldRequest{Name: "Talhão A", Area: 10, FarmID: farm.ID}
	if _, err := svc.CreateField(req, 1, "manager"); err != nil {
		t.Fatalf("primeira criação deveria funcionar: %v", err)
	}

	_, err := svc.CreateField(req, 1, "manager")
	if err == nil {
		t.Error("esperava erro de conflito para nome duplicado")
	}
	if !apperrors.IsConflict(err) {
		t.Errorf("esperava ErrConflict, got: %v", err)
	}
}

// ── InputService ──────────────────────────────────────────────────────────────

func TestInputService_Create_NegativeStock(t *testing.T) {
	db := setupTestDB(t)
	txRunner  := repository.NewGormTxRunner(db)
	inputRepo := repository.NewInputRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	svc := service.NewInputService(inputRepo, alertRepo, txRunner)

	_, err := svc.CreateInput(dto.CreateInputRequest{
		Name: "X", StockQty: -10, Unit: "L",
	}, 1)
	if err == nil {
		t.Error("esperava erro para estoque negativo")
	}
}

func TestInputService_AdjustStock_WouldGoNegative(t *testing.T) {
	db := setupTestDB(t)
	txRunner  := repository.NewGormTxRunner(db)
	inputRepo := repository.NewInputRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	svc := service.NewInputService(inputRepo, alertRepo, txRunner)

	input := mustCreateInput(t, db, 1, 10, 5)

	// Tenta debitar mais do que o estoque
	_, err := svc.AdjustStock(input.ID, dto.AdjustStockRequest{Quantity: -100}, 1, "admin")
	if err == nil {
		t.Error("esperava erro para estoque insuficiente")
	}
	if !apperrors.IsInsufficientStock(err) {
		t.Errorf("esperava ErrInsufficientStock, got: %v", err)
	}
}

// ── PlantingService ───────────────────────────────────────────────────────────

func TestPlantingService_Create_InactiveField(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewPlantingService(
		repository.NewPlantingRepository(db),
		repository.NewFieldRepository(db),
		repository.NewSeasonRepository(db),
		repository.NewCropRepository(db),
	)

	farm   := mustCreateFarm(t, db, 1)
	field  := mustCreateField(t, db, farm.ID, 1, entities.FieldStatusInactive)
	season := mustCreateSeason(t, db)
	crop   := mustCreateCrop(t, db)

	_, err := svc.CreatePlanting(dto.CreatePlantingRequest{
		FieldID: field.ID, SeasonID: season.ID, CropID: crop.ID,
		PlantingDate: time.Now(),
	}, 1, "manager")

	if err == nil {
		t.Error("esperava erro para talhão inativo")
	}
}

// ── AuthService ───────────────────────────────────────────────────────────────

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewAuthService(repository.NewUserRepository(db), "secret-teste-32-caracteres-aqui!", 1)

	_, err := svc.Login(dto.LoginRequest{Email: "nao@existe.com", Password: "qualquer"})
	if err == nil {
		t.Error("esperava erro para credenciais inválidas")
	}
}

func TestAuthService_RefreshToken_Invalid(t *testing.T) {
	db := setupTestDB(t)
	svc := service.NewAuthService(repository.NewUserRepository(db), "secret-teste-32-caracteres-aqui!", 1)

	_, err := svc.RefreshToken(dto.RefreshTokenRequest{RefreshToken: "token-invalido"})
	if err == nil {
		t.Error("esperava erro para refresh token inválido")
	}
}
