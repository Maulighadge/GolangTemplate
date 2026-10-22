package controllers

import (
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
		response.Unauthorized(c, "Unauthorized access")
		return
	}

	profile, err := ctl.userService.GetProfile(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "User profile retrieved", profile)
}

// UpdateProfile updates user profile information
func (ctl *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized access")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	updated, err := ctl.userService.UpdateProfile(userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "User profile updated", updated)
}

// ChangePassword changes user password
func (ctl *UserController) ChangePassword(c *gin.Context) {
	userID, ok := middlewares.GetAuthUserID(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized access")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := ctl.userService.ChangePassword(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, "Password changed successfully", nil)
}
