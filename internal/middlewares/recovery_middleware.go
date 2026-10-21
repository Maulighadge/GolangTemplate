package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/pkg/logger"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
)

// RecoveryMiddleware recovers from panics gracefully and returns a JSON 500
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errStr := fmt.Sprintf("%v", r)
				logger.Error("Panic recovered: %s", errStr)
				response.Error(c, http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
