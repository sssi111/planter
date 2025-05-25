package repository

import (
	"context"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user operations
type UserRepository interface {
	// GetByID gets a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	
	// GetByEmail gets a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error
	
	// Update updates a user
	Update(ctx context.Context, user *models.User) error
	
	// GetLocations gets a user's locations
	GetLocations(ctx context.Context, userID uuid.UUID) ([]string, error)
	
	// AddLocation adds a location to a user
	AddLocation(ctx context.Context, userID uuid.UUID, location string) error
	
	// RemoveLocation removes a location from a user
	RemoveLocation(ctx context.Context, userID uuid.UUID, location string) error
	
	// GetFavoritePlantIDs gets a user's favorite plant IDs
	GetFavoritePlantIDs(ctx context.Context, userID uuid.UUID) ([]string, error)
	
	// GetOwnedPlantIDs gets a user's owned plant IDs
	GetOwnedPlantIDs(ctx context.Context, userID uuid.UUID) ([]string, error)
}