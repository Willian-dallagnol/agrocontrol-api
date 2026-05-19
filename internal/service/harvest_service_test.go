package service

import (
	"errors"
	"testing"
	"time"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/dto"
)

// ── Mock HarvestRepository ─────────────────────────────────────────────────

type mockHarvestRepo struct {
	harvest *entities.Harvest
	err     error
	exists  bool
	belongs bool
}

func (m *mockHarvestRepo) FindByID(id uint) (*entities.Harvest, error)        { return m.harvest, m.err }
func (m *mockHarvestRepo) FindAll() ([]entities.Harvest, error)               { return nil, m.err }
func (m *mockHarvestRepo) FindByUser(userID uint) ([]entities.Harvest, error) { return nil, m.err }
func (m *mockHarvestRepo) FindAllPaginated(offset, limit int, fieldID uint) ([]entities.Harvest, int64, error) {
	return nil, 0, m.err
}
func (m *mockHarvestRepo) FindByUserPaginated(userID uint, offset, limit int, fieldID uint) ([]entities.Harvest, int64, error) {
	return nil, 0, m.err
}
func (m *mockHarvestRepo) ExistsByPlantingID(plantingID uint) (bool, error) { return m.exists, m.err }
func (m *mockHarvestRepo) Update(h *entities.Harvest) error                 { return m.err }
func (m *mockHarvestRepo) Delete(id uint) error                             { return m.err }
func (m *mockHarvestRepo) CreateTx(tx ports.TxRunner, h *entities.Harvest) error {
	if m.err != nil {
		return m.err
	}
	h.ID = 1
	return nil
}

// ── Mock PlantingRepository para Harvest ──────────────────────────────────

type mockPlantingRepoForHarvest struct {
	planting *entities.Planting
	err      error
	belongs  bool
}

func (m *mockPlantingRepoForHarvest) Create(p *entities.Planting) error { return m.err }
func (m *mockPlantingRepoForHarvest) FindByID(id uint) (*entities.Planting, error) {
	return m.planting, m.err
}
func (m *mockPlantingRepoForHarvest) FindAll() ([]entities.Planting, error) { return nil, m.err }
func (m *mockPlantingRepoForHarvest) FindByUser(userID uint) ([]entities.Planting, error) {
	return nil, m.err
}
func (m *mockPlantingRepoForHarvest) FindByFieldID(fieldID uint) ([]entities.Planting, error) {
	return nil, m.err
}
func (m *mockPlantingRepoForHarvest) HasActivePlanting(fieldID uint) (bool, error) {
	return false, m.err
}
func (m *mockPlantingRepoForHarvest) CountActive() (int64, error) { return 0, m.err }
func (m *mockPlantingRepoForHarvest) BelongsToUser(plantingID, userID uint) (bool, error) {
	return m.belongs, nil
}
func (m *mockPlantingRepoForHarvest) Update(p *entities.Planting) error { return m.err }
func (m *mockPlantingRepoForHarvest) Delete(id uint) error              { return m.err }

// ── Helpers ────────────────────────────────────────────────────────────────

func validHarvestRequest() dto.CreateHarvestRequest {
	return dto.CreateHarvestRequest{
		PlantingID:        1,
		HarvestDate:       time.Now(),
		ProductivityBagHa: 65,
		ProductivityKgHa:  3900,
		TotalBags:         500,
		GrainMoisture:     13.5,
	}
}

func activePlanting() *entities.Planting {
	return &entities.Planting{ID: 1, FieldID: 1, Status: entities.PlantingStatusActive, CreatedBy: 1}
}

// ── Harvest tests ──────────────────────────────────────────────────────────

