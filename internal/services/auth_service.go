package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo repository.UserRepository
	auth     *middleware.Auth
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, auth *middleware.Auth) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		auth:     auth,
	}
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	// Get the user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate a token
	token, err := s.auth.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Hide the password hash
	user.PasswordHash = ""

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Register creates a new user and returns a token
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.AuthResponse, error) {
	// Check if the email is already in use
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	user := &models.User{
		Name:                name,
		Email:               email,
		PasswordHash:        string(hashedPassword),
		Language:            models.LanguageRussian,
		NotificationsEnabled: true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate a token
	token, err := s.auth.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Hide the password hash
	user.PasswordHash = ""

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}