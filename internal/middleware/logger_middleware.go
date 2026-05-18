package middleware

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

// Logger registra cada request com método, path, status, latência, IP, user_id e request_id.
// O request_id permite correlacionar múltiplos logs de uma mesma requisição em produção.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userID := c.GetUint(ctxUserID)
		requestID := c.GetString(RequestIDKey)

		logFn := logger.Info
		if status >= 500 {
			logFn = logger.Error
		} else if status >= 400 {
			logFn = logger.Warn
		}

		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		logFn("http",
			"request_id", requestID,
			"method", method,
			"path", fullPath,
			"status", status,
			"latency_ms", fmt.Sprintf("%.2f", float64(latency.Microseconds())/1000),
			"ip", clientIP,
			"user_id", userID,
			"bytes_out", c.Writer.Size(),
		)
	}
}
