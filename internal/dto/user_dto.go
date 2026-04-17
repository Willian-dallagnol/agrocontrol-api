package dto

// 👤 Estrutura usada para criar um novo usuário (entrada da API)
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
	// 👉 nome do usuário
	// 🔥 obrigatório

	Email string `json:"email" binding:"required,email"`
	// 👉 email do usuário
	// 🔥 obrigatório e deve estar em formato válido
	// ⚠️ evita dados inválidos já na entrada

	Password string `json:"password" binding:"required,min=6"`
	// 👉 senha do usuário
	// 🔥 obrigatório
	// 🔐 mínimo de 6 caracteres (segurança básica)

	Role string `json:"role" binding:"required"`
	// 👉 nível de acesso do usuário
	// (admin, manager, operator)
	// 🔥 usado para controle de permissões
}

// 📤 Estrutura de resposta da API (dados retornados ao cliente)
type UserResponse struct {
	ID uint `json:"id"`
	// 👉 identificador do usuário

	Name string `json:"name"`
	// 👉 nome do usuário

	Email string `json:"email"`
	// 👉 email

	Role string `json:"role"`
	// 👉 nível de acesso
}
