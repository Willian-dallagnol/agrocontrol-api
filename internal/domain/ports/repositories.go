package ports

import (
	"agrocontrol-api/internal/domain/entities"
	"time"
)

// TxRunner abstrai a execução de transações.
// A implementação concreta fica em repository (usa *gorm.DB).
// Isso mantém gorm.DB fora dos serviços.
type TxRunner interface {
	RunInTx(fn func(tx TxRunner) error) error
}

// ── UserRepository ────────────────────────────────────────────────────────────

type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uint) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindAll() ([]entities.User, error)
	Update(user *entities.User) error
}

// ── FarmRepository ────────────────────────────────────────────────────────────

type FarmRepository interface {
	Create(farm *entities.Farm) error
	FindByID(id uint) (*entities.Farm, error)
	FindAllPaginated(offset, limit int, search string) ([]entities.Farm, int64, error)
	FindByCreatedByPaginated(userID uint, offset, limit int, search string) ([]entities.Farm, int64, error)
	FindAll() ([]entities.Farm, error)
	FindByCreatedBy(userID uint) ([]entities.Farm, error)
	Update(farm *entities.Farm) error
	Delete(id uint) error
	Count() (int64, error)
}

// ── FieldRepository ───────────────────────────────────────────────────────────

type FieldRepository interface {
	Create(field *entities.Field) error
	FindByID(id uint) (*entities.Field, error)
	FindAllPaginated(offset, limit int, search string) ([]entities.Field, int64, error)
	FindByUserPaginated(userID uint, offset, limit int, search string) ([]entities.Field, int64, error)
	FindAll() ([]entities.Field, error)
	FindByUser(userID uint) ([]entities.Field, error)
	FindByFarmID(farmID uint) ([]entities.Field, error)
	ExistsByNameAndFarm(name string, farmID uint, excludeID uint) (bool, error)
	Update(field *entities.Field) error
	Delete(id uint) error
	Count() (int64, error)
	BelongsToUser(fieldID, userID uint) (bool, error)
}

// ── CropRepository ────────────────────────────────────────────────────────────

type CropRepository interface {
	Create(crop *entities.Crop) error
	FindByID(id uint) (*entities.Crop, error)
	FindAllPaginated(offset, limit int, search string) ([]entities.Crop, int64, error)
	FindAll() ([]entities.Crop, error)
	Update(crop *entities.Crop) error
	Delete(id uint) error
}

// ── SeasonRepository ──────────────────────────────────────────────────────────

type SeasonRepository interface {
	Create(season *entities.Season) error
	FindByID(id uint) (*entities.Season, error)
	FindAllPaginated(offset, limit int, search string) ([]entities.Season, int64, error)
	FindAll() ([]entities.Season, error)
	Update(season *entities.Season) error
	Delete(id uint) error
}

// ── PlantingRepository ────────────────────────────────────────────────────────

type PlantingRepository interface {
	Create(planting *entities.Planting) error
	FindByID(id uint) (*entities.Planting, error)
	FindAll() ([]entities.Planting, error)
	FindByUser(userID uint) ([]entities.Planting, error)
	FindByFieldID(fieldID uint) ([]entities.Planting, error)
	HasActivePlanting(fieldID uint) (bool, error)
	CountActive() (int64, error)
	BelongsToUser(plantingID, userID uint) (bool, error)
	Update(planting *entities.Planting) error
	Delete(id uint) error
}

// ── InputRepository ───────────────────────────────────────────────────────────

type InputRepository interface {
	Create(input *entities.Input) error
	FindByID(id uint) (*entities.Input, error)
	FindAllPaginated(offset, limit int, search, category string) ([]entities.Input, int64, error)
	FindByUserPaginated(userID uint, offset, limit int, search, category string) ([]entities.Input, int64, error)
	FindAll() ([]entities.Input, error)
	FindByUser(userID uint) ([]entities.Input, error)
	FindLowStock() ([]entities.Input, error)
	FindExpiringSoon(days int) ([]entities.Input, error)
	Update(input *entities.Input) error
	Delete(id uint) error
	CountLowStock() (int64, error)
	// Métodos transacionais recebem TxRunner para manter gorm fora dos serviços
	CreateTx(tx TxRunner, input *entities.Input) error
	DeductStockTx(tx TxRunner, id uint, qty float64) error
	FindByIDTx(tx TxRunner, id uint) (*entities.Input, error)
}

// ── ApplicationRepository ─────────────────────────────────────────────────────

