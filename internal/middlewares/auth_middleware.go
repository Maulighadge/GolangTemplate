package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/models"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
	"github.com/Maulighadge/GolangTemplate/pkg/utils"
)

const (
	CtxUserIDKey = "userID"
	CtxEmailKey  = "userEmail"
	CtxRoleKey   = "userRole"
)

// AuthMiddleware verifies the JWT Bearer token and populates user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization header format. Expected 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user data in Gin context
		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxEmailKey, claims.Email)
		c.Set(CtxRoleKey, models.Role(claims.Role))

		c.Next()
	}
}

// RequireRole checks if the authenticated user has any of the required roles
func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(CtxRoleKey)
		if !exists {
			response.Unauthorized(c, "Unauthorized access")
			c.Abort()
			return
		}

		userRole := roleVal.(models.Role)
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "You do not have permission to perform this action")
		c.Abort()
	}
}

// Helper functions to get authenticated user info from context
func GetAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(CtxUserIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func GetAuthUserRole(c *gin.Context) (models.Role, bool) {
	val, exists := c.Get(CtxRoleKey)
	if !exists {
		return "", false
	}
	role, ok := val.(models.Role)
	return role, ok
}
