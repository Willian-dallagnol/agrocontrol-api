package service

import (
	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/dto"
	"fmt"
	"log/slog"
	"strings"
)

// ── CropService ───────────────────────────────────────────────────────────────

type CropService struct{ Repo ports.CropRepository }

func NewCropService(repo ports.CropRepository) *CropService { return &CropService{Repo: repo} }

func (s *CropService) CreateCrop(req dto.CreateCropRequest, userID uint) (*dto.CropResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("nome da cultura é obrigatório: %w", apperrors.ErrInvalidInput)
	}
	crop := &entities.Crop{
		Name: strings.TrimSpace(req.Name), Variety: strings.TrimSpace(req.Variety),
		Type: req.Type, CycleDays: req.CycleDays,
		SpacingCm: req.SpacingCm, PlantPopulation: req.PlantPopulation, CreatedBy: userID,
	}
	if err := s.Repo.Create(crop); err != nil {
		return nil, fmt.Errorf("erro ao criar cultura: %w", err)
	}
	return toCropResponse(crop), nil
}

func (s *CropService) GetCropsPaginated(q dto.CropQuery) (dto.PaginatedResponse, error) {
	crops, total, err := s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.Search)
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar culturas: %w", err)
	}
	resp := make([]dto.CropResponse, 0, len(crops))
	for i := range crops {
		resp = append(resp, *toCropResponse(&crops[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *CropService) GetCropByID(id uint) (*dto.CropResponse, error) {
	crop, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("cultura", id)
	}
	return toCropResponse(crop), nil
}

func (s *CropService) UpdateCrop(id uint, req dto.UpdateCropRequest, role string) (*dto.CropResponse, error) {
	if role != "admin" && role != "manager" {
		return nil, apperrors.ForbiddenError("editar cultura")
	}
	crop, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("cultura", id)
	}
	crop.Name            = strings.TrimSpace(req.Name)
	crop.Variety         = strings.TrimSpace(req.Variety)
	crop.Type            = req.Type
	crop.CycleDays       = req.CycleDays
	crop.SpacingCm       = req.SpacingCm
	crop.PlantPopulation = req.PlantPopulation
	if err := s.Repo.Update(crop); err != nil {
		return nil, fmt.Errorf("erro ao atualizar cultura: %w", err)
	}
	return toCropResponse(crop), nil
}

func (s *CropService) DeleteCrop(id uint) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		return apperrors.NotFoundError("cultura", id)
	}
	return s.Repo.Delete(id)
}

func toCropResponse(c *entities.Crop) *dto.CropResponse {
	return &dto.CropResponse{
		ID: c.ID, Name: c.Name, Variety: c.Variety, Type: c.Type,
		CycleDays: c.CycleDays, SpacingCm: c.SpacingCm, PlantPopulation: c.PlantPopulation,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

// ── SeasonService ─────────────────────────────────────────────────────────────

type SeasonService struct{ Repo ports.SeasonRepository }

func NewSeasonService(repo ports.SeasonRepository) *SeasonService {
	return &SeasonService{Repo: repo}
}

func (s *SeasonService) CreateSeason(req dto.CreateSeasonRequest, userID uint) (*dto.SeasonResponse, error) {
	if !req.EndDate.After(req.StartDate) {
		return nil, fmt.Errorf("data_fim deve ser posterior a data_inicio: %w", apperrors.ErrInvalidInput)
	}
	status := req.Status
	if status == "" {
		status = entities.SeasonStatusPlanning
	}
	season := &entities.Season{
		Name:      strings.TrimSpace(req.Name),
		StartDate: req.StartDate, EndDate: req.EndDate,
		Status: status, CreatedBy: userID,
	}
	if err := s.Repo.Create(season); err != nil {
		return nil, fmt.Errorf("erro ao criar safra: %w", err)
	}
	slog.Info("season: criada", "season_id", season.ID, "name", season.Name)
	return toSeasonResponse(season), nil
}

func (s *SeasonService) GetSeasonsPaginated(q dto.SeasonQuery) (dto.PaginatedResponse, error) {
	seasons, total, err := s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.Search)
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar safras: %w", err)
	}
	resp := make([]dto.SeasonResponse, 0, len(seasons))
	for i := range seasons {
		resp = append(resp, *toSeasonResponse(&seasons[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *SeasonService) GetSeasonByID(id uint) (*dto.SeasonResponse, error) {
	season, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("safra", id)
	}
	return toSeasonResponse(season), nil
}

func (s *SeasonService) UpdateSeason(id uint, req dto.UpdateSeasonRequest) (*dto.SeasonResponse, error) {
	if !req.EndDate.After(req.StartDate) {
		return nil, fmt.Errorf("data_fim deve ser posterior a data_inicio: %w", apperrors.ErrInvalidInput)
	}
	season, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("safra", id)
	}
	season.Name      = strings.TrimSpace(req.Name)
	season.StartDate = req.StartDate
	season.EndDate   = req.EndDate
	season.Status    = req.Status
	if err := s.Repo.Update(season); err != nil {
		return nil, fmt.Errorf("erro ao atualizar safra: %w", err)
	}
	return toSeasonResponse(season), nil
}

func (s *SeasonService) DeleteSeason(id uint) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		return apperrors.NotFoundError("safra", id)
	}
	return s.Repo.Delete(id)
}

