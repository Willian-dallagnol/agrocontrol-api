package entities

import (
	"testing"
	"time"
)

// ── Harvest tests ──────────────────────────────────────────────────────────

func TestHarvest_CalculateProductivity(t *testing.T) {
	h := &Harvest{TotalBags: 300}
	h.CalculateProductivity(5.0)

	if h.ProductivityBagHa != 60 {
		t.Errorf("esperava 60 sc/ha, got %.2f", h.ProductivityBagHa)
	}
	if h.ProductivityKgHa != 3600 {
		t.Errorf("esperava 3600 kg/ha, got %.2f", h.ProductivityKgHa)
	}
}

func TestHarvest_CalculateProductivity_InvalidArea(t *testing.T) {
	h := &Harvest{TotalBags: 300}
	h.CalculateProductivity(0)

	if h.ProductivityBagHa != 0 {
		t.Error("área zero não deve calcular produtividade")
	}
}

func TestHarvest_IsHighYield(t *testing.T) {
	h := &Harvest{ProductivityBagHa: 65}
	if !h.IsHighYield(60) {
		t.Error("65 sc/ha deveria ser alta produtividade acima de 60")
	}
	if h.IsHighYield(70) {
		t.Error("65 sc/ha não deveria ser alta produtividade acima de 70")
	}
}

func TestHarvest_AdjustedTotalBags(t *testing.T) {
	h := &Harvest{TotalBags: 100, FieldLoss: 10}
	adjusted := h.AdjustedTotalBags()
	if adjusted != 90 {
		t.Errorf("esperava 90 sacas ajustadas, got %.2f", adjusted)
	}
}

func TestHarvest_AdjustedTotalBags_NoLoss(t *testing.T) {
	h := &Harvest{TotalBags: 100, FieldLoss: 0}
	if h.AdjustedTotalBags() != 100 {
		t.Error("sem perda, total ajustado deve ser igual ao total")
	}
}

// ── Field tests ────────────────────────────────────────────────────────────

func TestField_IsActive(t *testing.T) {
	f := &Field{Status: FieldStatusActive}
	if !f.IsActive() {
		t.Error("campo ativo deveria retornar true")
	}
	f.Status = FieldStatusInactive
	if f.IsActive() {
		t.Error("campo inativo deveria retornar false")
	}
}

func TestField_IsAvailableForPlanting(t *testing.T) {
	f := &Field{Status: FieldStatusActive}
	if !f.IsAvailableForPlanting() {
		t.Error("campo ativo deveria estar disponível para plantio")
	}
	f.Status = FieldStatusFallow
	if !f.IsAvailableForPlanting() {
		t.Error("campo em pousio deveria estar disponível para plantio")
	}
	f.Status = FieldStatusInactive
	if f.IsAvailableForPlanting() {
		t.Error("campo inativo não deveria estar disponível para plantio")
	}
}

func TestField_AreaInSquareMeters(t *testing.T) {
	f := &Field{Area: 2.5}
	if f.AreaInSquareMeters() != 25000 {
		t.Errorf("esperava 25000 m², got %.2f", f.AreaInSquareMeters())
	}
}

func TestField_Deactivate(t *testing.T) {
	f := &Field{Status: FieldStatusActive}
	f.Deactivate()
	if f.Status != FieldStatusInactive {
		t.Error("campo deveria estar inativo após Deactivate")
	}
}

func TestField_SetFallow(t *testing.T) {
	f := &Field{Status: FieldStatusActive}
	f.SetFallow()
	if f.Status != FieldStatusFallow {
		t.Error("campo deveria estar em pousio após SetFallow")
	}
}

// ── Season tests ───────────────────────────────────────────────────────────

func TestSeason_IsActive(t *testing.T) {
	s := &Season{Status: SeasonStatusActive}
	if !s.IsActive() {
		t.Error("safra ativa deveria retornar true")
	}
}

func TestSeason_DurationDays(t *testing.T) {
	s := &Season{
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 6, 0),
	}
	days := s.DurationDays()
	if days < 180 || days > 184 {
		t.Errorf("esperava ~180 dias, got %d", days)
	}
}

func TestSeason_IsOngoing(t *testing.T) {
	s := &Season{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now().AddDate(0, 5, 0),
	}
	if !s.IsOngoing() {
		t.Error("safra deveria estar em andamento")
	}
}

func TestSeason_Activate(t *testing.T) {
	s := &Season{Status: SeasonStatusPlanning}
	s.Activate()
	if s.Status != SeasonStatusActive {
		t.Error("safra deveria estar ativa após Activate")
	}
}

func TestSeason_Finish(t *testing.T) {
	s := &Season{Status: SeasonStatusActive}
	s.Finish()
	if s.Status != SeasonStatusFinished {
		t.Error("safra deveria estar encerrada após Finish")
	}
}

// ── Planting tests ─────────────────────────────────────────────────────────

func TestPlanting_IsActive(t *testing.T) {
	p := &Planting{Status: PlantingStatusActive}
	if !p.IsActive() {
		t.Error("plantio ativo deveria retornar true")
	}
}

func TestPlanting_IsLate(t *testing.T) {
	p := &Planting{
		Status:          PlantingStatusActive,
		ExpectedHarvest: time.Now().AddDate(0, 0, -5),
	}
	if !p.IsLate() {
		t.Error("plantio com data passada deveria estar atrasado")
	}
}

func TestPlanting_IsLate_NotActive(t *testing.T) {
	p := &Planting{
		Status:          PlantingStatusHarvested,
		ExpectedHarvest: time.Now().AddDate(0, 0, -5),
	}
	if p.IsLate() {
		t.Error("plantio colhido não deveria ser considerado atrasado")
	}
}

func TestPlanting_MarkHarvested(t *testing.T) {
	p := &Planting{Status: PlantingStatusActive}
	p.MarkHarvested()
	if p.Status != PlantingStatusHarvested {
		t.Error("plantio deveria estar colhido após MarkHarvested")
	}
}

func TestPlanting_TotalSeedsForArea(t *testing.T) {
	p := &Planting{DensityKgHa: 60}
	total := p.TotalSeedsForArea(10)
	if total != 600 {
		t.Errorf("esperava 600 kg de sementes, got %.2f", total)
	}
}

func TestPlanting_TotalSeedsForArea_InvalidDensity(t *testing.T) {
	p := &Planting{DensityKgHa: 0}
	if p.TotalSeedsForArea(10) != 0 {
		t.Error("densidade zero deveria retornar 0")
	}
}

// ── Input tests ────────────────────────────────────────────────────────────

func TestInput_IsLowStock(t *testing.T) {
	i := &Input{StockQty: 5, MinStockQty: 10}
	if !i.IsLowStock() {
		t.Error("estoque abaixo do mínimo deveria ser low stock")
	}
	i.StockQty = 15
	if i.IsLowStock() {
		t.Error("estoque acima do mínimo não deveria ser low stock")
	}
}

func TestInput_IsExpired(t *testing.T) {
	past := time.Now().AddDate(0, 0, -1)
	i := &Input{ExpirationDate: &past}
	if !i.IsExpired() {
		t.Error("insumo com data passada deveria estar vencido")
	}
}

func TestInput_IsExpiringSoon(t *testing.T) {
	soon := time.Now().AddDate(0, 0, 5)
	i := &Input{ExpirationDate: &soon}
	if !i.IsExpiringSoon(10) {
		t.Error("insumo vencendo em 5 dias deveria ser expiring soon em 10 dias")
	}
	if i.IsExpiringSoon(3) {
		t.Error("insumo vencendo em 5 dias não deveria ser expiring soon em 3 dias")
	}
}
