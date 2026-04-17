package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

type CropRepository struct {
	DB *gorm.DB
}

func NewCropRepository(db *gorm.DB) *CropRepository {
	return &CropRepository{DB: db}
}

func (r *CropRepository) Create(crop *entities.Crop) error {
	return r.DB.Create(crop).Error
}

func (r *CropRepository) FindAll() ([]entities.Crop, error) {
	var crops []entities.Crop
	err := r.DB.Find(&crops).Error
	return crops, err
}

func (r *CropRepository) FindByID(id uint) (*entities.Crop, error) {
	var crop entities.Crop
	err := r.DB.First(&crop, id).Error
	if err != nil {
		return nil, err
	}
	return &crop, nil
}

func (r *CropRepository) Update(crop *entities.Crop) error {
	return r.DB.Save(crop).Error
}

func (r *CropRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Crop{}, id).Error
}
