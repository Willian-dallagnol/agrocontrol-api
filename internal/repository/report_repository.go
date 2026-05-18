package repository

import (
	"agrocontrol-api/internal/domain/ports"
	"time"

	"gorm.io/gorm"
)

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) FindProductivity(
	userID uint, role string,
	seasonID, farmID, cropID uint,
	offset, limit int,
) ([]ports.ProductivityRow, int64, error) {

	base := r.DB.Table("plantings p").
		Select(`
			p.id                    AS planting_id,
			fi.id                   AS field_id,
			fi.name                 AS field_name,
			fi.area                 AS area_ha,
			fa.id                   AS farm_id,
			fa.name                 AS farm_name,
			s.id                    AS season_id,
			s.name                  AS season_name,
			cr.id                   AS crop_id,
			cr.name                 AS crop_name,
			cr.variety              AS variety,
			p.planting_date         AS planting_date,
			h.harvest_date          AS harvest_date,
			COALESCE(h.total_bags, 0)            AS total_bags,
			COALESCE(h.productivity_bag_ha, 0)   AS productivity_bag_ha,
			COALESCE(h.productivity_kg_ha, 0)    AS productivity_kg_ha,
			COALESCE(h.grain_moisture, 0)        AS grain_moisture,
			COALESCE(h.field_loss, 0)            AS field_loss
		`).
		Joins("JOIN fields fi ON fi.id = p.field_id").
		Joins("JOIN farms fa  ON fa.id = fi.farm_id").
		Joins("JOIN seasons s  ON s.id = p.season_id").
		Joins("JOIN crops cr   ON cr.id = p.crop_id").
		Joins("LEFT JOIN harvests h ON h.planting_id = p.id")

	if role != "admin" {
		base = base.Where("fa.created_by = ?", userID)
	}
	if seasonID > 0 { base = base.Where("p.season_id = ?", seasonID) }
	if farmID > 0   { base = base.Where("fa.id = ?", farmID) }
	if cropID > 0   { base = base.Where("p.crop_id = ?", cropID) }

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []ports.ProductivityRow
	err := base.
		Order("h.productivity_bag_ha DESC NULLS LAST, fi.name ASC").
		Offset(offset).Limit(limit).
		Scan(&rows).Error

	return rows, total, err
}

func (r *ReportRepository) FindProductivitySummary(
	userID uint, role string,
	seasonID, farmID, cropID uint,
) (*ports.ProductivitySummaryRow, error) {

	type Row struct {
		TotalFields          int64
		TotalPlantedAreaHa   float64
		TotalBags            float64
		AvgProductivityBagHa float64
		AvgProductivityKgHa  float64
		SeasonName           string
	}

	base := r.DB.Table("plantings p").
		Select(`
			COUNT(DISTINCT fi.id)                        AS total_fields,
			COALESCE(SUM(fi.area), 0)                    AS total_planted_area_ha,
			COALESCE(SUM(h.total_bags), 0)               AS total_bags,
			COALESCE(AVG(h.productivity_bag_ha), 0)      AS avg_productivity_bag_ha,
			COALESCE(AVG(h.productivity_kg_ha), 0)       AS avg_productivity_kg_ha,
			MAX(s.name)                                  AS season_name
		`).
		Joins("JOIN fields fi ON fi.id = p.field_id").
		Joins("JOIN farms fa  ON fa.id = fi.farm_id").
		Joins("JOIN seasons s  ON s.id = p.season_id").
		Joins("JOIN crops cr   ON cr.id = p.crop_id").
		Joins("LEFT JOIN harvests h ON h.planting_id = p.id")

	if role != "admin" { base = base.Where("fa.created_by = ?", userID) }
	if seasonID > 0    { base = base.Where("p.season_id = ?", seasonID) }
	if farmID > 0      { base = base.Where("fa.id = ?", farmID) }
	if cropID > 0      { base = base.Where("p.crop_id = ?", cropID) }

	var row Row
	if err := base.Scan(&row).Error; err != nil {
		return nil, err
	}

	type BestRow struct {
		FieldName         string
		ProductivityBagHa float64
	}
	bestQ := r.DB.Table("harvests h").
		Select("fi.name AS field_name, h.productivity_bag_ha").
		Joins("JOIN plantings p ON p.id = h.planting_id").
		Joins("JOIN fields fi ON fi.id = p.field_id").
		Joins("JOIN farms fa ON fa.id = fi.farm_id")
	if role != "admin" { bestQ = bestQ.Where("fa.created_by = ?", userID) }
	if seasonID > 0    { bestQ = bestQ.Where("p.season_id = ?", seasonID) }
	var best BestRow
	bestQ.Order("h.productivity_bag_ha DESC").Limit(1).Scan(&best)

	return &ports.ProductivitySummaryRow{
		SeasonName:           row.SeasonName,
		TotalFields:          int(row.TotalFields),
		TotalPlantedAreaHa:   row.TotalPlantedAreaHa,
		TotalBags:            row.TotalBags,
		AvgProductivityBagHa: row.AvgProductivityBagHa,
		AvgProductivityKgHa:  row.AvgProductivityKgHa,
		BestField:            best.FieldName,
		BestFieldBagHa:       best.ProductivityBagHa,
	}, nil
}

