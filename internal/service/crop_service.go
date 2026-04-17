package service

import (
	"errors"

	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
)

type CropService struct {
	Repo      *repository.CropRepository
	FieldRepo *repository.FieldRepository
}

func NewCropService(repo *repository.CropRepository, fieldRepo *repository.FieldRepository) *CropService {
	return &CropService{
		Repo:      repo,
		FieldRepo: fieldRepo,
	}
}

func (s *CropService) CreateCrop(req dto.CreateCropRequest) (*dto.CropResponse, error) {

	// 🔥 valida se o field existe
	_, err := s.FieldRepo.FindByID(req.FieldID)
	if err != nil {
		return nil, errors.New("talhão não encontrado")
	}

	crop := entities.Crop{
		Name:    req.Name,
		Type:    req.Type,
		FieldID: req.FieldID,
	}

	err = s.Repo.Create(&crop)
	if err != nil {
		return nil, err
	}

	return &dto.CropResponse{
		ID:      crop.ID,
		Name:    crop.Name,
		Type:    crop.Type,
		FieldID: crop.FieldID,
	}, nil
}

func (s *CropService) GetCrops() ([]dto.CropResponse, error) {
	crops, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.CropResponse

	for _, c := range crops {
		response = append(response, dto.CropResponse{
			ID:      c.ID,
			Name:    c.Name,
			Type:    c.Type,
			FieldID: c.FieldID,
		})
	}

	return response, nil
}

func (s *CropService) GetCropByID(id uint) (*dto.CropResponse, error) {
	crop, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.CropResponse{
		ID:      crop.ID,
		Name:    crop.Name,
		Type:    crop.Type,
		FieldID: crop.FieldID,
	}, nil
}

func (s *CropService) UpdateCrop(id uint, req dto.UpdateCropRequest) (*dto.CropResponse, error) {

	// 🔥 valida field
	_, err := s.FieldRepo.FindByID(req.FieldID)
	if err != nil {
		return nil, errors.New("talhão não encontrado")
	}

	crop, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	crop.Name = req.Name
	crop.Type = req.Type
	crop.FieldID = req.FieldID

	err = s.Repo.Update(crop)
	if err != nil {
		return nil, err
	}

	return &dto.CropResponse{
		ID:      crop.ID,
		Name:    crop.Name,
		Type:    crop.Type,
		FieldID: crop.FieldID,
	}, nil
}

func (s *CropService) DeleteCrop(id uint) error {
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.Repo.Delete(id)
}
