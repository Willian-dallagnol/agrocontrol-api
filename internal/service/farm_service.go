package service

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"errors"
)

type FarmService struct {
	Repo *repository.FarmRepository
}

func NewFarmService(repo *repository.FarmRepository) *FarmService {
	return &FarmService{Repo: repo}
}

func (s *FarmService) CreateFarm(req dto.CreateFarmRequest, userID uint) (*dto.FarmResponse, error) {
	if req.TotalArea <= 0 {
		return nil, errors.New("a área total da fazenda deve ser maior que zero")
	}

	farm := &entities.Farm{
		Name:      req.Name,
		OwnerName: req.OwnerName,
		Location:  req.Location,
		TotalArea: req.TotalArea,
		City:      req.City,
		State:     req.State,
		CreatedBy: userID,
	}

	err := s.Repo.Create(farm)
	if err != nil {
		return nil, err
	}

	return &dto.FarmResponse{
		ID:        farm.ID,
		Name:      farm.Name,
		OwnerName: farm.OwnerName,
		Location:  farm.Location,
		TotalArea: farm.TotalArea,
		City:      farm.City,
		State:     farm.State,
	}, nil
}

func (s *FarmService) GetFarms() ([]dto.FarmResponse, error) {
	farms, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.FarmResponse

	for _, farm := range farms {
		response = append(response, dto.FarmResponse{
			ID:        farm.ID,
			Name:      farm.Name,
			OwnerName: farm.OwnerName,
			Location:  farm.Location,
			TotalArea: farm.TotalArea,
			City:      farm.City,
			State:     farm.State,
		})
	}

	return response, nil
}

func (s *FarmService) GetFarmByID(id uint) (*dto.FarmResponse, error) {
	farm, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.FarmResponse{
		ID:        farm.ID,
		Name:      farm.Name,
		OwnerName: farm.OwnerName,
		Location:  farm.Location,
		TotalArea: farm.TotalArea,
		City:      farm.City,
		State:     farm.State,
	}, nil
}

func (s *FarmService) UpdateFarm(id uint, req dto.UpdateFarmRequest) (*dto.FarmResponse, error) {
	if req.TotalArea <= 0 {
		return nil, errors.New("a área total da fazenda deve ser maior que zero")
	}

	farm, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	farm.Name = req.Name
	farm.OwnerName = req.OwnerName
	farm.Location = req.Location
	farm.TotalArea = req.TotalArea
	farm.City = req.City
	farm.State = req.State

	err = s.Repo.Update(farm)
	if err != nil {
		return nil, err
	}

	return &dto.FarmResponse{
		ID:        farm.ID,
		Name:      farm.Name,
		OwnerName: farm.OwnerName,
		Location:  farm.Location,
		TotalArea: farm.TotalArea,
		City:      farm.City,
		State:     farm.State,
	}, nil
}

func (s *FarmService) DeleteFarm(id uint) error {
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.Repo.Delete(id)
}
