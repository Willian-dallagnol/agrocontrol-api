package service

import (
	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/dto"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// ── InputService ──────────────────────────────────────────────────────────────

type InputService struct {
	Repo      ports.InputRepository
	AlertRepo ports.AlertRepository
	TxRunner  ports.TxRunner
}

func NewInputService(repo ports.InputRepository, alertRepo ports.AlertRepository, tx ports.TxRunner) *InputService {
	return &InputService{Repo: repo, AlertRepo: alertRepo, TxRunner: tx}
}

func (s *InputService) CreateInput(req dto.CreateInputRequest, userID uint) (*dto.InputResponse, error) {
	if req.StockQty < 0 {
		return nil, fmt.Errorf("estoque inicial não pode ser negativo: %w", apperrors.ErrInvalidInput)
	}
	if req.CostPerUnit < 0 {
		return nil, fmt.Errorf("custo por unidade não pode ser negativo: %w", apperrors.ErrInvalidInput)
	}
	input := &entities.Input{
		Name:           strings.TrimSpace(req.Name),
		Category:       req.Category,
		Manufacturer:   strings.TrimSpace(req.Manufacturer),
		BatchNumber:    req.BatchNumber,
		ExpirationDate: req.ExpirationDate,
		Unit:           strings.TrimSpace(req.Unit),
		StockQty:       req.StockQty,
		MinStockQty:    req.MinStockQty,
		CostPerUnit:    req.CostPerUnit,
		Active:         true,
		CreatedBy:      userID,
	}
	if err := s.Repo.Create(input); err != nil {
		return nil, fmt.Errorf("erro ao criar insumo: %w", err)
	}
	slog.Info("input: criado", "input_id", input.ID, "category", input.Category)
	return toInputResponse(input), nil
}

func (s *InputService) GetInputsPaginated(q dto.InputQuery, userID uint, role string) (dto.PaginatedResponse, error) {
	var inputs []entities.Input
	var total int64
	var err error

	if role == "admin" {
		inputs, total, err = s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.Search, q.Category)
	} else {
		inputs, total, err = s.Repo.FindByUserPaginated(userID, q.Offset(), q.Limit, q.Search, q.Category)
	}
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar insumos: %w", err)
	}

	resp := make([]dto.InputResponse, 0, len(inputs))
	for i := range inputs {
		resp = append(resp, *toInputResponse(&inputs[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *InputService) GetInputByID(id uint, userID uint, role string) (*dto.InputResponse, error) {
	input, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("insumo", id)
	}
	if role != "admin" && input.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("acessar este insumo")
	}
	return toInputResponse(input), nil
}

func (s *InputService) UpdateInput(id uint, req dto.UpdateInputRequest, userID uint, role string) (*dto.InputResponse, error) {
	input, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("insumo", id)
	}
	if role != "admin" && input.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("editar este insumo")
	}
	input.Name = strings.TrimSpace(req.Name)
	input.Manufacturer = strings.TrimSpace(req.Manufacturer)
	input.BatchNumber = req.BatchNumber
	input.ExpirationDate = req.ExpirationDate
	input.Unit = strings.TrimSpace(req.Unit)
	input.MinStockQty = req.MinStockQty
	input.CostPerUnit = req.CostPerUnit
	input.Active = req.Active

	if err := s.Repo.Update(input); err != nil {
		return nil, fmt.Errorf("erro ao atualizar insumo: %w", err)
	}
	return toInputResponse(input), nil
}

func (s *InputService) DeleteInput(id uint, userID uint, role string) error {
	input, err := s.Repo.FindByID(id)
	if err != nil {
		return apperrors.NotFoundError("insumo", id)
	}
	if role != "admin" && input.CreatedBy != userID {
		return apperrors.ForbiddenError("excluir este insumo")
	}
	return s.Repo.Delete(id)
}

