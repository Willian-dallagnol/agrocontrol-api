package service

import (
	"errors"
	"testing"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
)

type mockFarmRepoForField struct {
	farm *entities.Farm
	err  error
}

func (m *mockFarmRepoForField) Create(f *entities.Farm) error            { return m.err }
func (m *mockFarmRepoForField) FindByID(id uint) (*entities.Farm, error) { return m.farm, m.err }
func (m *mockFarmRepoForField) FindAll() ([]entities.Farm, error)        { return nil, m.err }
func (m *mockFarmRepoForField) FindByCreatedBy(userID uint) ([]entities.Farm, error) {
	return nil, m.err
}
func (m *mockFarmRepoForField) FindAllPaginated(offset, limit int, search string) ([]entities.Farm, int64, error) {
	return nil, 0, m.err
}
func (m *mockFarmRepoForField) FindByCreatedByPaginated(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error) {
	return nil, 0, m.err
}
func (m *mockFarmRepoForField) Update(f *entities.Farm) error { return m.err }
func (m *mockFarmRepoForField) Delete(id uint) error          { return m.err }
func (m *mockFarmRepoForField) Count() (int64, error)         { return 0, m.err }

type mockFieldRepo struct {
	field   *entities.Field
	err     error
	exists  bool
	belongs bool
}

func (m *mockFieldRepo) Create(f *entities.Field) error {
	if m.err != nil {
		return m.err
	}
	f.ID = 1
	return nil
}
func (m *mockFieldRepo) FindByID(id uint) (*entities.Field, error)        { return m.field, m.err }
func (m *mockFieldRepo) FindAll() ([]entities.Field, error)               { return nil, m.err }
func (m *mockFieldRepo) FindByUser(userID uint) ([]entities.Field, error) { return nil, m.err }
func (m *mockFieldRepo) FindByFarmID(farmID uint) ([]entities.Field, error) {
	if m.field != nil {
		return []entities.Field{*m.field}, nil
	}
	return nil, m.err
}
func (m *mockFieldRepo) FindAllPaginated(offset, limit int, search string) ([]entities.Field, int64, error) {
	if m.field != nil {
		return []entities.Field{*m.field}, 1, nil
	}
	return nil, 0, m.err
}
func (m *mockFieldRepo) FindByUserPaginated(userID uint, offset, limit int, search string) ([]entities.Field, int64, error) {
	if m.field != nil {
		return []entities.Field{*m.field}, 1, nil
	}
	return nil, 0, m.err
}
func (m *mockFieldRepo) ExistsByNameAndFarm(name string, farmID, excludeID uint) (bool, error) {
	return m.exists, m.err
}
func (m *mockFieldRepo) Update(f *entities.Field) error { return m.err }
func (m *mockFieldRepo) Delete(id uint) error           { return m.err }
func (m *mockFieldRepo) Count() (int64, error)          { return 0, m.err }
func (m *mockFieldRepo) BelongsToUser(fieldID, userID uint) (bool, error) {
	return m.belongs, nil
}

func TestCreateField_Success(t *testing.T) {
	farmRepo := &mockFarmRepoForField{farm: &entities.Farm{ID: 1, CreatedBy: 1}}
	fieldRepo := &mockFieldRepo{}
	svc := NewFieldService(fieldRepo, farmRepo)

	resp, err := svc.CreateField(dto.CreateFieldRequest{
		Name:   "Talhão Norte",
		Area:   50,
		FarmID: 1,
	}, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Talhão Norte" {
		t.Errorf("esperava 'Talhão Norte', got '%s'", resp.Name)
	}
}

func TestCreateField_InvalidArea(t *testing.T) {
	farmRepo := &mockFarmRepoForField{farm: &entities.Farm{ID: 1, CreatedBy: 1}}
	fieldRepo := &mockFieldRepo{}
	svc := NewFieldService(fieldRepo, farmRepo)

	_, err := svc.CreateField(dto.CreateFieldRequest{
		Name:   "Talhão Inválido",
		Area:   0,
		FarmID: 1,
	}, 1, "manager")

	if err == nil {
		t.Fatal("esperava erro para area zero")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestCreateField_FarmNotFound(t *testing.T) {
	farmRepo := &mockFarmRepoForField{err: errors.New("not found")}
	fieldRepo := &mockFieldRepo{}
	svc := NewFieldService(fieldRepo, farmRepo)

	_, err := svc.CreateField(dto.CreateFieldRequest{
		Name:   "Talhão",
		Area:   50,
		FarmID: 99,
	}, 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de fazenda não encontrada")
	}
}

func TestCreateField_Forbidden(t *testing.T) {
	farmRepo := &mockFarmRepoForField{farm: &entities.Farm{ID: 1, CreatedBy: 2}}
	fieldRepo := &mockFieldRepo{}
	svc := NewFieldService(fieldRepo, farmRepo)

	_, err := svc.CreateField(dto.CreateFieldRequest{
		Name:   "Talhão",
		Area:   50,
		FarmID: 1,
	}, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestCreateField_DuplicateName(t *testing.T) {
	farmRepo := &mockFarmRepoForField{farm: &entities.Farm{ID: 1, CreatedBy: 1}}
	fieldRepo := &mockFieldRepo{exists: true}
	svc := NewFieldService(fieldRepo, farmRepo)

	_, err := svc.CreateField(dto.CreateFieldRequest{
		Name:   "Talhão Duplicado",
		Area:   50,
		FarmID: 1,
	}, 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de duplicidade")
	}
}

func TestGetFieldByID_NotFound(t *testing.T) {
	svc := NewFieldService(&mockFieldRepo{err: errors.New("not found")}, &mockFarmRepoForField{})
	_, err := svc.GetFieldByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetFieldByID_AdminCanAccessAny(t *testing.T) {
	field := &entities.Field{ID: 1, Name: "Talhão Admin", CreatedBy: 99}
	svc := NewFieldService(&mockFieldRepo{field: field, belongs: false}, &mockFarmRepoForField{})
	resp, err := svc.GetFieldByID(1, 1, "admin")

	if err != nil {
		t.Fatalf("admin deveria acessar qualquer talhão, got %v", err)
	}
	if resp.Name != "Talhão Admin" {
		t.Errorf("esperava 'Talhão Admin', got '%s'", resp.Name)
	}
}

func TestDeleteField_Forbidden(t *testing.T) {
	field := &entities.Field{ID: 1, CreatedBy: 2}
	svc := NewFieldService(&mockFieldRepo{field: field, belongs: false}, &mockFarmRepoForField{})
	err := svc.DeleteField(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestDeleteField_Success(t *testing.T) {
	field := &entities.Field{ID: 1, CreatedBy: 1}
	svc := NewFieldService(&mockFieldRepo{field: field, belongs: true}, &mockFarmRepoForField{})
	err := svc.DeleteField(1, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
}
