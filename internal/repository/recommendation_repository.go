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
	
	// SaveDetailedQuestionnaire saves a detailed plant questionnaire
	SaveDetailedQuestionnaire(ctx context.Context, questionnaire *models.DetailedQuestionnaireRequest) (*models.PlantQuestionnaire, error)
	
	// CreateChatSession creates a new chat session
	CreateChatSession(ctx context.Context, userID uuid.UUID, title string) (*models.ChatSession, error)
	
	// GetChatSession gets a chat session by ID
	GetChatSession(ctx context.Context, id uuid.UUID) (*models.ChatSession, error)
	
	// GetChatSessionsByUser gets all chat sessions for a user
	GetChatSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.ChatSession, error)
	
	// SaveChatMessage saves a chat message
	SaveChatMessage(ctx context.Context, message *models.ChatMessage) error
	
	// GetChatMessages gets all messages for a chat session
	GetChatMessages(ctx context.Context, sessionID uuid.UUID) ([]*models.ChatMessage, error)
	
	// UpdateChatSessionLastUsed updates the last used timestamp for a chat session
	UpdateChatSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error
}