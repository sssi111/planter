package services

import (
	"context"
	"fmt"

	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/repository"
	"github.com/google/uuid"
)

// PlantService handles plant operations
type PlantService struct {
	plantRepo repository.PlantRepository
}

// NewPlantService creates a new plant service
func NewPlantService(plantRepo repository.PlantRepository) *PlantService {
	return &PlantService{
		plantRepo: plantRepo,
	}
}

// GetAllPlants gets all plants
func (s *PlantService) GetAllPlants(ctx context.Context) ([]*models.Plant, error) {
	plants, err := s.plantRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get plants: %w", err)
	}
	return plants, nil
}

// GetPlant gets a plant by ID
func (s *PlantService) GetPlant(ctx context.Context, plantID uuid.UUID) (*models.Plant, error) {
	plant, err := s.plantRepo.GetByID(ctx, plantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plant: %w", err)
	}
	return plant, nil
}

// SearchPlants searches for plants by query
func (s *PlantService) SearchPlants(ctx context.Context, query string) ([]*models.Plant, error) {
	plants, err := s.plantRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search plants: %w", err)
	}
	return plants, nil
}

// GetFavoritePlants gets a user's favorite plants
func (s *PlantService) GetFavoritePlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	plants, err := s.plantRepo.GetFavorites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite plants: %w", err)
	}
	return plants, nil
}

// AddToFavorites adds a plant to a user's favorites
func (s *PlantService) AddToFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	// Check if the plant exists
	_, err := s.plantRepo.GetByID(ctx, plantID)
	if err != nil {
		return fmt.Errorf("plant not found: %w", err)
	}

	// Add to favorites
	err = s.plantRepo.AddToFavorites(ctx, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to add plant to favorites: %w", err)
	}
	return nil
}

// RemoveFromFavorites removes a plant from a user's favorites
func (s *PlantService) RemoveFromFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	err := s.plantRepo.RemoveFromFavorites(ctx, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to remove plant from favorites: %w", err)
	}
	return nil
}

// MarkAsWatered marks a plant as watered
func (s *PlantService) MarkAsWatered(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (*models.Plant, error) {
	// Check if the plant exists
	plant, err := s.plantRepo.GetByID(ctx, plantID)
	if err != nil {
		return nil, fmt.Errorf("plant not found: %w", err)
	}

	// Check if the user owns the plant
	userPlant, err := s.plantRepo.GetUserPlant(ctx, userID, plantID)
	if err != nil {
		return nil, fmt.Errorf("user does not own this plant: %w", err)
	}

	// Mark as watered
	err = s.plantRepo.MarkAsWatered(ctx, userID, plantID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark plant as watered: %w", err)
	}

	// Get the updated user plant
	userPlant, err = s.plantRepo.GetUserPlant(ctx, userID, plantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user plant: %w", err)
	}

	// Update the plant with user plant data
	plant.LastWatered = userPlant.LastWatered
	plant.NextWatering = userPlant.NextWatering
	plant.Location = userPlant.Location

	// Check if the plant is a favorite
	isFavorite, err := s.plantRepo.IsFavorite(ctx, userID, plantID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if plant is favorite: %w", err)
	}
	plant.IsFavorite = isFavorite

	return plant, nil
}

// GetUserPlants gets all plants owned by a user
func (s *PlantService) GetUserPlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	plants, err := s.plantRepo.GetUserPlants(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user plants: %w", err)
	}
	return plants, nil
}

// AddUserPlant adds a plant to a user's collection
func (s *PlantService) AddUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID, location string) error {
	// Check if the plant exists
	_, err := s.plantRepo.GetByID(ctx, plantID)
	if err != nil {
		return fmt.Errorf("plant not found: %w", err)
	}

	// Add the plant to the user's collection
	userPlant := &models.UserPlant{
		UserID:   userID,
		PlantID:  plantID,
		Location: &location,
	}

	err = s.plantRepo.AddUserPlant(ctx, userPlant)
	if err != nil {
		return fmt.Errorf("failed to add user plant: %w", err)
	}
	return nil
}

// UpdateUserPlant updates a user's plant
func (s *PlantService) UpdateUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID, location string) error {
	// Check if the user owns the plant
	userPlant, err := s.plantRepo.GetUserPlant(ctx, userID, plantID)
	if err != nil {
		return fmt.Errorf("user does not own this plant: %w", err)
	}

	// Update the location
	userPlant.Location = &location

	// Update the user plant
	err = s.plantRepo.UpdateUserPlant(ctx, userPlant)
	if err != nil {
		return fmt.Errorf("failed to update user plant: %w", err)
	}
	return nil
}

// RemoveUserPlant removes a plant from a user's collection
func (s *PlantService) RemoveUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	err := s.plantRepo.RemoveUserPlant(ctx, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to remove user plant: %w", err)
	}
	return nil
}

// CreatePlant creates a new plant
func (s *PlantService) CreatePlant(ctx context.Context, plant *models.Plant, careInstructions *models.CareInstructions) (*models.Plant, error) {
	// Validate plant data
	if plant.Name == "" {
		return nil, fmt.Errorf("plant name is required")
	}
	if plant.ScientificName == "" {
		return nil, fmt.Errorf("scientific name is required")
	}
	if plant.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if plant.ImageURL == "" {
		return nil, fmt.Errorf("image URL is required")
	}

	// Validate care instructions
	if careInstructions.WateringFrequency <= 0 {
		return nil, fmt.Errorf("watering frequency must be positive")
	}
	if careInstructions.Temperature.Min >= careInstructions.Temperature.Max {
		return nil, fmt.Errorf("minimum temperature must be less than maximum temperature")
	}
	if careInstructions.SoilType == "" {
		return nil, fmt.Errorf("soil type is required")
	}
	if careInstructions.FertilizerFrequency <= 0 {
		return nil, fmt.Errorf("fertilizer frequency must be positive")
	}

	// Create the plant
	createdPlant, err := s.plantRepo.CreatePlant(ctx, plant, careInstructions)
	if err != nil {
		return nil, fmt.Errorf("failed to create plant: %w", err)
	}

	return createdPlant, nil
}