type ApplicationRepository interface {
	FindByID(id uint) (*entities.Application, error)
	FindAllPaginated(offset, limit int, fieldID uint, appType string) ([]entities.Application, int64, error)
	FindByUserPaginated(userID uint, offset, limit int, fieldID uint, appType string) ([]entities.Application, int64, error)
	FindAll() ([]entities.Application, error)
	FindByUser(userID uint) ([]entities.Application, error)
	FindByFieldID(fieldID uint) ([]entities.Application, error)
	Update(app *entities.Application) error
	Delete(id uint) error
	CountThisMonth() (int64, error)
	CreateTx(tx TxRunner, app *entities.Application) error
}

// ── MonitoringRepository ──────────────────────────────────────────────────────

type MonitoringRepository interface {
	Create(monitoring *entities.Monitoring) error
	FindByID(id uint) (*entities.Monitoring, error)
	FindAll() ([]entities.Monitoring, error)
	FindByUser(userID uint) ([]entities.Monitoring, error)
	FindByFieldID(fieldID uint) ([]entities.Monitoring, error)
	Update(monitoring *entities.Monitoring) error
	Delete(id uint) error
}

// ── HarvestRepository ─────────────────────────────────────────────────────────

type HarvestRepository interface {
	FindByID(id uint) (*entities.Harvest, error)
	FindAllPaginated(offset, limit int, fieldID uint) ([]entities.Harvest, int64, error)
	FindByUserPaginated(userID uint, offset, limit int, fieldID uint) ([]entities.Harvest, int64, error)
	FindAll() ([]entities.Harvest, error)
	FindByUser(userID uint) ([]entities.Harvest, error)
	ExistsByPlantingID(plantingID uint) (bool, error)
	Update(harvest *entities.Harvest) error
	Delete(id uint) error
	CreateTx(tx TxRunner, harvest *entities.Harvest) error
}

// ── AlertRepository ───────────────────────────────────────────────────────────

type AlertRepository interface {
	FindByID(id uint) (*entities.Alert, error)
	FindAll() ([]entities.Alert, error)
	FindByUser(userID uint) ([]entities.Alert, error)
	FindOpen() ([]entities.Alert, error)
	FindOpenByUser(userID uint) ([]entities.Alert, error)
	Update(alert *entities.Alert) error
	Delete(id uint) error
	CountOpen() (int64, error)
	CountOpenByUser(userID uint) (int64, error)
	CreateTx(tx TxRunner, alert *entities.Alert) error
	Create(alert *entities.Alert) error
}

// ── ReportRepository ──────────────────────────────────────────────────────────

type ReportRepository interface {
	FindProductivity(userID uint, role string, seasonID, farmID, cropID uint, offset, limit int) ([]ProductivityRow, int64, error)
	FindProductivitySummary(userID uint, role string, seasonID, farmID, cropID uint) (*ProductivitySummaryRow, error)
	FindCostPerField(userID uint, role string, fieldID, seasonID, farmID uint) ([]CostRow, error)
	FindHarvestBagsPerField(userID uint, role string, fieldIDs []uint) (map[uint]float64, error)
	FindOverview(userID uint, role string) (*OverviewRow, error)
	FindTopFields(userID uint, role string) ([]TopFieldRow, error)
}

// ── Tipos de resultado de query ───────────────────────────────────────────────

type ProductivityRow struct {
	PlantingID        uint
	FieldID           uint
	FieldName         string
	AreaHa            float64
	FarmID            uint
	FarmName          string
	SeasonID          uint
	SeasonName        string
	CropID            uint
	CropName          string
	Variety           string
	PlantingDate      time.Time
	HarvestDate       time.Time
	TotalBags         float64
	ProductivityBagHa float64
	ProductivityKgHa  float64
	GrainMoisture     float64
	FieldLoss         float64
}

type ProductivitySummaryRow struct {
	SeasonName           string
	TotalFields          int
	TotalPlantedAreaHa   float64
	TotalBags            float64
	AvgProductivityBagHa float64
	AvgProductivityKgHa  float64
	BestField            string
	BestFieldBagHa       float64
}

type CostRow struct {
	FieldID   uint
	FieldName string
	FarmName  string
	AreaHa    float64
	Category  string
	TotalCost float64
	TotalUsed float64
	AppCount  int64
}

type OverviewRow struct {
	TotalFarms            int64
	TotalFields           int64
	TotalAreaHa           float64
	PlantedAreaHa         float64
	TotalSeasons          int64
	ActivePlantings       int64
	HarvestedThisYear     int64
	AvgProductivityBagHa  float64
	TotalBagsThisYear     float64
	TotalInputTypes       int64
	LowStockInputs        int64
	ExpiringInputs        int64
	TotalApplications     int64
	ApplicationsThisMonth int64
	EstimatedTotalCost    float64
	OpenAlerts            int64
	CriticalAlerts        int64
}

type TopFieldRow struct {
	FieldName         string
	FarmName          string
	ProductivityBagHa float64
}
