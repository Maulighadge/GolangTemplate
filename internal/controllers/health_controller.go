package controllers

import (
	"net/http"

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
	dbInstance := database.GetDB()
	if dbInstance == nil {
		dbStatus = "down"
	} else {
		sqlDB, err := dbInstance.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}
	}

	appName := "Golang-API"
	appEnv := "development"
	if config.AppConfig != nil {
		appName = config.AppConfig.AppName
		appEnv = config.AppConfig.AppEnv
	}

	if dbStatus == "down" {
		response.ErrorWithCode(c, http.StatusServiceUnavailable, "Service is degraded", "SERVICE_DEGRADED", gin.H{
			"app":      appName,
			"env":      appEnv,
			"status":   "degraded",
			"database": dbStatus,
		})
		return
	}

	response.OK(c, "Service is healthy", gin.H{
		"app":      appName,
		"env":      appEnv,
		"status":   "healthy",
		"database": dbStatus,
	})
}
