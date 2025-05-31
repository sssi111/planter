package repository

import (
	"context"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
)

// PlantRepository defines the interface for plant operations
type PlantRepository interface {
	// GetAll gets all plants
	GetAll(ctx context.Context) ([]*models.Plant, error)
	
	// GetByID gets a plant by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Plant, error)
	
	// Search searches for plants by query
	Search(ctx context.Context, query string) ([]*models.Plant, error)
	
	// GetFavorites gets a user's favorite plants
	GetFavorites(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error)
	
	// AddToFavorites adds a plant to a user's favorites
	AddToFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error
	
	// RemoveFromFavorites removes a plant from a user's favorites
	RemoveFromFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error
	
	// MarkAsWatered marks a plant as watered
	MarkAsWatered(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error
	
	// GetUserPlant gets a user's plant
	GetUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (*models.UserPlant, error)
	
	// GetUserPlants gets all plants owned by a user
	GetUserPlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error)
	
	// AddUserPlant adds a plant to a user's collection
	AddUserPlant(ctx context.Context, userPlant *models.UserPlant) error
	
	// UpdateUserPlant updates a user's plant
	UpdateUserPlant(ctx context.Context, userPlant *models.UserPlant) error
	
	// RemoveUserPlant removes a plant from a user's collection
	RemoveUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error
	
	// IsFavorite checks if a plant is a favorite of a user
	IsFavorite(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (bool, error)
	
	// CreatePlant creates a new plant
	CreatePlant(ctx context.Context, plant *models.Plant, careInstructions *models.CareInstructions) (*models.Plant, error)
	
	// GetAllUserPlantsForWateringCheck gets all user plants that need to be checked for watering
	GetAllUserPlantsForWateringCheck(ctx context.Context) ([]*models.UserPlant, error)
}