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

type FarmService struct {
	Repo ports.FarmRepository
}

func NewFarmService(repo ports.FarmRepository) *FarmService {
	return &FarmService{Repo: repo}
}

func (s *FarmService) CreateFarm(req dto.CreateFarmRequest, userID uint) (*dto.FarmResponse, error) {
	if req.TotalArea <= 0 {
		return nil, fmt.Errorf("total_area deve ser maior que zero: %w", apperrors.ErrInvalidInput)
	}
	farm := &entities.Farm{
		Name:      strings.TrimSpace(req.Name),
		OwnerName: strings.TrimSpace(req.OwnerName),
		Location:  strings.TrimSpace(req.Location),
		TotalArea: req.TotalArea,
		City:      strings.TrimSpace(req.City),
		State:     strings.ToUpper(strings.TrimSpace(req.State)),
		CreatedBy: userID,
	}
	if err := s.Repo.Create(farm); err != nil {
		slog.Error("farm: erro ao criar", "user_id", userID, "error", err)
		return nil, fmt.Errorf("erro ao criar fazenda: %w", err)
	}
	slog.Info("farm: criada", "farm_id", farm.ID, "user_id", userID)
	return toFarmResponse(farm), nil
}

func (s *FarmService) GetFarmsPaginated(q dto.FarmQuery, userID uint, role string) (dto.PaginatedResponse, error) {
	var farms []entities.Farm
	var total int64
	var err error

	if role == "admin" {
		farms, total, err = s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.Search)
	} else {
		farms, total, err = s.Repo.FindByCreatedByPaginated(userID, q.Offset(), q.Limit, q.Search)
	}
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar fazendas: %w", err)
	}

	resp := make([]dto.FarmResponse, 0, len(farms))
	for i := range farms {
		resp = append(resp, *toFarmResponse(&farms[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *FarmService) GetFarmByID(id uint, userID uint, role string) (*dto.FarmResponse, error) {
	farm, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("fazenda", id)
	}
	if role != "admin" && farm.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("acessar esta fazenda")
	}
	return toFarmResponse(farm), nil
}

func (s *FarmService) UpdateFarm(id uint, req dto.UpdateFarmRequest, userID uint, role string) (*dto.FarmResponse, error) {
	farm, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("fazenda", id)
	}
	if role != "admin" && farm.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("editar esta fazenda")
	}
	if req.TotalArea <= 0 {
		return nil, fmt.Errorf("total_area deve ser maior que zero: %w", apperrors.ErrInvalidInput)
	}
	farm.Name      = strings.TrimSpace(req.Name)
	farm.OwnerName = strings.TrimSpace(req.OwnerName)
	farm.Location  = strings.TrimSpace(req.Location)
	farm.TotalArea = req.TotalArea
	farm.City      = strings.TrimSpace(req.City)
	farm.State     = strings.ToUpper(strings.TrimSpace(req.State))

	if err := s.Repo.Update(farm); err != nil {
		return nil, fmt.Errorf("erro ao atualizar fazenda: %w", err)
	}
	return toFarmResponse(farm), nil
}

func (s *FarmService) DeleteFarm(id uint, userID uint, role string) error {
	farm, err := s.Repo.FindByID(id)
	if err != nil {
		return apperrors.NotFoundError("fazenda", id)
	}
	if role != "admin" && farm.CreatedBy != userID {
		return apperrors.ForbiddenError("excluir esta fazenda")
	}
	if err := s.Repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao excluir fazenda: %w", err)
	}
	slog.Info("farm: excluída", "farm_id", id, "by_user", userID)
	return nil
}

func toFarmResponse(f *entities.Farm) *dto.FarmResponse {
	return &dto.FarmResponse{
		ID: f.ID, Name: f.Name, OwnerName: f.OwnerName, Location: f.Location,
		TotalArea: f.TotalArea, City: f.City, State: f.State,
		CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}
