# 🌾 AgroControl API

[![CI](https://github.com/Willian-dallagnol/agrocontrol-api/actions/workflows/ci.yml/badge.svg)](https://github.com/Willian-dallagnol/agrocontrol-api/actions)
[![Coverage](https://img.shields.io/badge/coverage-42.1%25-brightgreen)](https://github.com/Willian-dallagnol/agrocontrol-api)
[![Codecov](https://codecov.io/gh/Willian-dallagnol/agrocontrol-api/branch/main/graph/badge.svg)](https://codecov.io/gh/Willian-dallagnol/agrocontrol-api)
[![Tests](https://img.shields.io/badge/tests-109%20passing-brightgreen)](https://github.com/Willian-dallagnol/agrocontrol-api)
[![Go](https://img.shields.io/badge/Go-1.25-blue?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

API REST completa para gestão de operações agrícolas, desenvolvida em Go com foco em arquitetura limpa, segurança e escalabilidade real.

API em produção: [agrocontrol-api-production.up.railway.app/health](https://agrocontrol-api-production.up.railway.app/health)

---

## Funcionalidades

- Autenticação JWT com access token (curta duração) e refresh token (7 dias)
- RBAC com três papéis: `admin`, `manager`, `operator`
- CRUD completo de fazendas, talhões, culturas, safras, plantios, insumos, aplicações, monitoramentos, colheitas e alertas
- Relatórios de produtividade (sc/ha, kg/ha) e custo por talhão
- Dashboard consolidado com cache Redis (TTL 5 minutos)
- Alertas automáticos para estoque baixo e monitoramentos urgentes
- Transações atômicas para operações críticas (debitação de estoque, registro de colheita)
- Rate limiting por IP (10 req/s, burst 30)
- Request ID rastreável em todos os logs
- Swagger UI disponível em `/swagger/index.html`

---

## Arquitetura
cmd/api/          → entrypoint (main.go)
configs/          → carregamento de configuração e conexão com banco
internal/
apperrors/      → erros de domínio com sentinel errors
cache/          → cliente Redis com helpers tipados
database/       → runner de migrations SQL
domain/
entities/     → modelos de domínio (sem dependência de infra)
ports/        → interfaces de repositório (inversão de dependência)
dto/            → request/response structs com tags de validação
handler/        → handlers HTTP (só conhecem serviços)
middleware/     → auth JWT, logger com request_id, rate limiter, RBAC
mocks/          → mocks de repositório para testes unitários
repository/     → implementações concretas das interfaces (GORM)
routes/         → registro de rotas
service/        → regras de negócio (dependem de ports, não de repository)
tests/          → testes de serviço e handler
migrations/       → SQL versionado (golang-migrate)
prometheus/       → configuração de scraping

### Decisões técnicas documentadas

**Por que interfaces de repositório (`domain/ports`)?**
Serviços dependem de interfaces, não de structs concretas de GORM. Isso permite testar serviços com mocks sem banco, e facilita trocar o ORM no futuro sem tocar na lógica de negócio. Ver [ADR 001](docs/adr/001-ports-and-adapters.md).

**Por que `TxRunner` em vez de passar `*gorm.DB` para os serviços?**
O serviço não deve saber que existe um banco de dados relacional — isso é detalhe de infraestrutura. O `TxRunner` expõe apenas a operação `RunInTx(fn)`, mantendo o GORM confinado na camada de repositório. Ver [ADR 002](docs/adr/002-txrunner.md).

**Por que access token de curta duração + refresh token?**
Access tokens de 24h são um risco de segurança. Com access de 1h e refresh de 7 dias, o impacto de um vazamento é limitado sem degradar a experiência do usuário. Ver [ADR 003](docs/adr/003-refresh-token.md).

**Por que Redis para cache do dashboard?**
O dashboard agrega dados de múltiplas tabelas em queries complexas. O cache de 5 minutos elimina esse custo para a maioria dos requests sem comprometer a atualidade dos dados.

---

## Testes

O projeto possui **112 testes** cobrindo services e entidades de domínio, com **42.1% de cobertura** nos services e **78% nas entidades**.

```bash
# Rodar todos os testes com race detector e cobertura
go test -race -coverprofile=coverage.out ./...

# Ver relatório de cobertura no browser
go tool cover -html=coverage.out

# Rodar apenas os testes de service
go test ./internal/service/... -cover -v
```

Serviços cobertos com testes:
- FarmService — validações de área, RBAC, paginação
- FieldService — validações de fazenda, duplicidade, permissão
- CropService — validações de nome, roles
- SeasonService — validações de datas de safra
- PlantingService — validações de talhão ativo, safra e cultura
- InputService — validações de estoque, custo, ajuste atômico
- AlertService — criação, prioridade padrão, resolução
- HarvestService — transação atômica, status de plantio
- MonitoringService — alertas automáticos para casos críticos
- AuthService — bcrypt e JWT reais
- UserService — validação de roles, email lowercase

**Entidades de domínio com regras de negócio testadas:**
- Harvest — CalculateProductivity, IsHighYield, AdjustedTotalBags
- Field — IsActive, IsAvailableForPlanting, Deactivate, SetFallow
- Season — IsOngoing, DurationDays, Activate, Finish
- Planting — IsLate, DaysUntilHarvest, MarkHarvested, TotalSeedsForArea
- Input — IsLowStock, IsExpired, IsExpiringSoon

---

## Rodando localmente

### Pré-requisitos
- Docker e Docker Compose

### Subindo tudo

```bash
git clone https://github.com/Willian-dallagnol/agrocontrol-api.git
cd agrocontrol-api

cp .env.example .env
# Editar .env e definir JWT_SECRET com no mínimo 32 caracteres

docker compose up --build
```

- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

### Desenvolvimento sem Docker

```bash
docker compose up postgres redis -d
cp .env.example .env
go run ./cmd/api
```

---

## Variáveis de ambiente

| Variável        | Descrição                             | Padrão          |
|-----------------|---------------------------------------|-----------------|
| `APP_PORT`      | Porta da API                          | `8080`          |
| `APP_ENV`       | Ambiente (`development`/`production`) | `production`    |
| `DB_HOST`       | Host do PostgreSQL                    | `localhost`     |
| `DB_PORT`       | Porta do PostgreSQL                   | `5432`          |
| `DB_USER`       | Usuário do banco                      | `postgres`      |
| `DB_PASSWORD`   | Senha do banco                        | —               |
| `DB_NAME`       | Nome do banco                         | `agro_control`  |
| `JWT_SECRET`    | Segredo JWT (mín. 32 chars)           | —               |
| `JWT_EXP_HOURS` | Duração do access token em horas      | `24`            |
| `REDIS_ADDR`    | Endereço do Redis                     | `localhost:6379`|

---

## Stack

| Camada      | Tecnologia                         |
|-------------|-------------------------------------|
| Linguagem   | Go 1.25                             |
| Framework   | Gin                                 |
| ORM         | GORM + driver PostgreSQL            |
| Banco       | PostgreSQL 16                       |
| Cache       | Redis 7                             |
| Auth        | JWT (golang-jwt/jwt v5)             |
| Migrations  | golang-migrate                      |
| Docs        | Swagger (swaggo)                    |
| Observ.     | slog (JSON), Prometheus, Grafana    |
| CI          | GitHub Actions                      |
| Container   | Docker + Docker Compose             |
| Deploy      | Railway                             |

---

## Licença

MIT