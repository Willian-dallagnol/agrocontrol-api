# ADR 002 — TxRunner para transações atômicas

**Status:** Aceito  
**Data:** 2025

## Contexto

Operações críticas (criar aplicação + debitar estoque + criar alerta; registrar colheita + atualizar plantio) precisam ser atômicas. A solução original passava `*gorm.DB` diretamente para os serviços, violando a separação de camadas.

## Decisão

Criamos a interface `ports.TxRunner` com o método `RunInTx(fn func(tx TxRunner) error) error`. A implementação concreta `GormTxRunner` usa `db.Transaction()` do GORM. Os serviços chamam `TxRunner.RunInTx()` sem saber que existe GORM.

## Consequências

- ✅ GORM confinado exclusivamente no pacote `repository`
- ✅ Serviços testáveis sem banco (TxRunner mockável)
- ✅ Transações atômicas preservadas
- ⚠️ Pequeno overhead de abstração (justificado pelo benefício)
