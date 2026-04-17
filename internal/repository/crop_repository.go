package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

// 🌾 Repository responsável por acessar a tabela de Crop no banco
type CropRepository struct {
	DB *gorm.DB
	// 👉 conexão com o banco de dados
}

// 🏗️ Construtor do repository
func NewCropRepository(db *gorm.DB) *CropRepository {
	return &CropRepository{DB: db}
}

// 🚀 Criar nova cultura no banco
func (r *CropRepository) Create(crop *entities.Crop) error {
	// 👉 insere o registro na tabela crops
	return r.DB.Create(crop).Error
}

// 📋 Buscar todas as culturas
func (r *CropRepository) FindAll() ([]entities.Crop, error) {
	var crops []entities.Crop

	// 👉 SELECT * FROM crops
	err := r.DB.Find(&crops).Error

	return crops, err
}

// 🔍 Buscar cultura por ID
func (r *CropRepository) FindByID(id uint) (*entities.Crop, error) {
	var crop entities.Crop

	// 👉 SELECT * FROM crops WHERE id = ?
	err := r.DB.First(&crop, id).Error
	if err != nil {
		// ❌ retorna erro se não encontrar
		return nil, err
	}

	return &crop, nil
}

// 🔄 Atualizar cultura
func (r *CropRepository) Update(crop *entities.Crop) error {
	// 👉 atualiza o registro no banco
	return r.DB.Save(crop).Error
}

// 🗑️ Deletar cultura
func (r *CropRepository) Delete(id uint) error {
	// 👉 DELETE FROM crops WHERE id = ?
	return r.DB.Delete(&entities.Crop{}, id).Error
}
