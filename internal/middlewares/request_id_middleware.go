package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	HeaderXRequestID = "X-Request-ID"
	CtxRequestIDKey  = "requestID"
)

// RequestIDMiddleware ensures every request has a unique Request ID for tracing and error correlation
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set(CtxRequestIDKey, reqID)
		c.Writer.Header().Set(HeaderXRequestID, reqID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from the Gin context
func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(CtxRequestIDKey); exists {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
