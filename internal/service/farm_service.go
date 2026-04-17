package service

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"errors"
)

// 🚜 Service responsável pela regra de negócio das fazendas
type FarmService struct {
	Repo *repository.FarmRepository
	// 👉 responsável por acessar o banco
}

// 🏗️ Construtor do service
func NewFarmService(repo *repository.FarmRepository) *FarmService {
	return &FarmService{Repo: repo}
}

// 🚀 Criar nova fazenda
func (s *FarmService) CreateFarm(req dto.CreateFarmRequest, userID uint) (*dto.FarmResponse, error) {

	// 🔥 regra de negócio: área deve ser maior que zero
	if req.TotalArea <= 0 {
		return nil, errors.New("a área total da fazenda deve ser maior que zero")
	}

	// 🧩 monta a entidade com os dados recebidos
	farm := &entities.Farm{
		Name:      req.Name,
		OwnerName: req.OwnerName,
		Location:  req.Location,
		TotalArea: req.TotalArea,
		City:      req.City,
		State:     req.State,
		CreatedBy: userID, // 🔐 vínculo com usuário logado
	}

	// 💾 salva no banco
	err := s.Repo.Create(farm)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta para API
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

// 📋 Listar todas as fazendas
func (s *FarmService) GetFarms() ([]dto.FarmResponse, error) {

	// 🔍 busca no banco
	farms, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.FarmResponse

	// 🔄 converte entity → DTO
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

// 🔍 Buscar fazenda por ID
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

// 🔄 Atualizar fazenda
func (s *FarmService) UpdateFarm(id uint, req dto.UpdateFarmRequest) (*dto.FarmResponse, error) {

	// 🔥 regra de negócio
	if req.TotalArea <= 0 {
		return nil, errors.New("a área total da fazenda deve ser maior que zero")
	}

	// 🔍 busca no banco
	farm, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// ✏️ atualiza campos
	farm.Name = req.Name
	farm.OwnerName = req.OwnerName
	farm.Location = req.Location
	farm.TotalArea = req.TotalArea
	farm.City = req.City
	farm.State = req.State

	// 💾 salva no banco
	err = s.Repo.Update(farm)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta atualizada
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

// 🗑️ Deletar fazenda
func (s *FarmService) DeleteFarm(id uint) error {

	// 🔍 verifica se existe
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	// 🗑️ remove do banco
	return s.Repo.Delete(id)
}