func toSeasonResponse(s *entities.Season) *dto.SeasonResponse {
	return &dto.SeasonResponse{
		ID: s.ID, Name: s.Name, StartDate: s.StartDate, EndDate: s.EndDate,
		Status: s.Status, CreatedBy: s.CreatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

// ── PlantingService ───────────────────────────────────────────────────────────

type PlantingService struct {
	Repo       ports.PlantingRepository
	FieldRepo  ports.FieldRepository
	SeasonRepo ports.SeasonRepository
	CropRepo   ports.CropRepository
}

func NewPlantingService(
	repo ports.PlantingRepository, fieldRepo ports.FieldRepository,
	seasonRepo ports.SeasonRepository, cropRepo ports.CropRepository,
) *PlantingService {
	return &PlantingService{Repo: repo, FieldRepo: fieldRepo, SeasonRepo: seasonRepo, CropRepo: cropRepo}
}

func (s *PlantingService) CreatePlanting(req dto.CreatePlantingRequest, userID uint, role string) (*dto.PlantingResponse, error) {
	field, err := s.FieldRepo.FindByID(req.FieldID)
	if err != nil {
		return nil, apperrors.NotFoundError("talhão", req.FieldID)
	}
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(req.FieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("plantar neste talhão")
		}
	}
	if field.Status != entities.FieldStatusActive {
		return nil, apperrors.ErrInactiveField
	}
	if _, err := s.SeasonRepo.FindByID(req.SeasonID); err != nil {
		return nil, apperrors.NotFoundError("safra", req.SeasonID)
	}
	if _, err := s.CropRepo.FindByID(req.CropID); err != nil {
		return nil, apperrors.NotFoundError("cultura", req.CropID)
	}
	planting := &entities.Planting{
		FieldID: req.FieldID, SeasonID: req.SeasonID, CropID: req.CropID,
		PlantingDate: req.PlantingDate, ExpectedHarvest: req.ExpectedHarvest,
		SeedsUsedKg: req.SeedsUsedKg, DensityKgHa: req.DensityKgHa,
		DepthCm: req.DepthCm, Spacing: req.Spacing,
		Responsible: strings.TrimSpace(req.Responsible),
		Status:      entities.PlantingStatusActive,
		Notes:       req.Notes, CreatedBy: userID,
	}
	if err := s.Repo.Create(planting); err != nil {
		return nil, fmt.Errorf("erro ao criar plantio: %w", err)
	}
	slog.Info("planting: criado", "planting_id", planting.ID, "field_id", req.FieldID)
	full, _ := s.Repo.FindByID(planting.ID)
	return toPlantingResponse(full), nil
}

func (s *PlantingService) GetPlantings(userID uint, role string) ([]dto.PlantingResponse, error) {
	var plantings []entities.Planting
	var err error
	if role == "admin" {
		plantings, err = s.Repo.FindAll()
	} else {
		plantings, err = s.Repo.FindByUser(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar plantios: %w", err)
	}
	resp := make([]dto.PlantingResponse, 0, len(plantings))
	for i := range plantings {
		resp = append(resp, *toPlantingResponse(&plantings[i]))
	}
	return resp, nil
}

func (s *PlantingService) GetPlantingByID(id uint, userID uint, role string) (*dto.PlantingResponse, error) {
	p, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("plantio", id)
	}
	if role != "admin" {
		ok, _ := s.Repo.BelongsToUser(id, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar este plantio")
		}
	}
	return toPlantingResponse(p), nil
}

func (s *PlantingService) UpdatePlanting(id uint, req dto.UpdatePlantingRequest, userID uint, role string) (*dto.PlantingResponse, error) {
	p, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("plantio", id)
	}
	if role != "admin" {
		ok, _ := s.Repo.BelongsToUser(id, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("editar este plantio")
		}
	}
	if !req.ExpectedHarvest.IsZero() { p.ExpectedHarvest = req.ExpectedHarvest }
	if req.Status != ""              { p.Status = req.Status }
	if req.Responsible != ""         { p.Responsible = strings.TrimSpace(req.Responsible) }
	if req.Notes != ""               { p.Notes = req.Notes }
	if err := s.Repo.Update(p); err != nil {
		return nil, fmt.Errorf("erro ao atualizar plantio: %w", err)
	}
	full, _ := s.Repo.FindByID(p.ID)
	return toPlantingResponse(full), nil
}

func (s *PlantingService) DeletePlanting(id uint, userID uint, role string) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		return apperrors.NotFoundError("plantio", id)
	}
	if role != "admin" {
		ok, _ := s.Repo.BelongsToUser(id, userID)
		if !ok {
			return apperrors.ForbiddenError("excluir este plantio")
		}
	}
	return s.Repo.Delete(id)
}

func toPlantingResponse(p *entities.Planting) *dto.PlantingResponse {
	r := &dto.PlantingResponse{
		ID: p.ID, FieldID: p.FieldID, SeasonID: p.SeasonID, CropID: p.CropID,
		PlantingDate: p.PlantingDate, ExpectedHarvest: p.ExpectedHarvest,
		SeedsUsedKg: p.SeedsUsedKg, DensityKgHa: p.DensityKgHa, DepthCm: p.DepthCm,
		Spacing: p.Spacing, Responsible: p.Responsible, Status: p.Status,
		Notes: p.Notes, CreatedBy: p.CreatedBy, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
	if p.Field.ID > 0  { r.FieldName = p.Field.Name }
	if p.Season.ID > 0 { r.SeasonName = p.Season.Name }
	if p.Crop.ID > 0   { r.CropName = p.Crop.Name }
	return r
}
