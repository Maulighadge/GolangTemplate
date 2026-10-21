package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/pkg/logger"
)

// RequestLogger logs incoming HTTP requests with timing and status
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("%3d | %13v | %15s | %-7s %#v",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
