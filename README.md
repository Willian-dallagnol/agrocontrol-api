# 🚜 AgroControl API

![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-green)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue?logo=postgresql)
![JWT](https://img.shields.io/badge/Auth-JWT-orange)
![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow)

---

## 📌 Sobre o projeto

O **AgroControl API** é uma API REST desenvolvida em **Go (Golang)** com foco em **gestão agrícola**, aplicando boas práticas de desenvolvimento backend e arquitetura em camadas.

O projeto simula um sistema real de gerenciamento rural, permitindo controle de usuários, autenticação segura e estrutura pronta para expansão com módulos agrícolas.

---

## 🎯 Objetivo

Criar uma API robusta, escalável e organizada, com foco em:

* autenticação segura (JWT)
* organização por camadas (clean architecture simplificada)
* regras de negócio reais do agro
* código pronto para produção e evolução

---

## 🧠 Arquitetura

O projeto segue uma arquitetura baseada em separação de responsabilidades:

```
Handler → Service → Repository → Database
```

### 📂 Estrutura

```
internal/
├── handler      # camada HTTP (entrada)
├── service      # regras de negócio
├── repository   # acesso ao banco
├── domain       # entidades do sistema
├── dto          # contratos de entrada/saída
├── middleware   # autenticação
└── utils        # helpers (JWT, bcrypt)
```

---

## ⚙️ Tecnologias

* **Go (Golang)**
* **Gin** (HTTP framework)
* **GORM** (ORM)
* **PostgreSQL**
* **JWT (auth)**
* **bcrypt (hash de senha)**

---

## 🔐 Funcionalidades implementadas

### 👤 Usuários

* Cadastro de usuários
* Validação de dados
* Controle de roles:

  * `admin`
  * `manager`
  * `operator`

---

### 🔑 Autenticação

* Login com email e senha
* Geração de token JWT
* Token com expiração

---

### 🛡️ Segurança

* Middleware de autenticação
* Rotas protegidas

---

### 🚜 Módulo Farm (Fazendas)

CRUD completo de fazendas com autenticação JWT:

* `POST /farms` → Criar fazenda
* `GET /farms` → Listar fazendas
* `GET /farms/:id` → Buscar por ID
* `PUT /farms/:id` → Atualizar
* `DELETE /farms/:id` → Remover

---

## 🚀 Endpoints principais

### Health Check

```
GET /health
```

---

### Criar usuário

```
POST /users
```

```json
{
  "name": "Willian",
  "email": "willian@email.com",
  "password": "123456",
  "role": "admin"
}
```

---

### Login

```
POST /login
```

```json
{
  "email": "willian@email.com",
  "password": "123456"
}
```

---

### Rota protegida

```
GET /auth/me
```

Header:

```
Authorization: Bearer TOKEN
```

---

## ▶️ Como rodar o projeto

### 1. Clonar

```bash
git clone https://github.com/Willian-dallagnol/agrocontrol-api.git
cd agrocontrol-api
```

---

### 2. Configurar ambiente

Crie um arquivo `.env`:

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

### 3. Banco de dados

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

Acesse:

```
http://localhost:8080/health
```

---

## 📈 Roadmap (próximas features)

* [x] Módulo de Fazendas (Farm)
* [ ] Talhões (Field)
* [ ] Culturas (Crop)
* [ ] Safras (Season)
* [ ] Plantios
* [ ] Insumos
* [ ] Aplicações agrícolas
* [ ] Alertas inteligentes
* [ ] Dashboard

---

## 💡 Diferenciais do projeto

* Estrutura organizada (nível mercado)
* Separação de responsabilidades clara
* Autenticação real com JWT
* Banco relacional com ORM
* CRUD completo implementado
* Base pronta para escalar

---

## 🧑‍💻 Autor

**Willian Dall Agnol**

* GitHub: https://github.com/Willian-dallagnol
* LinkedIn: https://www.linkedin.com/in/willian-dall-agnol-52161315a/

---

## 📌 Status

🚧 Em desenvolvimento contínuo — evoluindo para uma plataforma completa de gestão agrícola.
