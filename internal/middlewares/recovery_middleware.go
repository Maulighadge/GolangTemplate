package middlewares

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/pkg/logger"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
)

// RecoveryMiddleware recovers from panics gracefully, logs the stack trace with request ID, and returns a clean JSON 500
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				reqID := GetRequestID(c)
				stackTrace := string(debug.Stack())
				logger.Error("[ReqID: %s] Panic recovered: %v\nStack Trace:\n%s", reqID, r, stackTrace)

				response.ErrorWithCode(
					c,
					500,
					"An unexpected internal server error occurred",
					"INTERNAL_SERVER_ERROR",
					nil,
				)
				c.Abort()
			}
		}()
		c.Next()
	}
}
