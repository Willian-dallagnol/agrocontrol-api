package service

import (
	"errors"
	"testing"
	"time"

	"agrocontrol-api/internal/apperrors"
	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/domain/ports"
	"agrocontrol-api/internal/dto"
)

// ── Mock InputRepository ───────────────────────────────────────────────────

type mockInputRepo struct {
	input *entities.Input
	err   error
}

func (m *mockInputRepo) Create(i *entities.Input) error {
	if m.err != nil {
		return m.err
	}
	i.ID = 1
	return nil
}
func (m *mockInputRepo) FindByID(id uint) (*entities.Input, error)           { return m.input, m.err }
func (m *mockInputRepo) FindAll() ([]entities.Input, error)                  { return nil, m.err }
func (m *mockInputRepo) FindByUser(userID uint) ([]entities.Input, error)    { return nil, m.err }
func (m *mockInputRepo) FindLowStock() ([]entities.Input, error)             { return nil, m.err }
func (m *mockInputRepo) FindExpiringSoon(days int) ([]entities.Input, error) { return nil, m.err }
func (m *mockInputRepo) FindAllPaginated(offset, limit int, search, category string) ([]entities.Input, int64, error) {
	return nil, 0, m.err
}
func (m *mockInputRepo) FindByUserPaginated(userID uint, offset, limit int, search, category string) ([]entities.Input, int64, error) {
	return nil, 0, m.err
}
func (m *mockInputRepo) Update(i *entities.Input) error                              { return m.err }
func (m *mockInputRepo) Delete(id uint) error                                        { return m.err }
func (m *mockInputRepo) CountLowStock() (int64, error)                               { return 0, m.err }
func (m *mockInputRepo) CreateTx(tx ports.TxRunner, i *entities.Input) error         { return m.err }
func (m *mockInputRepo) DeductStockTx(tx ports.TxRunner, id uint, qty float64) error { return m.err }
func (m *mockInputRepo) FindByIDTx(tx ports.TxRunner, id uint) (*entities.Input, error) {
	return m.input, m.err
}

// ── Mock AlertRepository ───────────────────────────────────────────────────

type mockAlertRepo struct {
	alert *entities.Alert
	err   error
}

func (m *mockAlertRepo) Create(a *entities.Alert) error {
	if m.err != nil {
		return m.err
	}
	a.ID = 1
	return nil
}
func (m *mockAlertRepo) FindByID(id uint) (*entities.Alert, error)            { return m.alert, m.err }
func (m *mockAlertRepo) FindAll() ([]entities.Alert, error)                   { return nil, m.err }
func (m *mockAlertRepo) FindByUser(userID uint) ([]entities.Alert, error)     { return nil, m.err }
func (m *mockAlertRepo) FindOpen() ([]entities.Alert, error)                  { return nil, m.err }
func (m *mockAlertRepo) FindOpenByUser(userID uint) ([]entities.Alert, error) { return nil, m.err }
func (m *mockAlertRepo) Update(a *entities.Alert) error                       { return m.err }
func (m *mockAlertRepo) Delete(id uint) error                                 { return m.err }
func (m *mockAlertRepo) CountOpen() (int64, error)                            { return 0, m.err }
func (m *mockAlertRepo) CountOpenByUser(userID uint) (int64, error)           { return 0, m.err }
func (m *mockAlertRepo) CreateTx(tx ports.TxRunner, a *entities.Alert) error  { return m.err }

// ── Mock TxRunner ──────────────────────────────────────────────────────────

type mockTxRunner struct{}

func (m *mockTxRunner) RunInTx(fn func(tx ports.TxRunner) error) error {
	return fn(m)
}

// ── Input tests ────────────────────────────────────────────────────────────

