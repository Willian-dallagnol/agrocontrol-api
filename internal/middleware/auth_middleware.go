package middleware

import (
	"agrocontrol-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "user_id"
	ctxEmail  = "email"
	ctxRole   = "role"
)

// AuthMiddleware valida o token JWT e injeta user_id, email e role no contexto Gin.
// Deve ser aplicado ANTES de qualquer handler que precise de autenticação.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token não fornecido",
				"code":  "missing_token",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "formato inválido — use: Authorization: Bearer <token>",
				"code":  "invalid_token_format",
			})
			return
		}

		claims, err := utils.ValidateToken(parts[1], secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido ou expirado",
				"code":  "invalid_token",
			})
			return
		}

		// Injeta dados do usuário no contexto para os handlers
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxEmail, claims.Email)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}