func (s *InputService) AdjustStock(id uint, req dto.AdjustStockRequest, userID uint, role string) (*dto.InputResponse, error) {
	input, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("insumo", id)
	}
	if role != "admin" && input.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("ajustar estoque deste insumo")
	}

	newQty := input.StockQty + req.Quantity
	if newQty < 0 {
		return nil, fmt.Errorf("ajuste resultaria em estoque negativo (atual: %.2f, ajuste: %.2f): %w",
			input.StockQty, req.Quantity, apperrors.ErrInsufficientStock)
	}

	input.StockQty = newQty
	if err := s.Repo.Update(input); err != nil {
		return nil, fmt.Errorf("erro ao ajustar estoque: %w", err)
	}
	slog.Info("input: estoque ajustado", "input_id", id, "delta", req.Quantity, "new_qty", newQty)
	return toInputResponse(input), nil
}

func toInputResponse(i *entities.Input) *dto.InputResponse {
	return &dto.InputResponse{
		ID: i.ID, Name: i.Name, Category: i.Category, Manufacturer: i.Manufacturer,
		BatchNumber: i.BatchNumber, ExpirationDate: i.ExpirationDate, Unit: i.Unit,
		StockQty: i.StockQty, MinStockQty: i.MinStockQty, CostPerUnit: i.CostPerUnit,
		Active: i.Active, LowStock: i.IsLowStock(),
		CreatedAt: i.CreatedAt, UpdatedAt: i.UpdatedAt,
	}
}

// ── ApplicationService ────────────────────────────────────────────────────────

type ApplicationService struct {
	Repo      ports.ApplicationRepository
	FieldRepo ports.FieldRepository
	InputRepo ports.InputRepository
	AlertRepo ports.AlertRepository
	TxRunner  ports.TxRunner
}

func NewApplicationService(
	repo ports.ApplicationRepository, fieldRepo ports.FieldRepository,
	inputRepo ports.InputRepository, alertRepo ports.AlertRepository, tx ports.TxRunner,
) *ApplicationService {
	return &ApplicationService{Repo: repo, FieldRepo: fieldRepo, InputRepo: inputRepo, AlertRepo: alertRepo, TxRunner: tx}
}

