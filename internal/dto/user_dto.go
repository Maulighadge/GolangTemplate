package dto

type UpdateUserRequest struct {
	Name string `json:"name" binding:"omitempty,min=2,max=100"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
