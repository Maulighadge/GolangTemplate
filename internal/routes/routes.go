package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/internal/controllers"
	"github.com/Maulighadge/GolangTemplate/internal/middlewares"
	"github.com/Maulighadge/GolangTemplate/internal/repositories"
	"github.com/Maulighadge/GolangTemplate/internal/services"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
	"gorm.io/gorm"
)

// SetupRouter initializes Gin router, injects dependencies, and configures routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()

	// Global Middlewares
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.RecoveryMiddleware())
	router.Use(middlewares.CORSMiddleware())

	// Fallback handlers for 404 (Route not found) and 405 (Method not allowed)
	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "The requested endpoint does not exist")
	})
	router.NoMethod(func(c *gin.Context) {
		response.ErrorWithCode(c, 405, "Method not allowed", "METHOD_NOT_ALLOWED", nil)
	})

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)

	// Initialize Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)

	// Initialize Controllers
	healthController := controllers.NewHealthController()
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)

	// Health check route
	router.GET("/health", healthController.Check)

	// API v1 group
	apiV1 := router.Group("/api/v1")
	{
		// Authentication Routes
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// User Routes (Protected)
		users := apiV1.Group("/users")
		users.Use(middlewares.AuthMiddleware())
		{
			users.GET("/me", userController.GetProfile)
			users.PUT("/me", userController.UpdateProfile)
			users.PUT("/change-password", userController.ChangePassword)
		}

		// Products Routes
		products := apiV1.Group("/products")
		{
			// Public product endpoints
			products.GET("", productController.List)
			products.GET("/:id", productController.GetByID)

			// Protected product endpoints
			protectedProducts := products.Group("")
			protectedProducts.Use(middlewares.AuthMiddleware())
			{
				protectedProducts.POST("", productController.Create)
				protectedProducts.PUT("/:id", productController.Update)
				protectedProducts.DELETE("/:id", productController.Delete)
			}
		}
	}

	return router
}
