package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/config"
	"github.com/Maulighadge/GolangTemplate/internal/routes"
	"github.com/Maulighadge/GolangTemplate/pkg/database"
	"github.com/Maulighadge/GolangTemplate/pkg/logger"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	logger.Info("Starting application initialization...")

	// 2. Load Configuration
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		logger.Error("Configuration validation error: %v", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully. Environment: %s", cfg.AppEnv)

	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 3. Initialize Database & GORM
	db, err := database.InitDatabase(cfg)
	if err != nil {
		logger.Error("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	// 4. Run Auto Migrations
	if err := database.AutoMigrate(db); err != nil {
		logger.Error("Migration error: %v", err)
		os.Exit(1)
	}

	// 5. Setup Router and Dependencies
	router := routes.SetupRouter(db)

	// 6. Start HTTP Server with Graceful Shutdown
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:           serverAddr,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Run server in background goroutine
	go func() {
		logger.Info("Server listening on port %s (URL: http://localhost:%s)", cfg.AppPort, cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server listen error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server gracefully...")

	// Context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	// Close database connection pool
	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Error closing database connection pool: %v", err)
		} else {
			logger.Info("Database connection pool closed.")
		}
	}

	logger.Info("Server shutdown complete.")
}
