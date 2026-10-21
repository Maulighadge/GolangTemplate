package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/services"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary User registration
// @Tags Auth
func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	result, err := ctl.authService.Register(req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", result)
}

// Login godoc
// @Summary User login
// @Tags Auth
func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	result, err := ctl.authService.Login(req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", result)
}
