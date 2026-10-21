package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/config"
	"github.com/Maulighadge/GolangTemplate/pkg/database"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check returns overall application health status
func (h *HealthController) Check(c *gin.Context) {
	dbStatus := "up"
	sqlDB, err := database.GetDB().DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "down"
	}

	response.Success(c, http.StatusOK, "Service is healthy", gin.H{
		"app":       config.AppConfig.AppName,
		"env":       config.AppConfig.AppEnv,
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"database":  dbStatus,
	})
}
