package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

type FarmRepository struct {
	DB *gorm.DB
}

func NewFarmRepository(db *gorm.DB) *FarmRepository {
	return &FarmRepository{DB: db}
}

func (r *FarmRepository) Create(farm *entities.Farm) error {
	return r.DB.Create(farm).Error
}

func (r *FarmRepository) FindAll() ([]entities.Farm, error) {
	var farms []entities.Farm
	err := r.DB.Find(&farms).Error
	return farms, err
}

func (r *FarmRepository) FindByID(id uint) (*entities.Farm, error) {
	var farm entities.Farm
	err := r.DB.First(&farm, id).Error
	if err != nil {
		return nil, err
	}
	return &farm, nil
}

func (r *FarmRepository) Update(farm *entities.Farm) error {
	return r.DB.Save(farm).Error
}

func (r *FarmRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Farm{}, id).Error
}