func (s *ApplicationService) CreateApplication(req dto.CreateApplicationRequest, userID uint, role string) (*dto.ApplicationResponse, error) {
	field, err := s.FieldRepo.FindByID(req.FieldID)
	if err != nil {
		return nil, apperrors.NotFoundError("talhão", req.FieldID)
	}
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(req.FieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("registrar aplicação neste talhão")
		}
	}
	if field.Status == entities.FieldStatusInactive {
		return nil, apperrors.ErrInactiveField
	}

	input, err := s.InputRepo.FindByID(req.InputID)
	if err != nil {
		return nil, apperrors.NotFoundError("insumo", req.InputID)
	}
	if !input.Active {
		return nil, fmt.Errorf("insumo inativo: %w", apperrors.ErrInvalidInput)
	}

	totalUsed := req.DosePerHa * field.Area
	if totalUsed <= 0 {
		return nil, fmt.Errorf("dose_por_ha e área do talhão resultam em uso zero: %w", apperrors.ErrInvalidInput)
	}
	if input.StockQty < totalUsed {
		return nil, fmt.Errorf("estoque insuficiente: disponível %.2f %s, necessário %.2f: %w",
			input.StockQty, input.Unit, totalUsed, apperrors.ErrInsufficientStock)
	}

	var app entities.Application

	// ── Transação via TxRunner — sem *gorm.DB no serviço ──────────────────
	err = s.TxRunner.RunInTx(func(tx ports.TxRunner) error {
		app = entities.Application{
			FieldID: req.FieldID, PlantingID: req.PlantingID, InputID: req.InputID,
			ApplicationType: req.ApplicationType, ApplicationDate: req.ApplicationDate,
			DosePerHa: req.DosePerHa, TotalUsed: totalUsed, SprayVolume: req.SprayVolume,
			Target: req.Target, Equipment: req.Equipment, Operator: req.Operator,
			WindSpeed: req.WindSpeed, Temperature: req.Temperature, Humidity: req.Humidity,
			Notes: req.Notes, CreatedBy: userID,
		}
		if err := s.Repo.CreateTx(tx, &app); err != nil {
			return err
		}
		if err := s.InputRepo.DeductStockTx(tx, req.InputID, totalUsed); err != nil {
			return err
		}
		updatedInput, _ := s.InputRepo.FindByIDTx(tx, req.InputID)
		if updatedInput != nil && updatedInput.IsLowStock() {
			alert := &entities.Alert{
				Title:       "Estoque baixo: " + updatedInput.Name,
				Type:        entities.AlertTypeLowStock,
				Description: fmt.Sprintf("Estoque atual (%.2f %s) abaixo do mínimo (%.2f)", updatedInput.StockQty, updatedInput.Unit, updatedInput.MinStockQty),
				Priority:    entities.AlertPriorityHigh,
				Status:      entities.AlertStatusOpen,
				RefID:       &updatedInput.ID,
				RefType:     "input",
				CreatedBy:   userID,
			}
			_ = s.AlertRepo.CreateTx(tx, alert)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao registrar aplicação: %w", err)
	}

	slog.Info("application: criada", "app_id", app.ID, "field_id", req.FieldID, "input_id", req.InputID, "total_used", totalUsed)

	full, _ := s.Repo.FindByID(app.ID)
	return toApplicationResponse(full), nil
}

func (s *ApplicationService) GetApplicationsPaginated(q dto.ApplicationQuery, userID uint, role string) (dto.PaginatedResponse, error) {
	var apps []entities.Application
	var total int64
	var err error

	if role == "admin" {
		apps, total, err = s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.FieldID, q.Type)
	} else {
		apps, total, err = s.Repo.FindByUserPaginated(userID, q.Offset(), q.Limit, q.FieldID, q.Type)
	}
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar aplicações: %w", err)
	}
	resp := make([]dto.ApplicationResponse, 0, len(apps))
	for i := range apps {
		resp = append(resp, *toApplicationResponse(&apps[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *ApplicationService) GetApplicationByID(id uint, userID uint, role string) (*dto.ApplicationResponse, error) {
	app, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("aplicação", id)
	}
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(app.FieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar esta aplicação")
		}
	}
	return toApplicationResponse(app), nil
}

func (s *ApplicationService) GetApplicationsByField(fieldID uint, userID uint, role string) ([]dto.ApplicationResponse, error) {
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(fieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar aplicações deste talhão")
		}
	}
	apps, err := s.Repo.FindByFieldID(fieldID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aplicações: %w", err)
	}
	resp := make([]dto.ApplicationResponse, 0, len(apps))
	for i := range apps {
		resp = append(resp, *toApplicationResponse(&apps[i]))
	}
	return resp, nil
}

func toApplicationResponse(a *entities.Application) *dto.ApplicationResponse {
	r := &dto.ApplicationResponse{
		ID: a.ID, FieldID: a.FieldID, PlantingID: a.PlantingID, InputID: a.InputID,
		ApplicationType: a.ApplicationType, ApplicationDate: a.ApplicationDate,
		DosePerHa: a.DosePerHa, TotalUsed: a.TotalUsed, SprayVolume: a.SprayVolume,
		Target: a.Target, Equipment: a.Equipment, Operator: a.Operator,
		WindSpeed: a.WindSpeed, Temperature: a.Temperature, Humidity: a.Humidity,
		Notes: a.Notes, CreatedBy: a.CreatedBy, CreatedAt: a.CreatedAt,
	}
	if a.Field.ID > 0 {
		r.FieldName = a.Field.Name
	}
	if a.Input.ID > 0 {
		r.InputName = a.Input.Name
	}
	return r
}

// ── MonitoringService ─────────────────────────────────────────────────────────

type MonitoringService struct {
	Repo      ports.MonitoringRepository
	FieldRepo ports.FieldRepository
	AlertRepo ports.AlertRepository
}

func NewMonitoringService(repo ports.MonitoringRepository, fieldRepo ports.FieldRepository, alertRepo ports.AlertRepository) *MonitoringService {
	return &MonitoringService{Repo: repo, FieldRepo: fieldRepo, AlertRepo: alertRepo}
}

