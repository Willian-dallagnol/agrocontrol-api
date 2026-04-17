package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

type FieldRepository struct {
	DB *gorm.DB
}

func NewFieldRepository(db *gorm.DB) *FieldRepository {
	return &FieldRepository{DB: db}
}

func (r *FieldRepository) Create(field *entities.Field) error {
	return r.DB.Create(field).Error
}

func (r *FieldRepository) FindAll() ([]entities.Field, error) {
	var fields []entities.Field
	err := r.DB.Find(&fields).Error
	return fields, err
}

func (r *FieldRepository) FindByID(id uint) (*entities.Field, error) {
	var field entities.Field
	err := r.DB.First(&field, id).Error
	if err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *FieldRepository) Update(field *entities.Field) error {
	return r.DB.Save(field).Error
}

func (r *FieldRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Field{}, id).Error
}
func (r *FieldRepository) FindByFarmID(farmID uint) ([]entities.Field, error) {
	var fields []entities.Field

	err := r.DB.Where("farm_id = ?", farmID).Find(&fields).Error
	if err != nil {
		return nil, err
	}

	return fields, nil
}
