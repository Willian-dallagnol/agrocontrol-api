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

// FindAllPaginated retorna fazendas paginadas com total (admin)
func (r *FarmRepository) FindAllPaginated(offset, limit int, search string) ([]entities.Farm, int64, error) {
	var farms []entities.Farm
	var total int64

	q := r.DB.Model(&entities.Farm{})
	if search != "" {
		q = q.Where("name ILIKE ? OR owner_name ILIKE ? OR city ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("name asc").Offset(offset).Limit(limit).Find(&farms).Error
	return farms, total, err
}

// FindByCreatedByPaginated retorna fazendas do usuário paginadas
func (r *FarmRepository) FindByCreatedByPaginated(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error) {
	var farms []entities.Farm
	var total int64

	q := r.DB.Model(&entities.Farm{}).Where("created_by = ?", userID)
	if search != "" {
		q = q.Where("name ILIKE ? OR owner_name ILIKE ? OR city ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("name asc").Offset(offset).Limit(limit).Find(&farms).Error
	return farms, total, err
}

// Mantém FindAll e FindByCreatedBy sem paginação para uso interno (seed, dashboard)
func (r *FarmRepository) FindAll() ([]entities.Farm, error) {
	var farms []entities.Farm
	err := r.DB.Order("name asc").Find(&farms).Error
	return farms, err
}

func (r *FarmRepository) FindByCreatedBy(userID uint) ([]entities.Farm, error) {
	var farms []entities.Farm
	err := r.DB.Where("created_by = ?", userID).Order("name asc").Find(&farms).Error
	return farms, err
}

func (r *FarmRepository) FindByID(id uint) (*entities.Farm, error) {
	var farm entities.Farm
	if err := r.DB.First(&farm, id).Error; err != nil {
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

func (r *FarmRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&entities.Farm{}).Count(&count).Error
	return count, err
}
