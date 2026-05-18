package repository

import (
	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
	"time"

	"gorm.io/gorm"
)

// dbFromTx extrai o *gorm.DB de um ports.TxRunner.
// Funciona porque GormTxRunner é o único implementador de TxRunner.
func dbFromTx(tx ports.TxRunner) *gorm.DB {
	return tx.(*GormTxRunner).DB()
}

// ---- InputRepository ----

type InputRepository struct{ DB *gorm.DB }

func NewInputRepository(db *gorm.DB) *InputRepository { return &InputRepository{DB: db} }

func (r *InputRepository) Create(i *entities.Input) error { return r.DB.Create(i).Error }
func (r *InputRepository) Update(i *entities.Input) error { return r.DB.Save(i).Error }
func (r *InputRepository) Delete(id uint) error           { return r.DB.Delete(&entities.Input{}, id).Error }

func (r *InputRepository) FindAllPaginated(offset, limit int, search, category string) ([]entities.Input, int64, error) {
	var inputs []entities.Input
	var total int64
	q := r.DB.Model(&entities.Input{}).Where("active = true")
	if search != "" {
		q = q.Where("name ILIKE ? OR manufacturer ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("name asc").Offset(offset).Limit(limit).Find(&inputs).Error
	return inputs, total, err
}

func (r *InputRepository) FindByUserPaginated(userID uint, offset, limit int, search, category string) ([]entities.Input, int64, error) {
	var inputs []entities.Input
	var total int64
	q := r.DB.Model(&entities.Input{}).Where("active = true AND created_by = ?", userID)
	if search != "" {
		q = q.Where("name ILIKE ? OR manufacturer ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("name asc").Offset(offset).Limit(limit).Find(&inputs).Error
	return inputs, total, err
}

func (r *InputRepository) FindAll() ([]entities.Input, error) {
	var inputs []entities.Input
	err := r.DB.Where("active = true").Order("name asc").Find(&inputs).Error
	return inputs, err
}

func (r *InputRepository) FindByUser(userID uint) ([]entities.Input, error) {
	var inputs []entities.Input
	err := r.DB.Where("active = true AND created_by = ?", userID).Order("name asc").Find(&inputs).Error
	return inputs, err
}

func (r *InputRepository) FindByID(id uint) (*entities.Input, error) {
	var input entities.Input
	if err := r.DB.First(&input, id).Error; err != nil {
		return nil, err
	}
	return &input, nil
}

func (r *InputRepository) FindByIDTx(tx ports.TxRunner, id uint) (*entities.Input, error) {
	var input entities.Input
	if err := dbFromTx(tx).First(&input, id).Error; err != nil {
		return nil, err
	}
	return &input, nil
}

func (r *InputRepository) CreateTx(tx ports.TxRunner, input *entities.Input) error {
	return dbFromTx(tx).Create(input).Error
}

func (r *InputRepository) DeductStockTx(tx ports.TxRunner, id uint, qty float64) error {
	db := dbFromTx(tx)
	result := db.Model(&entities.Input{}).
		Where("id = ? AND stock_qty >= ?", id, qty).
		UpdateColumn("stock_qty", gorm.Expr("stock_qty - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrInsufficientStock
	}
	return nil
}

func (r *InputRepository) CountLowStock() (int64, error) {
	var count int64
	err := r.DB.Model(&entities.Input{}).
		Where("active = true AND stock_qty <= min_stock_qty").
		Count(&count).Error
	return count, err
}

func (r *InputRepository) FindLowStock() ([]entities.Input, error) {
	var inputs []entities.Input
	err := r.DB.Where("active = true AND stock_qty <= min_stock_qty").Find(&inputs).Error
	return inputs, err
}

func (r *InputRepository) FindExpiringSoon(days int) ([]entities.Input, error) {
	var inputs []entities.Input
	deadline := time.Now().AddDate(0, 0, days)
	err := r.DB.Where("active = true AND expiration_date IS NOT NULL AND expiration_date <= ?", deadline).
		Find(&inputs).Error
	return inputs, err
}

// ---- ApplicationRepository ----

type ApplicationRepository struct{ DB *gorm.DB }

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{DB: db}
}

func (r *ApplicationRepository) CreateTx(tx ports.TxRunner, a *entities.Application) error {
	return dbFromTx(tx).Create(a).Error
}

func (r *ApplicationRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Application{}, id).Error
}

func (r *ApplicationRepository) FindAllPaginated(offset, limit int, fieldID uint, appType string) ([]entities.Application, int64, error) {
	var apps []entities.Application
	var total int64
	q := r.DB.Model(&entities.Application{})
	if fieldID > 0 {
		q = q.Where("field_id = ?", fieldID)
	}
	if appType != "" {
		q = q.Where("application_type = ?", appType)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Preload("Field").Preload("Input").
		Order("application_date desc").Offset(offset).Limit(limit).Find(&apps).Error
	return apps, total, err
}

func (r *ApplicationRepository) FindByUserPaginated(userID uint, offset, limit int, fieldID uint, appType string) ([]entities.Application, int64, error) {
	var apps []entities.Application
	var total int64
	q := r.DB.Model(&entities.Application{}).
		Joins("JOIN fields ON fields.id = applications.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID)
	if fieldID > 0 {
		q = q.Where("applications.field_id = ?", fieldID)
	}
	if appType != "" {
		q = q.Where("applications.application_type = ?", appType)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Preload("Field").Preload("Input").
		Order("applications.application_date desc").Offset(offset).Limit(limit).Find(&apps).Error
	return apps, total, err
}

func (r *ApplicationRepository) FindAll() ([]entities.Application, error) {
	var apps []entities.Application
	err := r.DB.Preload("Field").Preload("Input").
		Order("application_date desc").Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) FindByUser(userID uint) ([]entities.Application, error) {
	var apps []entities.Application
	err := r.DB.
		Joins("JOIN fields ON fields.id = applications.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID).
		Preload("Field").Preload("Input").
		Order("applications.application_date desc").
		Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) FindByID(id uint) (*entities.Application, error) {
	var app entities.Application
	if err := r.DB.Preload("Field").Preload("Input").First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) FindByFieldID(fieldID uint) ([]entities.Application, error) {
	var apps []entities.Application
	err := r.DB.Where("field_id = ?", fieldID).
		Preload("Input").
		Order("application_date desc").Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) Update(app *entities.Application) error {
	return r.DB.Save(app).Error
}

func (r *ApplicationRepository) CountThisMonth() (int64, error) {
	var count int64
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := r.DB.Model(&entities.Application{}).
		Where("application_date >= ?", start).
		Count(&count).Error
	return count, err
}

// ---- MonitoringRepository ----

type MonitoringRepository struct{ DB *gorm.DB }

func NewMonitoringRepository(db *gorm.DB) *MonitoringRepository {
	return &MonitoringRepository{DB: db}
}

func (r *MonitoringRepository) Create(m *entities.Monitoring) error { return r.DB.Create(m).Error }
func (r *MonitoringRepository) Update(m *entities.Monitoring) error { return r.DB.Save(m).Error }
func (r *MonitoringRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Monitoring{}, id).Error
}

func (r *MonitoringRepository) FindAll() ([]entities.Monitoring, error) {
	var mons []entities.Monitoring
	err := r.DB.Preload("Field").Order("inspection_date desc").Find(&mons).Error
	return mons, err
}

func (r *MonitoringRepository) FindByUser(userID uint) ([]entities.Monitoring, error) {
	var mons []entities.Monitoring
	err := r.DB.
		Joins("JOIN fields ON fields.id = monitorings.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID).
		Preload("Field").
		Order("monitorings.inspection_date desc").
		Find(&mons).Error
	return mons, err
}

func (r *MonitoringRepository) FindByID(id uint) (*entities.Monitoring, error) {
	var m entities.Monitoring
	if err := r.DB.Preload("Field").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MonitoringRepository) FindByFieldID(fieldID uint) ([]entities.Monitoring, error) {
	var mons []entities.Monitoring
	err := r.DB.Where("field_id = ?", fieldID).Order("inspection_date desc").Find(&mons).Error
	return mons, err
}

// ---- HarvestRepository ----

type HarvestRepository struct{ DB *gorm.DB }

func NewHarvestRepository(db *gorm.DB) *HarvestRepository { return &HarvestRepository{DB: db} }

func (r *HarvestRepository) CreateTx(tx ports.TxRunner, h *entities.Harvest) error {
	return dbFromTx(tx).Create(h).Error
}
func (r *HarvestRepository) Update(h *entities.Harvest) error { return r.DB.Save(h).Error }
func (r *HarvestRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Harvest{}, id).Error
}

func (r *HarvestRepository) FindAllPaginated(offset, limit int, fieldID uint) ([]entities.Harvest, int64, error) {
	var harvests []entities.Harvest
	var total int64
	q := r.DB.Model(&entities.Harvest{})
	if fieldID > 0 {
		q = q.Where("field_id = ?", fieldID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Preload("Planting").Preload("Field").
		Order("harvest_date desc").Offset(offset).Limit(limit).Find(&harvests).Error
	return harvests, total, err
}

func (r *HarvestRepository) FindByUserPaginated(userID uint, offset, limit int, fieldID uint) ([]entities.Harvest, int64, error) {
	var harvests []entities.Harvest
	var total int64
	q := r.DB.Model(&entities.Harvest{}).
		Joins("JOIN fields ON fields.id = harvests.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID)
	if fieldID > 0 {
		q = q.Where("harvests.field_id = ?", fieldID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Preload("Planting").Preload("Field").
		Order("harvests.harvest_date desc").Offset(offset).Limit(limit).Find(&harvests).Error
	return harvests, total, err
}

func (r *HarvestRepository) FindAll() ([]entities.Harvest, error) {
	var harvests []entities.Harvest
	err := r.DB.Preload("Planting").Preload("Field").
		Order("harvest_date desc").Find(&harvests).Error
	return harvests, err
}

func (r *HarvestRepository) FindByUser(userID uint) ([]entities.Harvest, error) {
	var harvests []entities.Harvest
	err := r.DB.
		Joins("JOIN fields ON fields.id = harvests.field_id").
		Joins("JOIN farms ON farms.id = fields.farm_id").
		Where("farms.created_by = ?", userID).
		Preload("Planting").Preload("Field").
		Order("harvests.harvest_date desc").
		Find(&harvests).Error
	return harvests, err
}

func (r *HarvestRepository) FindByID(id uint) (*entities.Harvest, error) {
	var h entities.Harvest
	if err := r.DB.Preload("Planting").Preload("Field").First(&h, id).Error; err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HarvestRepository) ExistsByPlantingID(plantingID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Harvest{}).Where("planting_id = ?", plantingID).Count(&count).Error
	return count > 0, err
}

// ---- AlertRepository ----

type AlertRepository struct{ DB *gorm.DB }

func NewAlertRepository(db *gorm.DB) *AlertRepository { return &AlertRepository{DB: db} }

func (r *AlertRepository) Create(a *entities.Alert) error { return r.DB.Create(a).Error }
func (r *AlertRepository) CreateTx(tx ports.TxRunner, a *entities.Alert) error {
	return dbFromTx(tx).Create(a).Error
}
func (r *AlertRepository) Update(a *entities.Alert) error { return r.DB.Save(a).Error }
func (r *AlertRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Alert{}, id).Error
}

func (r *AlertRepository) FindAll() ([]entities.Alert, error) {
	var alerts []entities.Alert
	err := r.DB.Order("created_at desc").Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) FindByUser(userID uint) ([]entities.Alert, error) {
	var alerts []entities.Alert
	err := r.DB.Where("created_by = ?", userID).Order("created_at desc").Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) FindOpen() ([]entities.Alert, error) {
	var alerts []entities.Alert
	err := r.DB.Where("status = ?", entities.AlertStatusOpen).
		Order("priority desc, created_at desc").
		Limit(20).Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) FindOpenByUser(userID uint) ([]entities.Alert, error) {
	var alerts []entities.Alert
	err := r.DB.Where("status = ? AND created_by = ?", entities.AlertStatusOpen, userID).
		Order("priority desc, created_at desc").
		Limit(20).Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) FindByID(id uint) (*entities.Alert, error) {
	var alert entities.Alert
	if err := r.DB.First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertRepository) CountOpen() (int64, error) {
	var count int64
	err := r.DB.Model(&entities.Alert{}).Where("status = ?", entities.AlertStatusOpen).Count(&count).Error
	return count, err
}

func (r *AlertRepository) CountOpenByUser(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&entities.Alert{}).
		Where("status = ? AND created_by = ?", entities.AlertStatusOpen, userID).
		Count(&count).Error
	return count, err
}
