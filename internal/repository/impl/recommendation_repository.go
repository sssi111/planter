package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anpanovv/planter/internal/db"
	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
)

// RecommendationRepository is the implementation of the recommendation repository
type RecommendationRepository struct {
	db *db.DB
}

// NewRecommendationRepository creates a new recommendation repository
func NewRecommendationRepository(db *db.DB) *RecommendationRepository {
	return &RecommendationRepository{
		db: db,
	}
}

// SaveQuestionnaire saves a plant questionnaire
func (r *RecommendationRepository) SaveQuestionnaire(ctx context.Context, questionnaire *models.PlantQuestionnaire) error {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO plant_questionnaires (user_id, sunlight_preference, pet_friendly, care_level, preferred_location, additional_preferences)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, questionnaire.UserID, questionnaire.SunlightPreference, questionnaire.PetFriendly, questionnaire.CareLevel,
		questionnaire.PreferredLocation, questionnaire.AdditionalPreferences).
		Scan(&questionnaire.ID, &questionnaire.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save questionnaire: %w", err)
	}
	return nil
}

// GetQuestionnaire gets a plant questionnaire by ID
func (r *RecommendationRepository) GetQuestionnaire(ctx context.Context, id uuid.UUID) (*models.PlantQuestionnaire, error) {
	var questionnaire models.PlantQuestionnaire
	err := r.db.GetContext(ctx, &questionnaire, `
		SELECT id, user_id, sunlight_preference, pet_friendly, care_level, preferred_location, additional_preferences, created_at
		FROM plant_questionnaires
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("questionnaire not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get questionnaire: %w", err)
	}
	return &questionnaire, nil
}

// SaveRecommendation saves a plant recommendation
func (r *RecommendationRepository) SaveRecommendation(ctx context.Context, recommendation *models.PlantRecommendation) error {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO plant_recommendations (questionnaire_id, plant_id, score, reasoning)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, recommendation.QuestionnaireID, recommendation.PlantID, recommendation.Score, recommendation.Reasoning).
		Scan(&recommendation.ID, &recommendation.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save recommendation: %w", err)
	}
	return nil
}

// GetRecommendations gets all recommendations for a questionnaire
func (r *RecommendationRepository) GetRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.PlantRecommendation, error) {
	var recommendations []*models.PlantRecommendation
	err := r.db.SelectContext(ctx, &recommendations, `
		SELECT id, questionnaire_id, plant_id, score, reasoning, created_at
		FROM plant_recommendations
		WHERE questionnaire_id = $1
		ORDER BY score DESC
	`, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}
	return recommendations, nil
}

// GetRecommendedPlants gets all recommended plants for a questionnaire
func (r *RecommendationRepository) GetRecommendedPlants(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, p.price, p.shop_id,
			   p.created_at, p.updated_at,
			   c.id as "care_instructions.id", c.watering_frequency as "care_instructions.watering_frequency",
			   c.sunlight as "care_instructions.sunlight", c.min_temperature, c.max_temperature,
			   c.humidity as "care_instructions.humidity", c.soil_type as "care_instructions.soil_type",
			   c.fertilizer_frequency as "care_instructions.fertilizer_frequency",
			   c.additional_notes as "care_instructions.additional_notes",
			   pr.score, pr.reasoning
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		JOIN plant_recommendations pr ON p.id = pr.plant_id
		WHERE pr.questionnaire_id = $1
		ORDER BY pr.score DESC
	`, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended plants: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		var careInstructions models.CareInstructions
		var minTemp, maxTemp int
		var score float64
		var reasoning string

		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
			&plant.Price, &plant.ShopID, &plant.CreatedAt, &plant.UpdatedAt,
			&careInstructions.ID, &careInstructions.WateringFrequency, &careInstructions.Sunlight,
			&minTemp, &maxTemp, &careInstructions.Humidity, &careInstructions.SoilType,
			&careInstructions.FertilizerFrequency, &careInstructions.AdditionalNotes,
			&score, &reasoning,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plant: %w", err)
		}

		careInstructions.Temperature = models.TemperatureRange{
			Min: minTemp,
			Max: maxTemp,
		}
		plant.CareInstructions = careInstructions
		plants = append(plants, &plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating plants: %w", err)
	}

	return plants, nil
}

// CreateChatSession creates a new chat session
func (r *RecommendationRepository) CreateChatSession(ctx context.Context, userID uuid.UUID, title string) (*models.ChatSession, error) {
	session := &models.ChatSession{
		UserID: userID,
		Title:  title,
	}
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO chat_sessions (user_id, title)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at, last_used
	`, session.UserID, session.Title).
		Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt, &session.LastUsed)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	return session, nil
}

// GetChatSession gets a chat session by ID
func (r *RecommendationRepository) GetChatSession(ctx context.Context, id uuid.UUID) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.GetContext(ctx, &session, `
		SELECT id, user_id, title, created_at, updated_at, last_used
		FROM chat_sessions
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("chat session not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}
	return &session, nil
}

// GetChatSessionsByUser gets all chat sessions for a user
func (r *RecommendationRepository) GetChatSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.ChatSession, error) {
	var sessions []*models.ChatSession
	err := r.db.SelectContext(ctx, &sessions, `
		SELECT id, user_id, title, created_at, updated_at, last_used
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY last_used DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat sessions: %w", err)
	}
	return sessions, nil
}

// SaveChatMessage saves a chat message
func (r *RecommendationRepository) SaveChatMessage(ctx context.Context, message *models.ChatMessage) error {
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO chat_messages (session_id, user_id, role, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, message.SessionID, message.UserID, message.Role, message.Content).
		Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save chat message: %w", err)
	}
	return nil
}

// GetChatMessages gets all messages for a chat session
func (r *RecommendationRepository) GetChatMessages(ctx context.Context, sessionID uuid.UUID) ([]*models.ChatMessage, error) {
	var messages []*models.ChatMessage
	err := r.db.SelectContext(ctx, &messages, `
		SELECT id, session_id, user_id, role, content, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	return messages, nil
}

// UpdateChatSessionLastUsed updates the last used timestamp for a chat session
func (r *RecommendationRepository) UpdateChatSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE chat_sessions
		SET last_used = NOW(), updated_at = NOW()
		WHERE id = $1
	`, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update chat session last used: %w", err)
	}
	return nil
}

// SaveDetailedQuestionnaire saves a detailed plant questionnaire
func (r *RecommendationRepository) SaveDetailedQuestionnaire(ctx context.Context, questionnaire *models.DetailedQuestionnaireRequest) (*models.PlantQuestionnaire, error) {
	// This method is not needed as we're using the standard SaveQuestionnaire method
	// with additional preferences text that includes all the detailed information
	return nil, fmt.Errorf("not implemented: use SaveQuestionnaire instead")
}
