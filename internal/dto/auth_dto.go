package dto

// 🔐 Estrutura usada para receber os dados de login (entrada da API)
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	// 👉 email do usuário
	// 🔥 obrigatório e deve estar em formato válido de email

	Password string `json:"password" binding:"required"`
	// 👉 senha do usuário
	// 🔥 obrigatório
}

// 🔑 Estrutura de resposta após login bem-sucedido
type LoginResponse struct {
	Token string `json:"token"`
	// 👉 token JWT gerado após autenticação
	// 🔐 será usado para acessar rotas protegidas

	User UserResponse `json:"user"`
	// 👉 dados do usuário autenticado (sem expor senha)
}
