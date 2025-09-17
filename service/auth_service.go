package service

import (
	"chatapp/models"
	"chatapp/repository"
	"chatapp/utils"
	"errors"
)

// AuthService handles authentication business logic
type AuthService interface {
	Login(username, password string) (*models.User, string, error)
	GetUserProfile(userID uint) (*models.User, error)
	CreateUser(user *models.User) error
	ValidateToken(token string) (*utils.Claims, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(username, password string) (*models.User, string, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

func (s *authService) GetUserProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *authService) CreateUser(user *models.User) error {
	// Check if username already exists
	existingUser, _ := s.userRepo.GetByUsername(user.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if user.Email != "" {
		existingUser, _ = s.userRepo.GetByEmail(user.Email)
		if existingUser != nil {
			return errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Create user
	return s.userRepo.Create(user)
}

func (s *authService) ValidateToken(token string) (*utils.Claims, error) {
	return utils.ValidateToken(token)
}
