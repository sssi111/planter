package services

import (
	"context"
	"fmt"

	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/repository"
	"github.com/google/uuid"
)

// UserService handles user operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser gets a user by ID
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Get the existing user to ensure it exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update only the allowed fields
	existingUser.Name = user.Name
	existingUser.ProfileImageURL = user.ProfileImageURL
	existingUser.Language = user.Language
	existingUser.NotificationsEnabled = user.NotificationsEnabled
	existingUser.Locations = user.Locations

	// Update the user
	err = s.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}

// AddLocation adds a location to a user
func (s *UserService) AddLocation(ctx context.Context, userID uuid.UUID, location string) error {
	err := s.userRepo.AddLocation(ctx, userID, location)
	if err != nil {
		return fmt.Errorf("failed to add location: %w", err)
	}
	return nil
}

// RemoveLocation removes a location from a user
func (s *UserService) RemoveLocation(ctx context.Context, userID uuid.UUID, location string) error {
	err := s.userRepo.RemoveLocation(ctx, userID, location)
	if err != nil {
		return fmt.Errorf("failed to remove location: %w", err)
	}
	return nil
}

// GetLocations gets a user's locations
func (s *UserService) GetLocations(ctx context.Context, userID uuid.UUID) ([]string, error) {
	locations, err := s.userRepo.GetLocations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}
	return locations, nil
}