func (s *MonitoringService) CreateMonitoring(req dto.CreateMonitoringRequest, userID uint, role string) (*dto.MonitoringResponse, error) {
	if _, err := s.FieldRepo.FindByID(req.FieldID); err != nil {
		return nil, apperrors.NotFoundError("talhão", req.FieldID)
	}
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(req.FieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("registrar monitoramento neste talhão")
		}
	}

	m := &entities.Monitoring{
		FieldID: req.FieldID, PlantingID: req.PlantingID, InspectionDate: req.InspectionDate,
		Type: req.Type, ProblemName: strings.TrimSpace(req.ProblemName),
		InfestationLevel: req.InfestationLevel, Severity: req.Severity,
		CropStage: req.CropStage, TechnicalRec: req.TechnicalRec,
		Urgent: req.Urgent, Inspector: strings.TrimSpace(req.Inspector),
		Notes: req.Notes, CreatedBy: userID,
	}
	if err := s.Repo.Create(m); err != nil {
		return nil, fmt.Errorf("erro ao criar monitoramento: %w", err)
	}

	if req.Urgent || req.Severity == entities.SeverityCritical {
		fieldID := req.FieldID
		alert := &entities.Alert{
			Title:       "Monitoramento urgente: " + req.ProblemName,
			Type:        entities.AlertTypePest,
			Description: fmt.Sprintf("Severidade %s detectada no talhão %d — ação imediata necessária", req.Severity, req.FieldID),
			Priority:    entities.AlertPriorityHigh,
			Status:      entities.AlertStatusOpen,
			RefID:       &fieldID,
			RefType:     "field",
			CreatedBy:   userID,
		}
		if err := s.AlertRepo.Create(alert); err != nil {
			slog.Warn("monitoring: falha ao criar alerta automático", "monitoring_id", m.ID, "error", err)
		}
	}

	full, _ := s.Repo.FindByID(m.ID)
	return toMonitoringResponse(full), nil
}

func (s *MonitoringService) GetMonitorings(userID uint, role string) ([]dto.MonitoringResponse, error) {
	var mons []entities.Monitoring
	var err error
	if role == "admin" {
		mons, err = s.Repo.FindAll()
	} else {
		mons, err = s.Repo.FindByUser(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar monitoramentos: %w", err)
	}
	resp := make([]dto.MonitoringResponse, 0, len(mons))
	for i := range mons {
		resp = append(resp, *toMonitoringResponse(&mons[i]))
	}
	return resp, nil
}

func (s *MonitoringService) GetMonitoringByID(id uint, userID uint, role string) (*dto.MonitoringResponse, error) {
	m, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("monitoramento", id)
	}
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(m.FieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar este monitoramento")
		}
	}
	return toMonitoringResponse(m), nil
}

func (s *MonitoringService) GetMonitoringsByField(fieldID uint, userID uint, role string) ([]dto.MonitoringResponse, error) {
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(fieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar monitoramentos deste talhão")
		}
	}
	mons, err := s.Repo.FindByFieldID(fieldID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar monitoramentos: %w", err)
	}
	resp := make([]dto.MonitoringResponse, 0, len(mons))
	for i := range mons {
		resp = append(resp, *toMonitoringResponse(&mons[i]))
	}
	return resp, nil
}

