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

O sistema simula um cenário real do agronegócio, permitindo o gerenciamento de:

- Usuários
- Fazendas (Farm)
- Talhões (Field)
- Culturas (Crop)

Tudo com autenticação segura via JWT e estrutura pronta para evolução.

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
├── handler
├── service
├── repository
├── domain
├── dto
├── middleware
└── utils
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

---

### 🔑 Autenticação
- JWT
- Middleware de proteção de rotas

---

## 🚜 Módulo Farm (Fazendas)

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

## 🌾 Módulo Crop (Culturas)

```
POST   /crops
GET    /crops
GET    /crops/:id
PUT    /crops/:id
DELETE /crops/:id
```

✔ Validação de relacionamento com Field  
✔ Regra de negócio aplicada  

---

## 🔗 Relacionamento entre entidades

```
Farm (1) → (N) Field → (N) Crop
```

Cada:

- Fazenda possui vários talhões  
- Talhão possui várias culturas  

---

## 🚀 Endpoints avançados

### 📌 Listar talhões por fazenda

```
GET /farms/:id/fields
```

✔ Filtra dados por relacionamento  
✔ Valida existência da fazenda  

---

## ▶️ Como rodar o projeto

### 1. Clonar

```bash
git clone https://github.com/Willian-dallagnol/agrocontrol-api.git
cd agrocontrol-api
```

---

### 2. Configurar ambiente

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

### 3. Criar banco

```sql
CREATE DATABASE agro_control;
```

---

### 4. Rodar aplicação

```bash
go run cmd/api/main.go
```

---

### 5. Testar

```
http://localhost:8080/health
```

---

## 📸 Exemplo de resposta

```json
[
  {
    "id": 1,
    "name": "Soja",
    "type": "Grão",
    "field_id": 1
  }
]
```

---

## 📈 Roadmap

- [x] Farm
- [x] Field
- [x] Crop
- [ ] Season
- [ ] Planting
- [ ] Insumos
- [ ] Aplicações agrícolas
- [ ] Alertas
- [ ] Dashboard

---

## 💡 Diferenciais

- 🧠 Arquitetura limpa
- 🔐 Autenticação com JWT
- ⚙️ Regras de negócio aplicadas
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