func (r *ReportRepository) FindCostPerField(
	userID uint, role string,
	fieldID, seasonID, farmID uint,
) ([]ports.CostRow, error) {

	base := r.DB.Table("applications a").
		Select(`
			fi.id                                                   AS field_id,
			fi.name                                                 AS field_name,
			fa.name                                                 AS farm_name,
			fi.area                                                 AS area_ha,
			i.category                                              AS category,
			COALESCE(SUM(a.total_used * i.cost_per_unit), 0)        AS total_cost,
			COALESCE(SUM(a.total_used), 0)                          AS total_used,
			COUNT(a.id)                                             AS app_count
		`).
		Joins("JOIN fields fi  ON fi.id = a.field_id").
		Joins("JOIN farms fa   ON fa.id = fi.farm_id").
		Joins("JOIN inputs i   ON i.id  = a.input_id").
		Group("fi.id, fi.name, fa.name, fi.area, i.category")

	if role != "admin" { base = base.Where("fa.created_by = ?", userID) }
	if fieldID > 0     { base = base.Where("a.field_id = ?", fieldID) }
	if farmID > 0      { base = base.Where("fa.id = ?", farmID) }
	if seasonID > 0 {
		base = base.
			Joins("JOIN seasons se ON se.id = (SELECT season_id FROM plantings WHERE field_id = a.field_id AND status != 'lost' ORDER BY planting_date DESC LIMIT 1)").
			Where("a.application_date BETWEEN se.start_date AND se.end_date")
	}

	var rows []ports.CostRow
	err := base.Order("fi.name ASC, i.category ASC").Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) FindHarvestBagsPerField(
	userID uint, role string,
	fieldIDs []uint,
) (map[uint]float64, error) {

	type Row struct {
		FieldID   uint
		TotalBags float64
	}

	q := r.DB.Table("harvests h").
		Select("h.field_id, COALESCE(SUM(h.total_bags), 0) AS total_bags").
		Joins("JOIN fields fi ON fi.id = h.field_id").
		Joins("JOIN farms fa ON fa.id = fi.farm_id").
		Group("h.field_id")

	if role != "admin"      { q = q.Where("fa.created_by = ?", userID) }
	if len(fieldIDs) > 0   { q = q.Where("h.field_id IN ?", fieldIDs) }

	var rows []Row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[uint]float64, len(rows))
	for _, r := range rows {
		result[r.FieldID] = r.TotalBags
	}
	return result, nil
}