func toMonitoringResponse(m *entities.Monitoring) *dto.MonitoringResponse {
	r := &dto.MonitoringResponse{
		ID: m.ID, FieldID: m.FieldID, PlantingID: m.PlantingID, InspectionDate: m.InspectionDate,
		Type: m.Type, ProblemName: m.ProblemName, InfestationLevel: m.InfestationLevel,
		Severity: m.Severity, CropStage: m.CropStage, TechnicalRec: m.TechnicalRec,
		Urgent: m.Urgent, Inspector: m.Inspector, Notes: m.Notes,
		CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
	if m.Field.ID > 0 {
		r.FieldName = m.Field.Name
	}
	return r
}

// ── HarvestService ────────────────────────────────────────────────────────────

type HarvestService struct {
	Repo         ports.HarvestRepository
	PlantingRepo ports.PlantingRepository
	FieldRepo    ports.FieldRepository
	TxRunner     ports.TxRunner
}

func NewHarvestService(repo ports.HarvestRepository, plantingRepo ports.PlantingRepository, fieldRepo ports.FieldRepository, tx ports.TxRunner) *HarvestService {
	return &HarvestService{Repo: repo, PlantingRepo: plantingRepo, FieldRepo: fieldRepo, TxRunner: tx}
}

func (s *HarvestService) CreateHarvest(req dto.CreateHarvestRequest, userID uint, role string) (*dto.HarvestResponse, error) {
	planting, err := s.PlantingRepo.FindByID(req.PlantingID)
	if err != nil {
		return nil, apperrors.NotFoundError("plantio", req.PlantingID)
	}
	if role != "admin" {
		ok, _ := s.PlantingRepo.BelongsToUser(req.PlantingID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("registrar colheita deste plantio")
		}
	}
	if planting.Status != entities.PlantingStatusActive {
		return nil, apperrors.ErrNoActivePlanting
	}
	exists, err := s.Repo.ExistsByPlantingID(req.PlantingID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar colheita existente: %w", err)
	}
	if exists {
		return nil, apperrors.ErrAlreadyHarvested
	}

	var harvest entities.Harvest

	err = s.TxRunner.RunInTx(func(tx ports.TxRunner) error {
		harvest = entities.Harvest{
			PlantingID:    req.PlantingID,
			FieldID:       planting.FieldID,
			HarvestDate:   req.HarvestDate,
			TotalBags:     req.TotalBags,
			GrainMoisture: req.GrainMoisture,
			Impurity:      req.Impurity,
			FieldLoss:     req.FieldLoss,
			Notes:         req.Notes,
			CreatedBy:     userID,
		}
		// Calcula produtividade usando regra da entidade + área real do talhão
		field, err := s.FieldRepo.FindByID(planting.FieldID)
		if err == nil && field.Area > 0 {
			harvest.CalculateProductivity(field.Area)
		}
		if err := s.Repo.CreateTx(tx, &harvest); err != nil {
			return err
		}
		planting.Status = entities.PlantingStatusHarvested
		return s.PlantingRepo.Update(planting)
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao registrar colheita: %w", err)
	}

	slog.Info("harvest: registrada", "harvest_id", harvest.ID, "planting_id", req.PlantingID,
		"productivity_bag_ha", harvest.ProductivityBagHa)

	full, _ := s.Repo.FindByID(harvest.ID)
	return toHarvestResponse(full), nil
}

func (s *HarvestService) GetHarvestsPaginated(q dto.HarvestQuery, userID uint, role string) (dto.PaginatedResponse, error) {
	var harvests []entities.Harvest
	var total int64
	var err error

	if role == "admin" {
		harvests, total, err = s.Repo.FindAllPaginated(q.Offset(), q.Limit, q.FieldID)
	} else {
		harvests, total, err = s.Repo.FindByUserPaginated(userID, q.Offset(), q.Limit, q.FieldID)
	}
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("erro ao buscar colheitas: %w", err)
	}
	resp := make([]dto.HarvestResponse, 0, len(harvests))
	for i := range harvests {
		resp = append(resp, *toHarvestResponse(&harvests[i]))
	}
	return dto.NewPaginatedResponse(resp, total, q.Page, q.Limit), nil
}

func (s *HarvestService) GetHarvestByID(id uint, userID uint, role string) (*dto.HarvestResponse, error) {
	h, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("colheita", id)
	}
	if role != "admin" {
		ok, _ := s.FieldRepo.BelongsToUser(h.FieldID, userID)
		if !ok {
			return nil, apperrors.ForbiddenError("acessar esta colheita")
		}
	}
	return toHarvestResponse(h), nil
}

