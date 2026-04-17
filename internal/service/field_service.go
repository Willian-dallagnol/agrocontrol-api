package service

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"errors"
)

type FieldService struct {
	Repo     *repository.FieldRepository
	FarmRepo *repository.FarmRepository
}

func NewFieldService(repo *repository.FieldRepository, farmRepo *repository.FarmRepository) *FieldService {
	return &FieldService{
		Repo:     repo,
		FarmRepo: farmRepo,
	}
}

func (s *FieldService) CreateField(req dto.CreateFieldRequest) (*dto.FieldResponse, error) {

	// 🔹 validação de regra
	if req.Area <= 0 {
		return nil, errors.New("a área do talhão deve ser maior que zero")
	}

	// 🔹 valida se a farm existe
	_, err := s.FarmRepo.FindByID(req.FarmID)
	if err != nil {
		return nil, errors.New("fazenda não encontrada")
	}

	field := &entities.Field{
		Name:     req.Name,
		Area:     req.Area,
		SoilType: req.SoilType,
		FarmID:   req.FarmID,
	}

	err = s.Repo.Create(field)
	if err != nil {
		return nil, err
	}

	return &dto.FieldResponse{
		ID:       field.ID,
		Name:     field.Name,
		Area:     field.Area,
		SoilType: field.SoilType,
		FarmID:   field.FarmID,
	}, nil
}

func (s *FieldService) GetFields() ([]dto.FieldResponse, error) {
	fields, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.FieldResponse

	for _, f := range fields {
		response = append(response, dto.FieldResponse{
			ID:       f.ID,
			Name:     f.Name,
			Area:     f.Area,
			SoilType: f.SoilType,
			FarmID:   f.FarmID,
		})
	}

	return response, nil
}

func (s *FieldService) GetFieldByID(id uint) (*dto.FieldResponse, error) {
	field, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.FieldResponse{
		ID:       field.ID,
		Name:     field.Name,
		Area:     field.Area,
		SoilType: field.SoilType,
		FarmID:   field.FarmID,
	}, nil
}

func (s *FieldService) UpdateField(id uint, req dto.UpdateFieldRequest) (*dto.FieldResponse, error) {

	if req.Area <= 0 {
		return nil, errors.New("a área do talhão deve ser maior que zero")
	}

	// 🔹 valida farm
	_, err := s.FarmRepo.FindByID(req.FarmID)
	if err != nil {
		return nil, errors.New("fazenda não encontrada")
	}

	field, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	field.Name = req.Name
	field.Area = req.Area
	field.SoilType = req.SoilType
	field.FarmID = req.FarmID

	err = s.Repo.Update(field)
	if err != nil {
		return nil, err
	}

	return &dto.FieldResponse{
		ID:       field.ID,
		Name:     field.Name,
		Area:     field.Area,
		SoilType: field.SoilType,
		FarmID:   field.FarmID,
	}, nil
}

func (s *FieldService) DeleteField(id uint) error {
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.Repo.Delete(id)
}

func (s *FieldService) GetFieldsByFarmID(farmID uint) ([]dto.FieldResponse, error) {

	// 🔹 valida se a farm existe
	_, err := s.FarmRepo.FindByID(farmID)
	if err != nil {
		return nil, errors.New("fazenda não encontrada")
	}

	fields, err := s.Repo.FindByFarmID(farmID)
	if err != nil {
		return nil, err
	}

	var response []dto.FieldResponse

	for _, f := range fields {
		response = append(response, dto.FieldResponse{
			ID:       f.ID,
			Name:     f.Name,
			Area:     f.Area,
			SoilType: f.SoilType,
			FarmID:   f.FarmID,
		})
	}

	return response, nil
}