func (r *ReportRepository) FindOverview(userID uint, role string) (*ports.OverviewRow, error) {
	now := time.Now()
	thisYear  := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	in30Days  := now.AddDate(0, 0, 30)

	userFilter := ""
	userArgs := []interface{}{}
	if role != "admin" {
		userFilter = " AND fa.created_by = ?"
		userArgs = append(userArgs, userID)
	}

	row := &ports.OverviewRow{}

	farmQ := r.DB.Table("farms")
	if role != "admin" { farmQ = farmQ.Where("created_by = ?", userID) }
	farmQ.Count(&row.TotalFarms)

	type AreaRow struct { TotalFields int64; TotalArea float64 }
	var areaRow AreaRow
	fieldQ := r.DB.Table("fields fi").
		Select("COUNT(fi.id) AS total_fields, COALESCE(SUM(fi.area), 0) AS total_area").
		Joins("JOIN farms fa ON fa.id = fi.farm_id")
	if role != "admin" { fieldQ = fieldQ.Where("fa.created_by = ?", userID) }
	fieldQ.Scan(&areaRow)
	row.TotalFields = areaRow.TotalFields
	row.TotalAreaHa = areaRow.TotalArea

	type PlantedAreaRow struct { PlantedArea float64; ActiveCount int64 }
	var plantedRow PlantedAreaRow
	r.DB.Raw(`
		SELECT COALESCE(SUM(fi.area), 0) AS planted_area, COUNT(p.id) AS active_count
		FROM plantings p
		JOIN fields fi ON fi.id = p.field_id
		JOIN farms fa ON fa.id = fi.farm_id
		WHERE p.status = 'active'`+userFilter, userArgs...).Scan(&plantedRow)
	row.PlantedAreaHa  = plantedRow.PlantedArea
	row.ActivePlantings = plantedRow.ActiveCount

	r.DB.Table("seasons").Count(&row.TotalSeasons)

	r.DB.Raw(`
		SELECT COUNT(h.id) AS harvested_this_year,
		       COALESCE(AVG(h.productivity_bag_ha), 0) AS avg_productivity_bag_ha,
		       COALESCE(SUM(h.total_bags), 0) AS total_bags_this_year
		FROM harvests h
		JOIN fields fi ON fi.id = h.field_id
		JOIN farms fa ON fa.id = fi.farm_id
		WHERE h.harvest_date >= ?`+userFilter,
		append([]interface{}{thisYear}, userArgs...)...).Scan(row)

	inputQ := r.DB.Table("inputs")
	if role != "admin" { inputQ = inputQ.Where("created_by = ?", userID) }
	inputQ.Where("active = true").Count(&row.TotalInputTypes)
	inputQ.Where("active = true AND stock_qty <= min_stock_qty").Count(&row.LowStockInputs)
	inputQ.Where("active = true AND expiration_date IS NOT NULL AND expiration_date <= ?", in30Days).Count(&row.ExpiringInputs)

	r.DB.Raw(`
		SELECT COUNT(a.id) AS total_applications,
		       COALESCE(SUM(CASE WHEN a.application_date >= ? THEN 1 ELSE 0 END), 0) AS applications_this_month,
		       COALESCE(SUM(a.total_used * i.cost_per_unit), 0) AS estimated_total_cost
		FROM applications a
		JOIN fields fi ON fi.id = a.field_id
		JOIN farms fa ON fa.id = fi.farm_id
		JOIN inputs i ON i.id = a.input_id
		WHERE 1=1`+userFilter,
		append([]interface{}{thisMonth}, userArgs...)...).Scan(row)

	alertQ := r.DB.Table("alerts")
	if role != "admin" { alertQ = alertQ.Where("created_by = ?", userID) }
	alertQ.Where("status = 'open'").Count(&row.OpenAlerts)
	alertQ.Where("status = 'open' AND priority = 'high'").Count(&row.CriticalAlerts)

	return row, nil
}

func (r *ReportRepository) FindTopFields(userID uint, role string) ([]ports.TopFieldRow, error) {
	q := r.DB.Table("harvests h").
		Select("fi.name AS field_name, fa.name AS farm_name, h.productivity_bag_ha").
		Joins("JOIN fields fi ON fi.id = h.field_id").
		Joins("JOIN farms fa ON fa.id = fi.farm_id").
		Where("h.productivity_bag_ha > 0")
	if role != "admin" { q = q.Where("fa.created_by = ?", userID) }

	var items []ports.TopFieldRow
	err := q.Order("h.productivity_bag_ha DESC").Limit(5).Scan(&items).Error
	return items, err
}
