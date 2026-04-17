package service

import (
	"errors"

	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
)

// 🌾 Service responsável pela regra de negócio do módulo Crop
type CropService struct {
	Repo *repository.CropRepository
	// 👉 usado para salvar, buscar, atualizar e deletar culturas no banco

	FieldRepo *repository.FieldRepository
	// 👉 usado para validar se o Field existe antes de criar/atualizar Crop
}

// 🏗️ Construtor do service
func NewCropService(repo *repository.CropRepository, fieldRepo *repository.FieldRepository) *CropService {
	return &CropService{
		Repo:      repo,
		FieldRepo: fieldRepo,
	}
}

// 🚀 Criar nova cultura
func (s *CropService) CreateCrop(req dto.CreateCropRequest) (*dto.CropResponse, error) {

	// 🔍 valida se o talhão (Field) existe
	_, err := s.FieldRepo.FindByID(req.FieldID)
	if err != nil {
		// ❌ não permite criar cultura sem talhão válido
		return nil, errors.New("talhão não encontrado")
	}

	// 🧩 monta a entidade Crop com os dados recebidos
	crop := entities.Crop{
		Name:    req.Name,
		Type:    req.Type,
		FieldID: req.FieldID,
	}

	// 💾 salva no banco
	err = s.Repo.Create(&crop)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta formatada para a API
	return &dto.CropResponse{
		ID:      crop.ID,
		Name:    crop.Name,
		Type:    crop.Type,
		FieldID: crop.FieldID,
	}, nil
}

// 📋 Listar todas as culturas
func (s *CropService) GetCrops() ([]dto.CropResponse, error) {
	// 🔍 busca tudo no banco
	crops, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.CropResponse

	// 🔄 converte entidades para DTO de resposta
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

// 🔍 Buscar cultura por ID
func (s *CropService) GetCropByID(id uint) (*dto.CropResponse, error) {
	// busca no banco
	crop, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// retorna resposta formatada
	return &dto.CropResponse{
		ID:      crop.ID,
		Name:    crop.Name,
		Type:    crop.Type,
		FieldID: crop.FieldID,
	}, nil
}

// 🔄 Atualizar cultura
func (s *CropService) UpdateCrop(id uint, req dto.UpdateCropRequest) (*dto.CropResponse, error) {

	// 🔍 valida novamente se o Field existe
	_, err := s.FieldRepo.FindByID(req.FieldID)
	if err != nil {
		return nil, errors.New("talhão não encontrado")
	}

	// 🔍 busca a cultura atual no banco
	crop, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// ✏️ atualiza os campos
	crop.Name = req.Name
	crop.Type = req.Type
	crop.FieldID = req.FieldID

	// 💾 salva atualização no banco
	err = s.Repo.Update(crop)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta atualizada
	return &dto.CropResponse{
		ID:      crop.ID,
		Name:    crop.Name,
		Type:    crop.Type,
		FieldID: crop.FieldID,
	}, nil
}

// 🗑️ Deletar cultura
func (s *CropService) DeleteCrop(id uint) error {
	// 🔍 primeiro verifica se a cultura existe
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	// 🗑️ deleta do banco
	return s.Repo.Delete(id)
}
