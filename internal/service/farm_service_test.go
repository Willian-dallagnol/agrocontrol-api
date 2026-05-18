package service_test

import (
	"errors"
	"testing"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/service"
)

type mockFarmRepo struct {
	farm *entities.Farm
	err  error
}

func (m *mockFarmRepo) Create(f *entities.Farm) error {
	if m.err != nil {
		return m.err
	}
	f.ID = 1
	return nil
}
func (m *mockFarmRepo) FindByID(id uint) (*entities.Farm, error) { return m.farm, m.err }
func (m *mockFarmRepo) FindAll() ([]entities.Farm, error) {
	if m.farm != nil {
		return []entities.Farm{*m.farm}, nil
	}
	return nil, m.err
}
func (m *mockFarmRepo) FindByCreatedBy(userID uint) ([]entities.Farm, error) {
	if m.farm != nil {
		return []entities.Farm{*m.farm}, nil
	}
	return nil, m.err
}
func (m *mockFarmRepo) FindAllPaginated(offset, limit int, search string) ([]entities.Farm, int64, error) {
	if m.farm != nil {
		return []entities.Farm{*m.farm}, 1, nil
	}
	return nil, 0, m.err
}
func (m *mockFarmRepo) FindByCreatedByPaginated(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error) {
	if m.farm != nil {
		return []entities.Farm{*m.farm}, 1, nil
	}
	return nil, 0, m.err
}
func (m *mockFarmRepo) Update(f *entities.Farm) error { return m.err }
func (m *mockFarmRepo) Delete(id uint) error          { return m.err }
func (m *mockFarmRepo) Count() (int64, error)         { return 0, m.err }

func TestCreateFarm_Success(t *testing.T) {
	svc := service.NewFarmService(&mockFarmRepo{})
	resp, err := svc.CreateFarm(dto.CreateFarmRequest{
		Name:      "Fazenda Teste",
		OwnerName: "Willian",
		Location:  "Paraná",
		TotalArea: 100,
		City:      "São Pedro do Iguaçu",
		State:     "pr",
	}, 1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Fazenda Teste" {
		t.Errorf("esperava 'Fazenda Teste', got '%s'", resp.Name)
	}
	if resp.State != "PR" {
		t.Errorf("state deveria ser uppercase 'PR', got '%s'", resp.State)
	}
}

func TestCreateFarm_InvalidArea(t *testing.T) {
	svc := service.NewFarmService(&mockFarmRepo{})
	_, err := svc.CreateFarm(dto.CreateFarmRequest{
		Name:      "Fazenda Inválida",
		TotalArea: 0,
	}, 1)

	if err == nil {
		t.Fatal("esperava erro para total_area zero")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestCreateFarm_RepoError(t *testing.T) {
	svc := service.NewFarmService(&mockFarmRepo{err: errors.New("db error")})
	_, err := svc.CreateFarm(dto.CreateFarmRequest{
		Name:      "Fazenda",
		TotalArea: 50,
	}, 1)

	if err == nil {
		t.Fatal("esperava erro do repositório")
	}
}

func TestGetFarmByID_NotFound(t *testing.T) {
	svc := service.NewFarmService(&mockFarmRepo{err: errors.New("not found")})
	_, err := svc.GetFarmByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetFarmByID_Forbidden(t *testing.T) {
	farm := &entities.Farm{ID: 1, CreatedBy: 2}
	svc := service.NewFarmService(&mockFarmRepo{farm: farm})
	_, err := svc.GetFarmByID(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestGetFarmByID_AdminCanAccessAny(t *testing.T) {
	farm := &entities.Farm{ID: 1, Name: "Fazenda Admin", CreatedBy: 99}
	svc := service.NewFarmService(&mockFarmRepo{farm: farm})
	resp, err := svc.GetFarmByID(1, 1, "admin")

	if err != nil {
		t.Fatalf("admin deveria acessar qualquer fazenda, got %v", err)
	}
	if resp.Name != "Fazenda Admin" {
		t.Errorf("esperava 'Fazenda Admin', got '%s'", resp.Name)
	}
}

func TestDeleteFarm_Forbidden(t *testing.T) {
	farm := &entities.Farm{ID: 1, CreatedBy: 2}
	svc := service.NewFarmService(&mockFarmRepo{farm: farm})
	err := svc.DeleteFarm(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestDeleteFarm_AdminCanDeleteAny(t *testing.T) {
	farm := &entities.Farm{ID: 1, CreatedBy: 99}
	svc := service.NewFarmService(&mockFarmRepo{farm: farm})
	err := svc.DeleteFarm(1, 1, "admin")

	if err != nil {
		t.Fatalf("admin deveria deletar qualquer fazenda, got %v", err)
	}
}