func TestCreateInput_Success(t *testing.T) {
	svc := NewInputService(&mockInputRepo{}, &mockAlertRepo{}, &mockTxRunner{})
	resp, err := svc.CreateInput(dto.CreateInputRequest{
		Name:        "Glifosato",
		Category:    "Herbicida",
		Unit:        "L",
		StockQty:    100,
		CostPerUnit: 25.50,
	}, 1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Name != "Glifosato" {
		t.Errorf("esperava 'Glifosato', got '%s'", resp.Name)
	}
}

func TestCreateInput_NegativeStock(t *testing.T) {
	svc := NewInputService(&mockInputRepo{}, &mockAlertRepo{}, &mockTxRunner{})
	_, err := svc.CreateInput(dto.CreateInputRequest{
		Name:     "Insumo Inválido",
		StockQty: -10,
	}, 1)

	if err == nil {
		t.Fatal("esperava erro para estoque negativo")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestCreateInput_NegativeCost(t *testing.T) {
	svc := NewInputService(&mockInputRepo{}, &mockAlertRepo{}, &mockTxRunner{})
	_, err := svc.CreateInput(dto.CreateInputRequest{
		Name:        "Insumo Inválido",
		StockQty:    10,
		CostPerUnit: -5,
	}, 1)

	if err == nil {
		t.Fatal("esperava erro para custo negativo")
	}
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("esperava ErrInvalidInput, got %v", err)
	}
}

func TestGetInputByID_NotFound(t *testing.T) {
	svc := NewInputService(&mockInputRepo{err: errors.New("not found")}, &mockAlertRepo{}, &mockTxRunner{})
	_, err := svc.GetInputByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetInputByID_Forbidden(t *testing.T) {
	input := &entities.Input{ID: 1, CreatedBy: 2}
	svc := NewInputService(&mockInputRepo{input: input}, &mockAlertRepo{}, &mockTxRunner{})
	_, err := svc.GetInputByID(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestAdjustStock_Success(t *testing.T) {
	input := &entities.Input{ID: 1, CreatedBy: 1, StockQty: 100}
	svc := NewInputService(&mockInputRepo{input: input}, &mockAlertRepo{}, &mockTxRunner{})
	resp, err := svc.AdjustStock(1, dto.AdjustStockRequest{Quantity: 50}, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.StockQty != 150 {
		t.Errorf("esperava estoque 150, got %.2f", resp.StockQty)
	}
}

func TestAdjustStock_NegativeResult(t *testing.T) {
	input := &entities.Input{ID: 1, CreatedBy: 1, StockQty: 10}
	svc := NewInputService(&mockInputRepo{input: input}, &mockAlertRepo{}, &mockTxRunner{})
	_, err := svc.AdjustStock(1, dto.AdjustStockRequest{Quantity: -50}, 1, "manager")

	if err == nil {
		t.Fatal("esperava erro para estoque negativo")
	}
	if !errors.Is(err, apperrors.ErrInsufficientStock) {
		t.Errorf("esperava ErrInsufficientStock, got %v", err)
	}
}

// ── Alert tests ────────────────────────────────────────────────────────────

func TestCreateAlert_Success(t *testing.T) {
	svc := NewAlertService(&mockAlertRepo{})
	resp, err := svc.CreateAlert(dto.CreateAlertRequest{
		Title: "Estoque baixo",
		Type:  entities.AlertTypeLowStock,
	}, 1)

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Title != "Estoque baixo" {
		t.Errorf("esperava 'Estoque baixo', got '%s'", resp.Title)
	}
	if resp.Priority != entities.AlertPriorityMedium {
		t.Errorf("esperava priority medium por padrão, got '%s'", resp.Priority)
	}
}

func TestGetAlertByID_NotFound(t *testing.T) {
	svc := NewAlertService(&mockAlertRepo{err: errors.New("not found")})
	_, err := svc.GetAlertByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetAlertByID_Forbidden(t *testing.T) {
	alert := &entities.Alert{ID: 1, CreatedBy: 2}
	svc := NewAlertService(&mockAlertRepo{alert: alert})
	_, err := svc.GetAlertByID(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestUpdateAlertStatus_Resolved(t *testing.T) {
	alert := &entities.Alert{ID: 1, CreatedBy: 1, Status: entities.AlertStatusOpen}
	svc := NewAlertService(&mockAlertRepo{alert: alert})
	resp, err := svc.UpdateStatus(1, dto.UpdateAlertStatusRequest{
		Status: entities.AlertStatusResolved,
	}, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp.Status != entities.AlertStatusResolved {
		t.Errorf("esperava status resolved, got '%s'", resp.Status)
	}
	if resp.ResolvedAt == nil {
		t.Error("esperava ResolvedAt preenchido")
	}
}

func TestUpdateAlertStatus_AdminCanUpdateAny(t *testing.T) {
	alert := &entities.Alert{ID: 1, CreatedBy: 99, Status: entities.AlertStatusOpen}
	svc := NewAlertService(&mockAlertRepo{alert: alert})
	_, err := svc.UpdateStatus(1, dto.UpdateAlertStatusRequest{
		Status: entities.AlertStatusResolved,
	}, 1, "admin")

	if err != nil {
		t.Fatalf("admin deveria atualizar qualquer alerta, got %v", err)
	}
}

func init() {
	_ = time.Now()
}
