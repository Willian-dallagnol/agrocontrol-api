# ADR 001 — Interfaces de repositório (Ports & Adapters)

**Status:** Aceito  
**Data:** 2025

## Contexto

Os serviços originalmente dependiam diretamente das structs concretas do pacote `repository` (ex: `*repository.FarmRepository`). Isso tornava impossível testar serviços sem um banco de dados real e criava acoplamento desnecessário entre regras de negócio e infraestrutura.

## Decisão

Criamos o pacote `internal/domain/ports` com interfaces Go para todos os repositórios. Os serviços dependem das interfaces, não das implementações concretas.

## Consequências

- ✅ Serviços testáveis com mocks sem banco
- ✅ Possível trocar GORM por qualquer outro ORM sem alterar serviços
- ✅ Inversão de dependência clara: domínio não conhece infraestrutura
- ⚠️ Mais arquivos para manter (interfaces + implementações)