func toHarvestResponse(h *entities.Harvest) *dto.HarvestResponse {
	r := &dto.HarvestResponse{
		ID: h.ID, PlantingID: h.PlantingID, FieldID: h.FieldID, HarvestDate: h.HarvestDate,
		ProductivityBagHa: h.ProductivityBagHa, ProductivityKgHa: h.ProductivityKgHa,
		TotalBags: h.TotalBags, GrainMoisture: h.GrainMoisture, Impurity: h.Impurity,
		FieldLoss: h.FieldLoss, Notes: h.Notes, CreatedBy: h.CreatedBy,
		CreatedAt: h.CreatedAt, UpdatedAt: h.UpdatedAt,
	}
	if h.Field.ID > 0 {
		r.FieldName = h.Field.Name
	}
	return r
}

// ── AlertService ──────────────────────────────────────────────────────────────

type AlertService struct {
	Repo ports.AlertRepository
}

func NewAlertService(repo ports.AlertRepository) *AlertService {
	return &AlertService{Repo: repo}
}

func (s *AlertService) CreateAlert(req dto.CreateAlertRequest, userID uint) (*dto.AlertResponse, error) {
	priority := req.Priority
	if priority == "" {
		priority = entities.AlertPriorityMedium
	}
	alert := &entities.Alert{
		Title:       strings.TrimSpace(req.Title),
		Type:        req.Type,
		Description: req.Description,
		Priority:    priority,
		Status:      entities.AlertStatusOpen,
		RefID:       req.RefID,
		RefType:     req.RefType,
		CreatedBy:   userID,
	}
	if err := s.Repo.Create(alert); err != nil {
		return nil, fmt.Errorf("erro ao criar alerta: %w", err)
	}
	return toAlertResponse(alert), nil
}

func (s *AlertService) GetAlerts(userID uint, role string) ([]dto.AlertResponse, error) {
	var alerts []entities.Alert
	var err error
	if role == "admin" {
		alerts, err = s.Repo.FindAll()
	} else {
		alerts, err = s.Repo.FindByUser(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar alertas: %w", err)
	}
	resp := make([]dto.AlertResponse, 0, len(alerts))
	for i := range alerts {
		resp = append(resp, *toAlertResponse(&alerts[i]))
	}
	return resp, nil
}

func (s *AlertService) GetOpenAlerts(userID uint, role string) ([]dto.AlertResponse, error) {
	var alerts []entities.Alert
	var err error
	if role == "admin" {
		alerts, err = s.Repo.FindOpen()
	} else {
		alerts, err = s.Repo.FindOpenByUser(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar alertas abertos: %w", err)
	}
	resp := make([]dto.AlertResponse, 0, len(alerts))
	for i := range alerts {
		resp = append(resp, *toAlertResponse(&alerts[i]))
	}
	return resp, nil
}

func (s *AlertService) GetAlertByID(id uint, userID uint, role string) (*dto.AlertResponse, error) {
	alert, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("alerta", id)
	}
	if role != "admin" && alert.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("acessar este alerta")
	}
	return toAlertResponse(alert), nil
}

func (s *AlertService) UpdateStatus(id uint, req dto.UpdateAlertStatusRequest, userID uint, role string) (*dto.AlertResponse, error) {
	alert, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFoundError("alerta", id)
	}
	if role != "admin" && alert.CreatedBy != userID {
		return nil, apperrors.ForbiddenError("atualizar este alerta")
	}
	alert.Status = req.Status
	if req.Status == entities.AlertStatusResolved {
		now := time.Now()
		alert.ResolvedAt = &now
	}
	if err := s.Repo.Update(alert); err != nil {
		return nil, fmt.Errorf("erro ao atualizar alerta: %w", err)
	}
	return toAlertResponse(alert), nil
}

func toAlertResponse(a *entities.Alert) *dto.AlertResponse {
	return &dto.AlertResponse{
		ID: a.ID, Title: a.Title, Type: a.Type, Description: a.Description,
		Priority: a.Priority, Status: a.Status, RefID: a.RefID, RefType: a.RefType,
		CreatedBy: a.CreatedBy, ResolvedAt: a.ResolvedAt, CreatedAt: a.CreatedAt,
	}
}
