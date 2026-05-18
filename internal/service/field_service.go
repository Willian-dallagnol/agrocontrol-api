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

type FieldService struct {
	Repo     ports.FieldRepository
	FarmRepo ports.FarmRepository
}

func NewFieldService(repo ports.FieldRepository, farmRepo ports.FarmRepository) *FieldService {
	return &FieldService{Repo: repo, FarmRepo: farmRepo}
}

func (s *FieldService) CreateField(req dto.CreateFieldRequest, userID uint, role string) (*dto.FieldResponse, error) {
	if req.Area <= 0 {
		return nil, fmt.Errorf("area deve ser maior que zero: %w", apperrors.ErrInvalidInput)
	}
	farm, err := s.FarmRepo.FindByID(req.FarmID)
	if err != nil {
		return nil, apperrors.NotFoundError("fazenda", req.FarmID)
	}
	if role != "admin" && farm.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("criar talhão nesta fazenda")
	}
	exists, err := s.Repo.ExistsByNameAndFarm(req.Name, req.FarmID, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar duplicidade: %w", err)
	}
	if exists {
		return nil, apperrors.ConflictError(fmt.Sprintf("talhão '%s' nesta fazenda", req.Name))
	}
	status := req.Status
	if status == "" {
		status = entities.FieldStatusActive
	}
	field := &entities.Field{
		Name:      strings.TrimSpace(req.Name),
		Area:      req.Area,
		SoilType:  strings.TrimSpace(req.SoilType),
		Status:    status,
		FarmID:    req.FarmID,
		CreatedBy: userID,
	}
	if err := s.Repo.Create(field); err != nil {
		return nil, fmt.Errorf("erro ao criar talhão: %w", err)
	}
	slog.Info("field: criado", "field_id", field.ID, "farm_id", req.FarmID, "user_id", userID)
	return toFieldResponse(field), nil
}

func (s *FieldService) GetFieldsPaginated(q dto.FieldQuery, userID uint, role string) (dto.PaginatedResponse, error) {
	var fields []entities.Field
	var total int64
	var err error
	if role == "admin" {
		fields, total, err = s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.Search)
	} else {
		fields, total, err = s.Repo.FindByUserPaginated(userID, q.Offset(), q.Limit, q.Search)
	}
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar talhões: %w", err)
	}
	resp := make([]dto.FieldResponse, 0, len(fields))
	for i := range fields {
		resp = append(resp, *toFieldResponse(&fields[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *FieldService) GetFieldByID(id uint, userID uint, role string) (*dto.FieldResponse, error) {
	field, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("talhão", id)
	}
	if role != "admin" {
		ok, _ := s.Repo.BelongsToUser(id, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar este talhão")
		}
	}
	return toFieldResponse(field), nil
}

func (s *FieldService) GetFieldsByFarmID(farmID uint, userID uint, role string) ([]dto.FieldResponse, error) {
	farm, err := s.FarmRepo.FindByID(farmID)
	if err != nil {
		return nil, apperrors.NotFoundError("fazenda", farmID)
	}
	if role != "admin" && farm.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("acessar talhões desta fazenda")
	}
	fields, err := s.Repo.FindByFarmID(farmID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar talhões: %w", err)
	}
	resp := make([]dto.FieldResponse, 0, len(fields))
	for i := range fields {
		resp = append(resp, *toFieldResponse(&fields[i]))
	}
	return resp, nil
}

func (s *FieldService) UpdateField(id uint, req dto.UpdateFieldRequest, userID uint, role string) (*dto.FieldResponse, error) {
	field, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("talhão", id)
	}
	if role != "admin" {
		ok, _ := s.Repo.BelongsToUser(id, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("editar este talhão")
		}
	}
	if req.Area <= 0 {
		return nil, fmt.Errorf("area deve ser maior que zero: %w", apperrors.ErrInvalidInput)
	}
	if _, err := s.FarmRepo.FindByID(req.FarmID); err != nil {
		return nil, apperrors.NotFoundError("fazenda", req.FarmID)
	}
	exists, err := s.Repo.ExistsByNameAndFarm(req.Name, req.FarmID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar duplicidade: %w", err)
	}
	if exists {
		return nil, apperrors.ConflictError(fmt.Sprintf("talhão '%s' nesta fazenda", req.Name))
	}
	field.Name     = strings.TrimSpace(req.Name)
	field.Area     = req.Area
	field.SoilType = strings.TrimSpace(req.SoilType)
	field.Status   = req.Status
	field.FarmID   = req.FarmID
	if err := s.Repo.Update(field); err != nil {
		return nil, fmt.Errorf("erro ao atualizar talhão: %w", err)
	}
	return toFieldResponse(field), nil
}

func (s *FieldService) DeleteField(id uint, userID uint, role string) error {
	if _, err := s.Repo.FindByID(id); err != nil {
		return apperrors.NotFoundError("talhão", id)
	}
	if role != "admin" {
		ok, _ := s.Repo.BelongsToUser(id, userID)
		if !ok {
			return apperrors.ForbiddenError("excluir este talhão")
		}
	}
	if err := s.Repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao excluir talhão: %w", err)
	}
	slog.Info("field: excluído", "field_id", id, "by_user", userID)
	return nil
}

func toFieldResponse(f *entities.Field) *dto.FieldResponse {
	return &dto.FieldResponse{
		ID: f.ID, Name: f.Name, Area: f.Area, SoilType: f.SoilType,
		Status: f.Status, FarmID: f.FarmID, CreatedBy: f.CreatedBy,
		CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}
