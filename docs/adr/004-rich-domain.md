# ADR 004 — Domínio Rico vs Domínio Anêmico

**Status:** Aceito  
**Data:** 2026

## Contexto

A versão inicial do projeto concentrava toda a lógica de negócio nos services, deixando as entidades como simples structs de persistência (domínio anêmico). Isso facilitava o desenvolvimento inicial, mas criava dois problemas:

1. Regras de negócio espalhadas em services dificultam reutilização e teste isolado
2. Entidades sem comportamento não expressam o domínio agrícola de forma clara

## Decisão

Enriquecemos as entidades de domínio com métodos que encapsulam regras de negócio reais:

**Harvest:** `CalculateProductivity(areaHa)`, `IsHighYield(threshold)`, `AdjustedTotalBags()`  
**Field:** `IsActive()`, `IsAvailableForPlanting()`, `Deactivate()`, `SetFallow()`  
**Season:** `IsOngoing()`, `DurationDays()`, `Activate()`, `Finish()`  
**Planting:** `IsLate()`, `DaysUntilHarvest()`, `MarkHarvested()`, `TotalSeedsForArea(areaHa)`  
**Input:** `IsLowStock()`, `IsExpired()`, `IsExpiringSoon(days)` (já existia)

## Consequências

- ✅ Regras de negócio testáveis isoladamente (78% de cobertura nas entidades)
- ✅ Entidades expressam o domínio agrícola de forma clara e legível
- ✅ Services ficam mais enxutos — delegam decisões às entidades
- ✅ Reutilização — qualquer camada pode chamar `planting.IsLate()` sem duplicar lógica
- ⚠️ Entidades ainda carregam tags GORM — separação total exigiria mapeamento explícito (trade-off aceito para este contexto)

## Alternativa considerada

Manter tudo nos services e usar funções utilitárias puras. Descartado porque funções soltas não expressam o domínio com a mesma clareza que métodos nas entidades.