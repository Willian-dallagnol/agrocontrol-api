# 🚜 AgroControl API

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-blue?logo=go" />
  <img src="https://img.shields.io/badge/Gin-Framework-green" />
  <img src="https://img.shields.io/badge/PostgreSQL-Database-blue?logo=postgresql" />
  <img src="https://img.shields.io/badge/Auth-JWT-orange" />
  <img src="https://img.shields.io/badge/status-em%20desenvolvimento-yellow" />
</p>

---

## 📌 Sobre o projeto

O **AgroControl API** é uma API REST desenvolvida em **Go (Golang)** com foco em **gestão agrícola**, aplicando boas práticas de backend e arquitetura em camadas.

O sistema simula um cenário real do agronegócio, permitindo o gerenciamento de usuários, fazendas e talhões, com autenticação segura e estrutura pronta para expansão.

---

## 🎯 Objetivo

Construir uma API:

- 🔐 Segura (JWT)
- 🧱 Organizada (arquitetura em camadas)
- ⚡ Escalável
- 🌾 Baseada em regras reais do agronegócio

---

## 🧠 Arquitetura

```
Handler → Service → Repository → Database
```

### 📂 Estrutura

```
internal/
├── handler      # camada HTTP
├── service      # regras de negócio
├── repository   # acesso ao banco
├── domain       # entidades
├── dto          # contratos da API
├── middleware   # autenticação JWT
└── utils        # helpers
```

---

## ⚙️ Tecnologias

- **Go (Golang)**
- **Gin**
- **GORM**
- **PostgreSQL**
- **JWT**
- **bcrypt**

---

## 🔐 Funcionalidades

### 👤 Usuários
- Cadastro
- Login
- Controle de roles (admin, manager, operator)

### 🔑 Autenticação
- JWT
- Middleware de proteção de rotas

---

## 🚜 Módulo Farm (Fazendas)

CRUD completo:

```
POST   /farms
GET    /farms
GET    /farms/:id
PUT    /farms/:id
DELETE /farms/:id
```

✔ Validação de regra (`total_area > 0`)

---

## 🌱 Módulo Field (Talhões)

CRUD completo:

```
POST   /fields
GET    /fields
GET    /fields/:id
PUT    /fields/:id
DELETE /fields/:id
```

✔ Validação de área  
✔ Validação de relacionamento com Farm  

---

## 🔗 Relacionamento

```
Farm (1) → (N) Field
```

---

## 🚀 Endpoint avançado

### 📌 Listar talhões por fazenda

```
GET /farms/:id/fields
```

✔ Filtra por farm_id  
✔ Valida existência da fazenda  
✔ Retorna apenas dados relacionados  

---

## ▶️ Como rodar

### 1. Clone o projeto

```bash
git clone https://github.com/Willian-dallagnol/agrocontrol-api.git
cd agrocontrol-api
```

---

### 2. Configure o ambiente

Crie um `.env`:

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=agro_control
JWT_SECRET=supersecret
```

---

### 3. Crie o banco

```sql
CREATE DATABASE agro_control;
```

---

### 4. Execute

```bash
go run cmd/api/main.go
```

---

### 5. Teste

```
http://localhost:8080/health
```

---

## 📈 Roadmap

- [x] Farm
- [x] Field
- [ ] Crop
- [ ] Season
- [ ] Planting
- [ ] Insumos
- [ ] Aplicações agrícolas
- [ ] Alertas
- [ ] Dashboard

---

## 💡 Diferenciais

- 🧠 Arquitetura limpa e organizada
- 🔐 Autenticação com JWT
- ⚙️ Regras de negócio bem definidas
- 🔗 Relacionamento entre entidades
- 🚀 Endpoint com filtro (nível mercado)
- 📦 Estrutura pronta para escalar

---

## 🧑‍💻 Autor

**Willian Dall Agnol**

- GitHub: https://github.com/Willian-dallagnol  
- LinkedIn: https://www.linkedin.com/in/willian-dall-agnol-52161315a/  

---

## 📌 Status

🚧 Em evolução contínua — rumo a uma plataforma completa de gestão agrícola