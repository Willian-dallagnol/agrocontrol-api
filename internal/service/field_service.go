package service

import (
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
	"agrocontrol-api/internal/repository"
	"errors"
)

// 🌱 Service responsável pela regra de negócio dos talhões (Field)
type FieldService struct {
	Repo *repository.FieldRepository
	// 👉 usado para acessar a tabela de fields no banco

	FarmRepo *repository.FarmRepository
	// 👉 usado para validar se a fazenda existe antes de criar/atualizar um field
}

// 🏗️ Construtor do service
func NewFieldService(repo *repository.FieldRepository, farmRepo *repository.FarmRepository) *FieldService {
	return &FieldService{
		Repo:     repo,
		FarmRepo: farmRepo,
	}
}

// 🚀 Criar novo talhão
func (s *FieldService) CreateField(req dto.CreateFieldRequest) (*dto.FieldResponse, error) {

	// 🔥 regra de negócio: área do talhão deve ser maior que zero
	if req.Area <= 0 {
		return nil, errors.New("a área do talhão deve ser maior que zero")
	}

	// 🔗 valida se a fazenda informada existe
	_, err := s.FarmRepo.FindByID(req.FarmID)
	if err != nil {
		return nil, errors.New("fazenda não encontrada")
	}

	// 🧩 monta a entidade Field com os dados recebidos
	field := &entities.Field{
		Name:     req.Name,
		Area:     req.Area,
		SoilType: req.SoilType,
		FarmID:   req.FarmID,
	}

	// 💾 salva no banco
	err = s.Repo.Create(field)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta formatada para a API
	return &dto.FieldResponse{
		ID:       field.ID,
		Name:     field.Name,
		Area:     field.Area,
		SoilType: field.SoilType,
		FarmID:   field.FarmID,
	}, nil
}

// 📋 Listar todos os talhões
func (s *FieldService) GetFields() ([]dto.FieldResponse, error) {

	// 🔍 busca todos no banco
	fields, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.FieldResponse

	// 🔄 converte entity → DTO
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

// 🔍 Buscar talhão por ID
func (s *FieldService) GetFieldByID(id uint) (*dto.FieldResponse, error) {

	// busca no banco
	field, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// retorna resposta formatada
	return &dto.FieldResponse{
		ID:       field.ID,
		Name:     field.Name,
		Area:     field.Area,
		SoilType: field.SoilType,
		FarmID:   field.FarmID,
	}, nil
}

// 🔄 Atualizar talhão
func (s *FieldService) UpdateField(id uint, req dto.UpdateFieldRequest) (*dto.FieldResponse, error) {

	// 🔥 regra de negócio: área deve ser maior que zero
	if req.Area <= 0 {
		return nil, errors.New("a área do talhão deve ser maior que zero")
	}

	// 🔗 valida novamente se a fazenda existe
	_, err := s.FarmRepo.FindByID(req.FarmID)
	if err != nil {
		return nil, errors.New("fazenda não encontrada")
	}

	// 🔍 busca o talhão atual no banco
	field, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// ✏️ atualiza os campos
	field.Name = req.Name
	field.Area = req.Area
	field.SoilType = req.SoilType
	field.FarmID = req.FarmID

	// 💾 salva atualização
	err = s.Repo.Update(field)
	if err != nil {
		return nil, err
	}

	// ✅ retorna resposta atualizada
	return &dto.FieldResponse{
		ID:       field.ID,
		Name:     field.Name,
		Area:     field.Area,
		SoilType: field.SoilType,
		FarmID:   field.FarmID,
	}, nil
}

// 🗑️ Deletar talhão
func (s *FieldService) DeleteField(id uint) error {

	// 🔍 primeiro verifica se o talhão existe
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	// 🗑️ remove do banco
	return s.Repo.Delete(id)
}

// 🔗 Listar todos os talhões de uma fazenda específica
func (s *FieldService) GetFieldsByFarmID(farmID uint) ([]dto.FieldResponse, error) {

	// 🔗 valida se a fazenda existe antes da consulta
	_, err := s.FarmRepo.FindByID(farmID)
	if err != nil {
		return nil, errors.New("fazenda não encontrada")
	}

	// 🔍 busca todos os fields relacionados à farm
	fields, err := s.Repo.FindByFarmID(farmID)
	if err != nil {
		return nil, err
	}

	var response []dto.FieldResponse

	// 🔄 converte entity → DTO
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
