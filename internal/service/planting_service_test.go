package service

import (
	"errors"
	"testing"
	"time"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
)

// ── Mock PlantingRepository ────────────────────────────────────────────────

type mockPlantingRepo struct {
	planting *entities.Planting
	err      error
	belongs  bool
}

func (m *mockPlantingRepo) Create(p *entities.Planting) error {
	if m.err != nil {
		return m.err
	}
	p.ID = 1
	return nil
}
func (m *mockPlantingRepo) FindByID(id uint) (*entities.Planting, error) { return m.planting, m.err }
func (m *mockPlantingRepo) FindAll() ([]entities.Planting, error)        { return nil, m.err }
func (m *mockPlantingRepo) FindByUser(userID uint) ([]entities.Planting, error) {
	return nil, m.err
}
func (m *mockPlantingRepo) FindByFieldID(fieldID uint) ([]entities.Planting, error) {
	return nil, m.err
}
func (m *mockPlantingRepo) HasActivePlanting(fieldID uint) (bool, error) { return false, m.err }
func (m *mockPlantingRepo) CountActive() (int64, error)                  { return 0, m.err }
func (m *mockPlantingRepo) BelongsToUser(plantingID, userID uint) (bool, error) {
	return m.belongs, nil
}
func (m *mockPlantingRepo) Update(p *entities.Planting) error { return m.err }
func (m *mockPlantingRepo) Delete(id uint) error              { return m.err }

// ── Mock FieldRepository para Planting ────────────────────────────────────

type mockFieldRepoForPlanting struct {
	field   *entities.Field
	err     error
	belongs bool
}

func (m *mockFieldRepoForPlanting) Create(f *entities.Field) error            { return m.err }
func (m *mockFieldRepoForPlanting) FindByID(id uint) (*entities.Field, error) { return m.field, m.err }
func (m *mockFieldRepoForPlanting) FindAll() ([]entities.Field, error)        { return nil, m.err }
func (m *mockFieldRepoForPlanting) FindByUser(userID uint) ([]entities.Field, error) {
	return nil, m.err
}
func (m *mockFieldRepoForPlanting) FindByFarmID(farmID uint) ([]entities.Field, error) {
	return nil, m.err
}
func (m *mockFieldRepoForPlanting) FindAllPaginated(offset, limit int, search string) ([]entities.Field, int64, error) {
	return nil, 0, m.err
}
func (m *mockFieldRepoForPlanting) FindByUserPaginated(userID uint, offset, limit int, search string) ([]entities.Field, int64, error) {
	return nil, 0, m.err
}
func (m *mockFieldRepoForPlanting) ExistsByNameAndFarm(name string, farmID, excludeID uint) (bool, error) {
	return false, nil
}
func (m *mockFieldRepoForPlanting) Update(f *entities.Field) error { return m.err }
func (m *mockFieldRepoForPlanting) Delete(id uint) error           { return m.err }
func (m *mockFieldRepoForPlanting) Count() (int64, error)          { return 0, m.err }
func (m *mockFieldRepoForPlanting) BelongsToUser(fieldID, userID uint) (bool, error) {
	return m.belongs, nil
}

// ── Helpers ────────────────────────────────────────────────────────────────

func validPlantingRequest() dto.CreatePlantingRequest {
	return dto.CreatePlantingRequest{
		FieldID:         1,
		SeasonID:        1,
		CropID:          1,
		PlantingDate:    time.Now(),
		ExpectedHarvest: time.Now().AddDate(0, 4, 0),
		SeedsUsedKg:     50,
		Responsible:     "Willian",
	}
}

func activeFarm() *entities.Farm {
	return &entities.Farm{ID: 1, CreatedBy: 1}
}

// ── Planting tests ─────────────────────────────────────────────────────────

func TestCreatePlanting_Success(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1, Status: entities.FieldStatusActive, CreatedBy: 1},
		belongs: true,
	}
	seasonRepo := &mockSeasonRepo{season: &entities.Season{ID: 1}}
	cropRepo := &mockCropRepo{crop: &entities.Crop{ID: 1}}
	plantingRepo := &mockPlantingRepo{planting: &entities.Planting{ID: 1}}

	svc := NewPlantingService(plantingRepo, fieldRepo, seasonRepo, cropRepo)
	resp, err := svc.CreatePlanting(validPlantingRequest(), 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp == nil {
		t.Fatal("esperava resposta válida")
	}
}

func TestCreatePlanting_FieldNotFound(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{err: errors.New("not found")}
	svc := NewPlantingService(&mockPlantingRepo{}, fieldRepo, &mockSeasonRepo{}, &mockCropRepo{})

	_, err := svc.CreatePlanting(validPlantingRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de talhão não encontrado")
	}
}

func TestCreatePlanting_FieldInactive(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1, Status: entities.FieldStatusInactive, CreatedBy: 1},
		belongs: true,
	}
	svc := NewPlantingService(&mockPlantingRepo{}, fieldRepo, &mockSeasonRepo{season: &entities.Season{ID: 1}}, &mockCropRepo{crop: &entities.Crop{ID: 1}})

	_, err := svc.CreatePlanting(validPlantingRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro para talhão inativo")
	}
	if !errors.Is(err, apperrors.ErrInactiveField) {
		t.Errorf("esperava ErrInactiveField, got %v", err)
	}
}

func TestCreatePlanting_Forbidden(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1, Status: entities.FieldStatusActive, CreatedBy: 2},
		belongs: false,
	}
	svc := NewPlantingService(&mockPlantingRepo{}, fieldRepo, &mockSeasonRepo{}, &mockCropRepo{})

	_, err := svc.CreatePlanting(validPlantingRequest(), 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestCreatePlanting_SeasonNotFound(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1, Status: entities.FieldStatusActive},
		belongs: true,
	}
	seasonRepo := &mockSeasonRepo{err: errors.New("not found")}
	svc := NewPlantingService(&mockPlantingRepo{}, fieldRepo, seasonRepo, &mockCropRepo{})

	_, err := svc.CreatePlanting(validPlantingRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de safra não encontrada")
	}
}

func TestCreatePlanting_CropNotFound(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1, Status: entities.FieldStatusActive},
		belongs: true,
	}
	cropRepo := &mockCropRepo{err: errors.New("not found")}
	svc := NewPlantingService(&mockPlantingRepo{}, fieldRepo, &mockSeasonRepo{season: &entities.Season{ID: 1}}, cropRepo)

	_, err := svc.CreatePlanting(validPlantingRequest(), 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de cultura não encontrada")
	}
}

func TestGetPlantingByID_NotFound(t *testing.T) {
	svc := NewPlantingService(&mockPlantingRepo{err: errors.New("not found")}, &mockFieldRepoForPlanting{}, &mockSeasonRepo{}, &mockCropRepo{})
	_, err := svc.GetPlantingByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetPlantingByID_Forbidden(t *testing.T) {
	planting := &entities.Planting{ID: 1, CreatedBy: 2}
	svc := NewPlantingService(&mockPlantingRepo{planting: planting, belongs: false}, &mockFieldRepoForPlanting{}, &mockSeasonRepo{}, &mockCropRepo{})
	_, err := svc.GetPlantingByID(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestDeletePlanting_Success(t *testing.T) {
	planting := &entities.Planting{ID: 1, CreatedBy: 1}
	svc := NewPlantingService(&mockPlantingRepo{planting: planting, belongs: true}, &mockFieldRepoForPlanting{}, &mockSeasonRepo{}, &mockCropRepo{})
	err := svc.DeletePlanting(1, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
}
