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

func (r *FieldRepository) FindAllPaginated(offset, limit int, search string) ([]entities.Field, int64, error) {
	var fields []entities.Field
	var total int64
	q := r.DB.Model(&entities.Field{})
	if search != "" {
		q = q.Where("name ILIKE ? OR soil_type ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("name asc").Offset(offset).Limit(limit).Find(&fields).Error
	return fields, total, err
}

func (r *FieldRepository) FindByUserPaginated(userID uint, offset, limit int, search string) ([]entities.Field, int64, error) {
	var fields []entities.Field
	var total int64
	q := r.DB.Model(&entities.Field{}).
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID)
	if search != "" {
		q = q.Where("fields.name ILIKE ? OR fields.soil_type ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("fields.name asc").Offset(offset).Limit(limit).Find(&fields).Error
	return fields, total, err
}

func (r *FieldRepository) FindAll() ([]entities.Field, error) {
	var fields []entities.Field
	err := r.DB.Order("name asc").Find(&fields).Error
	return fields, err
}

func (r *FieldRepository) FindByUser(userID uint) ([]entities.Field, error) {
	var fields []entities.Field
	err := r.DB.
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID).
		Order("fields.name asc").Find(&fields).Error
	return fields, err
}

func (r *FieldRepository) FindByID(id uint) (*entities.Field, error) {
	var field entities.Field
	if err := r.DB.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *FieldRepository) FindByFarmID(farmID uint) ([]entities.Field, error) {
	var fields []entities.Field
	err := r.DB.Where("farm_id = ?", farmID).Order("name asc").Find(&fields).Error
	return fields, err
}

func (r *FieldRepository) ExistsByNameAndFarm(name string, farmID uint, excludeID uint) (bool, error) {
	var count int64
	q := r.DB.Model(&entities.Field{}).Where("name = ? AND farm_id = ?", name, farmID)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *FieldRepository) Update(field *entities.Field) error {
	return r.DB.Save(field).Error
}

func (r *FieldRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Field{}, id).Error
}

func (r *FieldRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&entities.Field{}).Count(&count).Error
	return count, err
}

func (r *FieldRepository) BelongsToUser(fieldID, userID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Field{}).
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("fields.id = ? AND farms.created_by = ?", fieldID, userID).
		Count(&count).Error
	return count > 0, err
}
