package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

// 🚜 Repository responsável por acessar a tabela de Farm (fazendas)
type FarmRepository struct {
	DB *gorm.DB
	// 👉 conexão com o banco de dados
}

// 🏗️ Construtor do repository
func NewFarmRepository(db *gorm.DB) *FarmRepository {
	return &FarmRepository{DB: db}
}

// 🚀 Criar nova fazenda
func (r *FarmRepository) Create(farm *entities.Farm) error {
	// 👉 insere um novo registro na tabela farms
	return r.DB.Create(farm).Error
}

// 📋 Buscar todas as fazendas
func (r *FarmRepository) FindAll() ([]entities.Farm, error) {
	var farms []entities.Farm

	// 👉 SELECT * FROM farms
	err := r.DB.Find(&farms).Error

	return farms, err
}

// 🔍 Buscar fazenda por ID
func (r *FarmRepository) FindByID(id uint) (*entities.Farm, error) {
	var farm entities.Farm

	// 👉 SELECT * FROM farms WHERE id = ?
	err := r.DB.First(&farm, id).Error
	if err != nil {
		// ❌ retorna erro se não encontrar
		return nil, err
	}

	return &farm, nil
}

// 🔄 Atualizar fazenda
func (r *FarmRepository) Update(farm *entities.Farm) error {
	// 👉 atualiza o registro no banco
	return r.DB.Save(farm).Error
}

// 🗑️ Deletar fazenda
func (r *FarmRepository) Delete(id uint) error {
	// 👉 DELETE FROM farms WHERE id = ?
	return r.DB.Delete(&entities.Farm{}, id).Error
}
