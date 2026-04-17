# 🚜 AgroControl API

API de Gestão Agrícola desenvolvida em **Go**, com autenticação via **JWT**, banco de dados **PostgreSQL** e arquitetura em camadas (Handler, Service, Repository).

---

## 📌 Sobre o projeto

O **AgroControl API** é um backend voltado para o gerenciamento de operações agrícolas, permitindo controle de usuários, autenticação segura e estrutura escalável para expansão futura (fazendas, talhões, plantios, insumos e mais).

Este projeto foi desenvolvido com foco em:

* boas práticas de backend
* organização de código
* regras de negócio reais
* uso de tecnologias utilizadas no mercado

---

## 🧠 Arquitetura

O projeto segue o padrão de separação por camadas:

```
Handler → Service → Repository → Database
```

### 📂 Estrutura de pastas

```
agrocontrol-api
├── cmd/api                # ponto de entrada da aplicação
├── configs               # configuração (env, banco)
├── internal
│   ├── domain/entities   # entidades do sistema
│   ├── dto               # entrada/saída de dados
│   ├── repository        # acesso ao banco
│   ├── service           # regras de negócio
│   ├── handler           # camada HTTP
│   ├── middleware        # autenticação
│   └── utils             # helpers (JWT, bcrypt)
├── .env
├── go.mod
```

---

## ⚙️ Tecnologias utilizadas

* **Go (Golang)**
* **Gin** (framework HTTP)
* **GORM** (ORM)
* **PostgreSQL**
* **JWT (JSON Web Token)**
* **bcrypt** (hash de senha)

---

## 🔐 Funcionalidades implementadas

### 👤 Usuários

* Cadastro de usuário
* Validação de dados
* Senha criptografada com bcrypt

### 🔑 Autenticação

* Login com email e senha
* Geração de token JWT
* Expiração de token

### 🛡️ Segurança

* Middleware de autenticação
* Rotas protegidas

---

## 🚀 Endpoints disponíveis

### 🔹 Health Check

```
GET /health
```

---

### 🔹 Criar usuário

```
POST /users
```

**Body:**

```json
{
  "name": "Willian",
  "email": "willian@email.com",
  "password": "123456",
  "role": "admin"
}
```

---

### 🔹 Login

```
POST /login
```

**Body:**

```json
{
  "email": "willian@email.com",
  "password": "123456"
}
```

---

### 🔹 Rota protegida

```
GET /auth/me
```

**Header:**

```
Authorization: Bearer TOKEN
```

---

## ▶️ Como rodar o projeto

### 1. Clonar o repositório

```bash
git clone https://github.com/Willian-dallagnol/agrocontrol-api.git
cd agrocontrol-api
```

---

### 2. Configurar variáveis de ambiente

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

### 3. Criar banco de dados

```sql
CREATE DATABASE agro_control;
```

---

### 4. Rodar a aplicação

```bash
go run cmd/api/main.go
```

---

### 5. Testar

Abra no navegador:

```
http://localhost:8080/health
```

---

## 📈 Próximas funcionalidades

* Cadastro de **Fazendas (Farm)**
* Gestão de **Talhões (Field)**
* Controle de **Culturas (Crop)**
* Gestão de **Safras (Season)**
* Registro de **Plantios**
* Controle de **Insumos**
* Registro de **Aplicações**
* Sistema de **Alertas agrícolas**

---

## 🧑‍💻 Autor

**Willian Dall Agnol**

* GitHub: https://github.com/Willian-dallagnol
* LinkedIn: https://www.linkedin.com/in/willian-dall-agnol-52161315a/

---

## 📌 Status do projeto

🚧 Em desenvolvimento — novas funcionalidades sendo adicionadas continuamente.

---
