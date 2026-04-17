package middleware

import (
	"agrocontrol-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 🔐 Middleware responsável por validar o token JWT
func AuthMiddleware(secret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 📥 Pega o header Authorization da requisição
		authHeader := c.GetHeader("Authorization")

		// ❌ Se não tiver token
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token não fornecido",
			})
			c.Abort() // 👉 interrompe a requisição
			return
		}

		// 🔍 Divide o header (esperado: "Bearer TOKEN")
		parts := strings.Split(authHeader, " ")

		// ❌ Valida formato do token
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "formato do token inválido",
			})
			c.Abort()
			return
		}

		// 🔑 Valida o token JWT
		claims, err := utils.ValidateToken(parts[1], secret)
		if err != nil {
			// ❌ token inválido ou expirado
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido",
			})
			c.Abort()
			return
		}

		// 📦 Armazena dados do usuário no contexto da requisição
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		// 👉 esses dados podem ser usados nos handlers

		// ✅ continua execução da rota protegida
		c.Next()
	}
}
