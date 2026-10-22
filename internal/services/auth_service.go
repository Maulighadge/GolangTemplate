package services

import (
	"github.com/Maulighadge/GolangTemplate/internal/dto"
	"github.com/Maulighadge/GolangTemplate/internal/models"
	"github.com/Maulighadge/GolangTemplate/internal/repositories"
	"github.com/Maulighadge/GolangTemplate/pkg/apperror"
	"github.com/Maulighadge/GolangTemplate/pkg/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check existing user")
	}
	if existingUser != nil {
		return nil, apperror.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to hash password")
	}

	newUser := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     models.RoleUser,
		IsActive: true,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, apperror.NewInternal(err, "failed to create user")
	}

	token, err := utils.GenerateToken(newUser.ID, newUser.Email, string(newUser.Role))
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to generate authentication token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(newUser),
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to query user")
	}
	if user == nil {
		return nil, apperror.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, apperror.ErrAccountDeactivated
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, apperror.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to generate authentication token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}
