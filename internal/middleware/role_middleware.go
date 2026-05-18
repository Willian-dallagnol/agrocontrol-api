package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole bloqueia o acesso se o usuário não possui uma das roles permitidas.
// Deve ser usado APÓS AuthMiddleware.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role := c.GetString(ctxRole)
		if _, ok := allowed[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "você não tem permissão para realizar esta ação",
				"code":  "forbidden",
			})
			return
		}
		c.Next()
	}
}

// AdminOnly permite acesso apenas a usuários com role "admin"
func AdminOnly() gin.HandlerFunc {
	return RequireRole("admin")
}

// ManagerOrAbove permite acesso a "admin" e "manager"
func ManagerOrAbove() gin.HandlerFunc {
	return RequireRole("admin", "manager")
}

// AnyRole permite acesso a qualquer role autenticada (admin, manager, operator)
func AnyRole() gin.HandlerFunc {
	return RequireRole("admin", "manager", "operator")
}
