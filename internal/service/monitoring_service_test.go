package service

import (
	"errors"
	"testing"
	"time"

	"agrocontrol-api/internal/domain/entities"
	"agrocontrol-api/internal/dto"
)

// ── Mock MonitoringRepository ──────────────────────────────────────────────

type mockMonitoringRepo struct {
	monitoring *entities.Monitoring
	err        error
}

func (m *mockMonitoringRepo) Create(mon *entities.Monitoring) error {
	if m.err != nil {
		return m.err
	}
	mon.ID = 1
	return nil
}
func (m *mockMonitoringRepo) FindByID(id uint) (*entities.Monitoring, error) {
	return m.monitoring, m.err
}
func (m *mockMonitoringRepo) FindAll() ([]entities.Monitoring, error) { return nil, m.err }
func (m *mockMonitoringRepo) FindByUser(userID uint) ([]entities.Monitoring, error) {
	return nil, m.err
}
func (m *mockMonitoringRepo) FindByFieldID(fieldID uint) ([]entities.Monitoring, error) {
	if m.monitoring != nil {
		return []entities.Monitoring{*m.monitoring}, nil
	}
	return nil, m.err
}
func (m *mockMonitoringRepo) Update(mon *entities.Monitoring) error { return m.err }
func (m *mockMonitoringRepo) Delete(id uint) error                  { return m.err }

// ── Monitoring tests ───────────────────────────────────────────────────────

func TestCreateMonitoring_Success(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1, Status: entities.FieldStatusActive},
		belongs: true,
	}
	svc := NewMonitoringService(&mockMonitoringRepo{monitoring: &entities.Monitoring{ID: 1}}, fieldRepo, &mockAlertRepo{})

	resp, err := svc.CreateMonitoring(dto.CreateMonitoringRequest{
		FieldID:        1,
		InspectionDate: time.Now(),
		Type:           "Praga",
		ProblemName:    "Lagarta do cartucho",
		Severity:       entities.SeverityLow,
		Inspector:      "Willian",
	}, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if resp == nil {
		t.Fatal("esperava resposta válida")
	}
}

func TestCreateMonitoring_FieldNotFound(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{err: errors.New("not found")}
	svc := NewMonitoringService(&mockMonitoringRepo{}, fieldRepo, &mockAlertRepo{})

	_, err := svc.CreateMonitoring(dto.CreateMonitoringRequest{
		FieldID:     99,
		ProblemName: "Teste",
	}, 1, "manager")

	if err == nil {
		t.Fatal("esperava erro de talhão não encontrado")
	}
}

func TestCreateMonitoring_Forbidden(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1},
		belongs: false,
	}
	svc := NewMonitoringService(&mockMonitoringRepo{}, fieldRepo, &mockAlertRepo{})

	_, err := svc.CreateMonitoring(dto.CreateMonitoringRequest{
		FieldID:     1,
		ProblemName: "Teste",
	}, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestCreateMonitoring_UrgentGeraAlerta(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1},
		belongs: true,
	}
	alertRepo := &mockAlertRepo{}
	svc := NewMonitoringService(&mockMonitoringRepo{monitoring: &entities.Monitoring{ID: 1}}, fieldRepo, alertRepo)

	_, err := svc.CreateMonitoring(dto.CreateMonitoringRequest{
		FieldID:        1,
		InspectionDate: time.Now(),
		Type:           "Praga",
		ProblemName:    "Spodoptera",
		Severity:       entities.SeverityLow,
		Urgent:         true,
		Inspector:      "Willian",
	}, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
}

func TestCreateMonitoring_CriticalGeraAlerta(t *testing.T) {
	fieldRepo := &mockFieldRepoForPlanting{
		field:   &entities.Field{ID: 1},
		belongs: true,
	}
	alertRepo := &mockAlertRepo{}
	svc := NewMonitoringService(&mockMonitoringRepo{monitoring: &entities.Monitoring{ID: 1}}, fieldRepo, alertRepo)

	_, err := svc.CreateMonitoring(dto.CreateMonitoringRequest{
		FieldID:        1,
		InspectionDate: time.Now(),
		Type:           "Doença",
		ProblemName:    "Ferrugem asiática",
		Severity:       entities.SeverityCritical,
		Inspector:      "Willian",
	}, 1, "manager")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
}

func TestGetMonitoringByID_NotFound(t *testing.T) {
	svc := NewMonitoringService(
		&mockMonitoringRepo{err: errors.New("not found")},
		&mockFieldRepoForPlanting{},
		&mockAlertRepo{},
	)
	_, err := svc.GetMonitoringByID(99, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de not found")
	}
}

func TestGetMonitoringByID_AdminCanAccessAny(t *testing.T) {
	monitoring := &entities.Monitoring{ID: 1, FieldID: 1}
	fieldRepo := &mockFieldRepoForPlanting{field: &entities.Field{ID: 1}, belongs: false}
	svc := NewMonitoringService(&mockMonitoringRepo{monitoring: monitoring}, fieldRepo, &mockAlertRepo{})

	resp, err := svc.GetMonitoringByID(1, 1, "admin")

	if err != nil {
		t.Fatalf("admin deveria acessar qualquer monitoramento, got %v", err)
	}
	if resp == nil {
		t.Fatal("esperava resposta válida")
	}
}

func TestGetMonitoringByID_Forbidden(t *testing.T) {
	monitoring := &entities.Monitoring{ID: 1, FieldID: 1}
	fieldRepo := &mockFieldRepoForPlanting{field: &entities.Field{ID: 1}, belongs: false}
	svc := NewMonitoringService(&mockMonitoringRepo{monitoring: monitoring}, fieldRepo, &mockAlertRepo{})

	_, err := svc.GetMonitoringByID(1, 1, "operator")

	if err == nil {
		t.Fatal("esperava erro de forbidden")
	}
}

func TestGetMonitorings_Admin(t *testing.T) {
	svc := NewMonitoringService(&mockMonitoringRepo{}, &mockFieldRepoForPlanting{}, &mockAlertRepo{})
	result, err := svc.GetMonitorings(1, "admin")

	if err != nil {
		t.Fatalf("esperava nil, got %v", err)
	}
	if result == nil {
		t.Fatal("esperava slice vazio, não nil")
	}
}
