package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/middlewares"
	"github.com/Maulighadge/GolangTemplate/internal/services"
	"github.com/Maulighadge/GolangTemplate/pkg/response"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile returns current user profile
func (ctl *UserController) GetProfile(c *gin.Context) {
	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	profile, err := ctl.userService.GetProfile(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved", profile)
}

// UpdateProfile updates user profile information
func (ctl *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	updated, err := ctl.userService.UpdateProfile(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "User profile updated", updated)
}

// ChangePassword changes user password
func (ctl *UserController) ChangePassword(c *gin.Context) {
	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := ctl.userService.ChangePassword(userID, req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}
