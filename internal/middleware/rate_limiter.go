package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitor armazena o limiter e o timestamp do último acesso por IP
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter gerencia limiters por IP com limpeza automática
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rps      rate.Limit // requisições por segundo
	burst    int        // rajada máxima permitida
}

// NewRateLimiter cria um RateLimiter com os parâmetros configurados
// rps: requisições por segundo por IP
// burst: máximo de requisições em rajada
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
	// Goroutine que limpa visitors inativos a cada minuto
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupLoop remove visitors que não fazem requisições há mais de 3 minutos
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware retorna o handler Gin do rate limiter
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "muitas requisições — tente novamente em alguns segundos",
				"code":  "rate_limit_exceeded",
			})
			return
		}
		c.Next()
	}
}
