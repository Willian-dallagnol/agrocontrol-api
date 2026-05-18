package service

import (
	"agrocontrol-api/internal/cache"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/dto"
	"context"
	"fmt"
	"math"
	"time"
)

type ReportService struct {
	Repo  ports.ReportRepository
	Cache *cache.Client
}

func NewReportService(repo ports.ReportRepository, cache *cache.Client) *ReportService {
	return &ReportService{Repo: repo, Cache: cache}
}

func (s *ReportService) GetProductivityReport(q dto.ProductivityReportQuery, userID uint, role string) (*dto.ProductivityReportResponse, error) {
	rows, total, err := s.Repo.FindProductivity(userID, role, q.SeasonID, q.FarmID, q.CropID, q.Offset(), q.Limit)
	if err != nil {
		return nil, err
	}
	summaryRow, err := s.Repo.FindProductivitySummary(userID, role, q.SeasonID, q.FarmID, q.CropID)
	if err != nil {
		return nil, err
	}
	avgBagHa := summaryRow.AvgProductivityBagHa
	items := make([]dto.ProductivityReportItem, 0, len(rows))
	for _, r := range rows {
		harvestDate := ""
		if !r.HarvestDate.IsZero() {
			harvestDate = r.HarvestDate.Format("2006-01-02")
		}
		items = append(items, dto.ProductivityReportItem{
			PlantingID:        r.PlantingID,
			FieldID:           r.FieldID,
			FieldName:         r.FieldName,
			FarmID:            r.FarmID,
			FarmName:          r.FarmName,
			SeasonID:          r.SeasonID,
			SeasonName:        r.SeasonName,
			CropID:            r.CropID,
			CropName:          r.CropName,
			Variety:           r.Variety,
			PlantedAreaHa:     r.AreaHa,
			PlantingDate:      r.PlantingDate.Format("2006-01-02"),
			HarvestDate:       harvestDate,
			TotalBags:         r.TotalBags,
			ProductivityBagHa: r.ProductivityBagHa,
			ProductivityKgHa:  r.ProductivityKgHa,
			GrainMoisture:     r.GrainMoisture,
			FieldLossPct:      r.FieldLoss,
			EstimatedTotalKg:  round2(r.ProductivityKgHa * r.AreaHa),
			AboveAverageArea:  r.ProductivityBagHa > avgBagHa,
		})
	}
	totalPages := int(math.Ceil(float64(total) / float64(q.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return &dto.ProductivityReportResponse{
		Items:      items,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
		Summary: dto.ProductivitySummary{
			SeasonName:           summaryRow.SeasonName,
			TotalFields:          summaryRow.TotalFields,
			TotalPlantedAreaHa:   round2(summaryRow.TotalPlantedAreaHa),
			TotalBags:            round2(summaryRow.TotalBags),
			AvgProductivityBagHa: round2(summaryRow.AvgProductivityBagHa),
			AvgProductivityKgHa:  round2(summaryRow.AvgProductivityKgHa),
			BestField:            summaryRow.BestField,
			BestFieldBagHa:       round2(summaryRow.BestFieldBagHa),
		},
	}, nil
}

func (s *ReportService) GetCostPerFieldReport(q dto.CostPerFieldQuery, userID uint, role string) (*dto.CostPerFieldResponse, error) {
	rows, err := s.Repo.FindCostPerField(userID, role, q.FieldID, q.SeasonID, q.FarmID)
	if err != nil {
		return nil, err
	}
	type fieldAcc struct {
		item    dto.CostPerFieldItem
		hasData bool
	}
	fieldMap := make(map[uint]*fieldAcc)
	fieldOrder := []uint{}
	for _, r := range rows {
		acc, exists := fieldMap[r.FieldID]
		if !exists {
			acc = &fieldAcc{item: dto.CostPerFieldItem{FieldID: r.FieldID, FieldName: r.FieldName, FarmName: r.FarmName, AreaHa: r.AreaHa}, hasData: true}
			fieldMap[r.FieldID] = acc
			fieldOrder = append(fieldOrder, r.FieldID)
		}
		cost := r.TotalCost
		acc.item.TotalApplications += int(r.AppCount)
		switch r.Category {
		case "fertilizer":   acc.item.CostFertilizer += cost
		case "herbicide":    acc.item.CostHerbicide += cost
		case "fungicide":    acc.item.CostFungicide += cost
		case "insecticide":  acc.item.CostInsecticide += cost
		case "correctant":   acc.item.CostCorrectant += cost
		case "biological":   acc.item.CostBiological += cost
		default:             acc.item.CostOther += cost
		}
		acc.item.TotalCost += cost
	}
	fieldIDs := make([]uint, 0, len(fieldOrder))
	for _, id := range fieldOrder {
		fieldIDs = append(fieldIDs, id)
	}
	bagsMap, _ := s.Repo.FindHarvestBagsPerField(userID, role, fieldIDs)
	items := make([]dto.CostPerFieldItem, 0, len(fieldOrder))
	var totalCost, totalArea float64
	var mostExpensiveField string
	var mostExpensiveCostHa float64
	for _, id := range fieldOrder {
		acc := fieldMap[id]
		item := &acc.item
		item.TotalCost = round2(item.TotalCost)
		if item.AreaHa > 0 { item.CostPerHa = round2(item.TotalCost / item.AreaHa) }
		if bags, ok := bagsMap[id]; ok && bags > 0 { item.TotalBags = bags; item.CostPerBag = round2(item.TotalCost / bags) }
		item.CostFertilizer = round2(item.CostFertilizer)
		item.CostHerbicide  = round2(item.CostHerbicide)
		item.CostFungicide  = round2(item.CostFungicide)
		item.CostInsecticide = round2(item.CostInsecticide)
		item.CostCorrectant = round2(item.CostCorrectant)
		item.CostBiological = round2(item.CostBiological)
		item.CostOther      = round2(item.CostOther)
		totalCost += item.TotalCost
		totalArea += item.AreaHa
		if item.CostPerHa > mostExpensiveCostHa { mostExpensiveCostHa = item.CostPerHa; mostExpensiveField = item.FieldName }
		items = append(items, *item)
	}
	avgCostHa := 0.0
	if totalArea > 0 { avgCostHa = round2(totalCost / totalArea) }
	return &dto.CostPerFieldResponse{
		Items: items,
		Summary: dto.CostSummary{
			TotalFields: len(items), TotalAreaHa: round2(totalArea),
			TotalCost: round2(totalCost), AvgCostPerHa: avgCostHa,
			MostExpensiveField: mostExpensiveField, MostExpensiveCost: round2(mostExpensiveCostHa),
		},
	}, nil
}

func (s *ReportService) GetDashboardOverview(userID uint, role string) (*dto.DashboardOverviewResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("dashboard:overview:%d:%s", userID, role)
	var cached dto.DashboardOverviewResponse
	if found, _ := s.Cache.Get(ctx, cacheKey, &cached); found {
		return &cached, nil
	}
	row, err := s.Repo.FindOverview(userID, role)
	if err != nil {
		return nil, err
	}
	topFieldsRaw, err := s.Repo.FindTopFields(userID, role)
	if err != nil {
		topFieldsRaw = []ports.TopFieldRow{}
	}
	topFields := make([]dto.TopFieldItem, 0, len(topFieldsRaw))
	for _, f := range topFieldsRaw {
		topFields = append(topFields, dto.TopFieldItem{FieldName: f.FieldName, FarmName: f.FarmName, ProductivityBagHa: f.ProductivityBagHa})
	}
	result := &dto.DashboardOverviewResponse{
		TotalFarms:            row.TotalFarms,
		TotalFields:           row.TotalFields,
		TotalAreaHa:           round2(row.TotalAreaHa),
		PlantedAreaHa:         round2(row.PlantedAreaHa),
		TotalSeasons:          row.TotalSeasons,
		ActivePlantings:       row.ActivePlantings,
		HarvestedThisYear:     row.HarvestedThisYear,
		AvgProductivityBagHa:  round2(row.AvgProductivityBagHa),
		TotalBagsThisYear:     round2(row.TotalBagsThisYear),
		TotalInputTypes:       row.TotalInputTypes,
		LowStockInputs:        row.LowStockInputs,
		ExpiringInputs:        row.ExpiringInputs,
		TotalApplications:     row.TotalApplications,
		ApplicationsThisMonth: row.ApplicationsThisMonth,
		EstimatedTotalCost:    round2(row.EstimatedTotalCost),
		OpenAlerts:            row.OpenAlerts,
		CriticalAlerts:        row.CriticalAlerts,
		TopFields:             topFields,
	}
	_ = s.Cache.Set(ctx, cacheKey, result, 5*time.Minute)
	return result, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
