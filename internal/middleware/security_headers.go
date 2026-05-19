package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adiciona headers de segurança HTTP em todas as respostas.
// Protege contra clickjacking, XSS, sniffing de MIME type e exposição de informações.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Impede que a página seja exibida em iframe — proteção contra clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Ativa proteção XSS do browser e bloqueia o request se detectar ataque
		c.Header("X-XSS-Protection", "1; mode=block")

		// Impede que o browser adivinhe o MIME type — evita ataques de sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Remove o header Referrer em requisições cross-origin — protege URLs internas
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Controla quais features do browser podem ser usadas
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Remove o header que expõe que o servidor usa Go/Gin
		c.Header("Server", "")

		// Força HTTPS por 1 ano — só ativo em produção
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy básica — só permite recursos da própria origem
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		c.Next()
	}
}
