package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

// 🌱 Repository responsável por acessar a tabela de Field (talhões)
type FieldRepository struct {
	DB *gorm.DB
	// 👉 conexão com o banco de dados
}

// 🏗️ Construtor do repository
func NewFieldRepository(db *gorm.DB) *FieldRepository {
	return &FieldRepository{DB: db}
}

// 🚀 Criar novo talhão
func (r *FieldRepository) Create(field *entities.Field) error {
	// 👉 insere um novo registro na tabela fields
	return r.DB.Create(field).Error
}

// 📋 Buscar todos os talhões
func (r *FieldRepository) FindAll() ([]entities.Field, error) {
	var fields []entities.Field

	// 👉 SELECT * FROM fields
	err := r.DB.Find(&fields).Error

	return fields, err
}

// 🔍 Buscar talhão por ID
func (r *FieldRepository) FindByID(id uint) (*entities.Field, error) {
	var field entities.Field

	// 👉 SELECT * FROM fields WHERE id = ?
	err := r.DB.First(&field, id).Error
	if err != nil {
		// ❌ retorna erro se não encontrar
		return nil, err
	}

	return &field, nil
}

// 🔄 Atualizar talhão
func (r *FieldRepository) Update(field *entities.Field) error {
	// 👉 atualiza o registro no banco
	return r.DB.Save(field).Error
}

// 🗑️ Deletar talhão
func (r *FieldRepository) Delete(id uint) error {
	// 👉 DELETE FROM fields WHERE id = ?
	return r.DB.Delete(&entities.Field{}, id).Error
}

// 🔗 Buscar todos os talhões de uma fazenda específica
func (r *FieldRepository) FindByFarmID(farmID uint) ([]entities.Field, error) {
	var fields []entities.Field

	// 👉 SELECT * FROM fields WHERE farm_id = ?
	err := r.DB.Where("farm_id = ?", farmID).Find(&fields).Error
	if err != nil {
		return nil, err
	}

	return fields, nil
}
