package services

import (
	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/repositories"
	"github.com/Maulighadge/GolangTemplate/pkg/apperror"
	"github.com/Maulighadge/GolangTemplate/pkg/utils"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to query user profile")
	}
	if user == nil {
		return nil, apperror.ErrUserNotFound
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to query user for update")
	}
	if user == nil {
		return nil, apperror.ErrUserNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperror.NewInternal(err, "failed to update user profile")
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *userService) ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return apperror.NewInternal(err, "failed to query user for password change")
	}
	if user == nil {
		return apperror.ErrUserNotFound
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return apperror.ErrIncorrectPassword
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.NewInternal(err, "failed to hash new password")
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return apperror.NewInternal(err, "failed to update password")
	}

	return nil
}
