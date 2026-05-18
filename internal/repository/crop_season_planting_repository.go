package repository

import (
	"agrocontrol-api/internal/domain/entities"

	"gorm.io/gorm"
)

// ---- CropRepository ----

type CropRepository struct{ DB *gorm.DB }

func NewCropRepository(db *gorm.DB) *CropRepository { return &CropRepository{DB: db} }

func (r *CropRepository) Create(c *entities.Crop) error  { return r.DB.Create(c).Error }
func (r *CropRepository) Update(c *entities.Crop) error  { return r.DB.Save(c).Error }
func (r *CropRepository) Delete(id uint) error           { return r.DB.Delete(&entities.Crop{}, id).Error }

func (r *CropRepository) FindAllPaginated(offset, limit int, search string) ([]entities.Crop, int64, error) {
	var crops []entities.Crop
	var total int64
	q := r.DB.Model(&entities.Crop{})
	if search != "" {
		q = q.Where("name ILIKE ? OR variety ILIKE ? OR type ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("name asc").Offset(offset).Limit(limit).Find(&crops).Error
	return crops, total, err
}

func (r *CropRepository) FindAll() ([]entities.Crop, error) {
	var crops []entities.Crop
	err := r.DB.Order("name asc").Find(&crops).Error
	return crops, err
}

func (r *CropRepository) FindByID(id uint) (*entities.Crop, error) {
	var crop entities.Crop
	if err := r.DB.First(&crop, id).Error; err != nil {
		return nil, err
	}
	return &crop, nil
}

// ---- SeasonRepository ----

type SeasonRepository struct{ DB *gorm.DB }

func NewSeasonRepository(db *gorm.DB) *SeasonRepository { return &SeasonRepository{DB: db} }

func (r *SeasonRepository) Create(s *entities.Season) error { return r.DB.Create(s).Error }
func (r *SeasonRepository) Update(s *entities.Season) error { return r.DB.Save(s).Error }
func (r *SeasonRepository) Delete(id uint) error            { return r.DB.Delete(&entities.Season{}, id).Error }

func (r *SeasonRepository) FindAllPaginated(offset, limit int, search string) ([]entities.Season, int64, error) {
	var seasons []entities.Season
	var total int64
	q := r.DB.Model(&entities.Season{})
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("start_date desc").Offset(offset).Limit(limit).Find(&seasons).Error
	return seasons, total, err
}

func (r *SeasonRepository) FindAll() ([]entities.Season, error) {
	var seasons []entities.Season
	err := r.DB.Order("start_date desc").Find(&seasons).Error
	return seasons, err
}

func (r *SeasonRepository) FindByID(id uint) (*entities.Season, error) {
	var season entities.Season
	if err := r.DB.First(&season, id).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

// ---- PlantingRepository ----

type PlantingRepository struct{ DB *gorm.DB }

func NewPlantingRepository(db *gorm.DB) *PlantingRepository { return &PlantingRepository{DB: db} }

func (r *PlantingRepository) Create(p *entities.Planting) error { return r.DB.Create(p).Error }
func (r *PlantingRepository) Update(p *entities.Planting) error { return r.DB.Save(p).Error }
func (r *PlantingRepository) Delete(id uint) error              { return r.DB.Delete(&entities.Planting{}, id).Error }

func (r *PlantingRepository) FindAll() ([]entities.Planting, error) {
	var plantings []entities.Planting
	err := r.DB.Preload("Field").Preload("Season").Preload("Crop").
		Order("planting_date desc").Find(&plantings).Error
	return plantings, err
}

func (r *PlantingRepository) FindByUser(userID uint) ([]entities.Planting, error) {
	var plantings []entities.Planting
	err := r.DB.
		Joins("JOIN fields ON fields.id = plantings.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID).
		Preload("Field").Preload("Season").Preload("Crop").
		Order("plantings.planting_date desc").
		Find(&plantings).Error
	return plantings, err
}

func (r *PlantingRepository) FindByID(id uint) (*entities.Planting, error) {
	var p entities.Planting
	if err := r.DB.Preload("Field").Preload("Season").Preload("Crop").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PlantingRepository) FindByFieldID(fieldID uint) ([]entities.Planting, error) {
	var plantings []entities.Planting
	err := r.DB.Where("field_id = ?", fieldID).
		Preload("Crop").Preload("Season").
		Order("planting_date desc").Find(&plantings).Error
	return plantings, err
}

func (r *PlantingRepository) HasActivePlanting(fieldID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Planting{}).
		Where("field_id = ? AND status = ?", fieldID, entities.PlantingStatusActive).
		Count(&count).Error
	return count > 0, err
}

func (r *PlantingRepository) CountActive() (int64, error) {
	var count int64
	err := r.DB.Model(&entities.Planting{}).Where("status = ?", entities.PlantingStatusActive).Count(&count).Error
	return count, err
}

func (r *PlantingRepository) BelongsToUser(plantingID, userID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Planting{}).
		Joins("JOIN fields ON fields.id = plantings.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("plantings.id = ? AND farms.created_by = ?", plantingID, userID).
		Count(&count).Error
	return count > 0, err
}