func TestCreateHarvest_Success(t *testing.T) {
	plantingRepo := &mockPlantingRepoForHarvest{planting: activePlanting(), belongs: true}
	harvestRepo := &mockHarvestRepo{harvest: &entities.Harvest{ID: 1}}
	fieldRepo := &mockFieldRepoForPlanting{field: &entities.Field{ID: 1}, belongs: true}

	svc := NewHarvestService(harvestRepo, plantingRepo, fieldRepo, &mockTxRunner{})
	resp, err := svc.CreateHarvest(validHarvestRequest(), 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp == nil {
		t.Fatal("esperava resposta válida")
	}
}

func TestCreateHarvest_PlantingNotFound(t *testing.T) {
	plantingRepo := &mockPlantingRepoForHarvest{err: errors.New("not found")}
	svc := NewHarvestService(&mockHarvestRepo{}, plantingRepo, &mockFieldRepoForPlanting{}, &mockTxRunner{})

	_, err := svc.CreateHarvest(validHarvestRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de plantio não encontrado")
	}
}

func TestCreateHarvest_Forbidden(t *testing.T) {
	plantingRepo := &mockPlantingRepoForHarvest{planting: activePlanting(), belongs: false}
	svc := NewHarvestService(&mockHarvestRepo{}, plantingRepo, &mockFieldRepoForPlanting{}, &mockTxRunner{})

	_, err := svc.CreateHarvest(validHarvestRequest(), 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestCreateHarvest_PlantingNotActive(t *testing.T) {
	planting := &entities.Planting{ID: 1, Status: entities.PlantingStatusHarvested, CreatedBy: 1}
	plantingRepo := &mockPlantingRepoForHarvest{planting: planting, belongs: true}
	svc := NewHarvestService(&mockHarvestRepo{}, plantingRepo, &mockFieldRepoForPlanting{}, &mockTxRunner{})

	_, err := svc.CreateHarvest(validHarvestRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro para plantio não ativo")
	}
	if !errors.Is(err, apperrors.ErrNoActivePlanting) {
		t.Errorf("esperava ErrNoActivePlanting, got %v", err)
	}
}

func TestCreateHarvest_AlreadyHarvested(t *testing.T) {
	plantingRepo := &mockPlantingRepoForHarvest{planting: activePlanting(), belongs: true}
	harvestRepo := &mockHarvestRepo{exists: true}
	svc := NewHarvestService(harvestRepo, plantingRepo, &mockFieldRepoForPlanting{}, &mockTxRunner{})

	_, err := svc.CreateHarvest(validHarvestRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro para colheita já existente")
	}
	if !errors.Is(err, apperrors.ErrAlreadyHarvested) {
		t.Errorf("esperava ErrAlreadyHarvested, got %v", err)
	}
}

func TestCreateHarvest_AdminCanHarvestAny(t *testing.T) {
	planting := &entities.Planting{ID: 1, FieldID: 1, Status: entities.PlantingStatusActive, CreatedBy: 99}
	plantingRepo := &mockPlantingRepoForHarvest{planting: planting, belongs: false}
	harvestRepo := &mockHarvestRepo{harvest: &entities.Harvest{ID: 1}}
	fieldRepo := &mockFieldRepoForPlanting{field: &entities.Field{ID: 1}}

	svc := NewHarvestService(harvestRepo, plantingRepo, fieldRepo, &mockTxRunner{})
	_, err := svc.CreateHarvest(validHarvestRequest(), 1, "admin")

	if err != nil {
		t.Fatalf("admin deveria registrar colheita de qualquer plantio, got %v", err)
	}
}

func TestGetHarvestByID_NotFound(t *testing.T) {
	svc := NewHarvestService(
		&mockHarvestRepo{err: errors.New("not found")},
		&mockPlantingRepoForHarvest{},
		&mockFieldRepoForPlanting{},
		&mockTxRunner{},
	)
	_, err := svc.GetHarvestByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetHarvestByID_Forbidden(t *testing.T) {
	harvest := &entities.Harvest{ID: 1, FieldID: 1}
	fieldRepo := &mockFieldRepoForPlanting{field: &entities.Field{ID: 1}, belongs: false}
	svc := NewHarvestService(
		&mockHarvestRepo{harvest: harvest},
		&mockPlantingRepoForHarvest{},
		fieldRepo,
		&mockTxRunner{},
	)
	_, err := svc.GetHarvestByID(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}
