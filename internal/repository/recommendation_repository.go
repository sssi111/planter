package repository

import (
	"context"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
)

// RecommendationRepository defines the interface for plant recommendation operations
type RecommendationRepository interface {
	// SaveQuestionnaire saves a plant questionnaire
	SaveQuestionnaire(ctx context.Context, questionnaire *models.PlantQuestionnaire) error
	
	// GetQuestionnaire gets a plant questionnaire by ID
	GetQuestionnaire(ctx context.Context, id uuid.UUID) (*models.PlantQuestionnaire, error)
	
	// SaveRecommendation saves a plant recommendation
	SaveRecommendation(ctx context.Context, recommendation *models.PlantRecommendation) error
	
	// GetRecommendations gets all recommendations for a questionnaire
	GetRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.PlantRecommendation, error)
	
	// GetRecommendedPlants gets all recommended plants for a questionnaire
	GetRecommendedPlants(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error)
}