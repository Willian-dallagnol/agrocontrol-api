package service

import (
	"errors"
	"testing"
	"time"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
)

// ── Mock CropRepository ────────────────────────────────────────────────────

type mockCropRepo struct {
	crop *entities.Crop
	err  error
}

func (m *mockCropRepo) Create(c *entities.Crop) error {
	if m.err != nil {
		return m.err
	}
	c.ID = 1
	return nil
}
func (m *mockCropRepo) FindByID(id uint) (*entities.Crop, error) { return m.crop, m.err }
func (m *mockCropRepo) FindAll() ([]entities.Crop, error)        { return nil, m.err }
func (m *mockCropRepo) FindAllPaginated(offset, limit int, search string) ([]entities.Crop, int64, error) {
	if m.crop != nil {
		return []entities.Crop{*m.crop}, 1, nil
	}
	return nil, 0, m.err
}
func (m *mockCropRepo) Update(c *entities.Crop) error { return m.err }
func (m *mockCropRepo) Delete(id uint) error          { return m.err }

// ── Crop tests ─────────────────────────────────────────────────────────────

func TestCreateCrop_Success(t *testing.T) {
	svc := NewCropService(&mockCropRepo{})
	resp, err := svc.CreateCrop(dto.CreateCropRequest{
		Name:    "Soja",
		Variety: "Intacta",
		Type:    "Grão",
	}, 1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Soja" {
		t.Errorf("esperava 'Soja', got '%s'", resp.Name)
	}
}

func TestCreateCrop_EmptyName(t *testing.T) {
	svc := NewCropService(&mockCropRepo{})
	_, err := svc.CreateCrop(dto.CreateCropRequest{Name: "  "}, 1)

	if err == nil {
		t.Fatal("esperava erro para nome vazio")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestCreateCrop_RepoError(t *testing.T) {
	svc := NewCropService(&mockCropRepo{err: errors.New("db error")})
	_, err := svc.CreateCrop(dto.CreateCropRequest{Name: "Milho"}, 1)

	if err == nil {
		t.Fatal("esperava erro do repositório")
	}
}

func TestGetCropByID_NotFound(t *testing.T) {
	svc := NewCropService(&mockCropRepo{err: errors.New("not found")})
	_, err := svc.GetCropByID(99)

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestUpdateCrop_Forbidden(t *testing.T) {
	svc := NewCropService(&mockCropRepo{crop: &entities.Crop{ID: 1}})
	_, err := svc.UpdateCrop(1, dto.UpdateCropRequest{Name: "Soja"}, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden para operator")
	}
}

func TestUpdateCrop_AdminSuccess(t *testing.T) {
	svc := NewCropService(&mockCropRepo{crop: &entities.Crop{ID: 1, Name: "Soja"}})
	resp, err := svc.UpdateCrop(1, dto.UpdateCropRequest{Name: "Milho"}, "admin")

	if err != nil {
		t.Fatalf("admin deveria atualizar cultura, got %v", err)
	}
	if resp.Name != "Milho" {
		t.Errorf("esperava 'Milho', got '%s'", resp.Name)
	}
}

func TestDeleteCrop_NotFound(t *testing.T) {
	svc := NewCropService(&mockCropRepo{err: errors.New("not found")})
	err := svc.DeleteCrop(99)

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

// ── Mock SeasonRepository ──────────────────────────────────────────────────

type mockSeasonRepo struct {
	season *entities.Season
	err    error
}

func (m *mockSeasonRepo) Create(s *entities.Season) error {
	if m.err != nil {
		return m.err
	}
	s.ID = 1
	return nil
}
func (m *mockSeasonRepo) FindByID(id uint) (*entities.Season, error) { return m.season, m.err }
func (m *mockSeasonRepo) FindAll() ([]entities.Season, error)        { return nil, m.err }
func (m *mockSeasonRepo) FindAllPaginated(offset, limit int, search string) ([]entities.Season, int64, error) {
	if m.season != nil {
		return []entities.Season{*m.season}, 1, nil
	}
	return nil, 0, m.err
}
func (m *mockSeasonRepo) Update(s *entities.Season) error { return m.err }
func (m *mockSeasonRepo) Delete(id uint) error            { return m.err }

// ── Season tests ───────────────────────────────────────────────────────────

func TestCreateSeason_Success(t *testing.T) {
	svc := NewSeasonService(&mockSeasonRepo{})
	start := time.Now()
	end := start.AddDate(0, 6, 0)

	resp, err := svc.CreateSeason(dto.CreateSeasonRequest{
		Name:      "Safra 2026/1",
		StartDate: start,
		EndDate:   end,
	}, 1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Safra 2026/1" {
		t.Errorf("esperava 'Safra 2026/1', got '%s'", resp.Name)
	}
}

func TestCreateSeason_InvalidDates(t *testing.T) {
	svc := NewSeasonService(&mockSeasonRepo{})
	now := time.Now()

	_, err := svc.CreateSeason(dto.CreateSeasonRequest{
		Name:      "Safra Inválida",
		StartDate: now,
		EndDate:   now.AddDate(0, 0, -1),
	}, 1)

	if err == nil {
		t.Fatal("esperava erro para data_fim anterior a data_inicio")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestCreateSeason_EqualDates(t *testing.T) {
	svc := NewSeasonService(&mockSeasonRepo{})
	now := time.Now()

	_, err := svc.CreateSeason(dto.CreateSeasonRequest{
		Name:      "Safra Inválida",
		StartDate: now,
		EndDate:   now,
	}, 1)

	if err == nil {
		t.Fatal("esperava erro para datas iguais")
	}
}

func TestGetSeasonByID_NotFound(t *testing.T) {
	svc := NewSeasonService(&mockSeasonRepo{err: errors.New("not found")})
	_, err := svc.GetSeasonByID(99)

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestDeleteSeason_NotFound(t *testing.T) {
	svc := NewSeasonService(&mockSeasonRepo{err: errors.New("not found")})
	err := svc.DeleteSeason(99)

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestDeleteSeason_Success(t *testing.T) {
	svc := NewSeasonService(&mockSeasonRepo{season: &entities.Season{ID: 1}})
	err := svc.DeleteSeason(1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
}
