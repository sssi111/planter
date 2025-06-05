package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/repository"
	"github.com/google/uuid"
)

// YandexGPTRequest represents a request to the Yandex GPT API
type YandexGPTRequest struct {
	ModelURI    string              `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages    []Message           `json:"messages"`
}

// CompletionOptions represents the completion options for the Yandex GPT API
type CompletionOptions struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

// Message represents a message in the Yandex GPT API request
type Message struct {
	Role    string `json:"role"`
	Text    string `json:"text"`
}

// YandexGPTResponse represents a response from the Yandex GPT API
type YandexGPTResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Role    string `json:"role"`
				Text    string `json:"text"`
			} `json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

// RecommendationService handles plant recommendation operations
type RecommendationService struct {
	recommendationRepo repository.RecommendationRepository
	plantRepo          repository.PlantRepository
	yandexGPTAPIKey    string
	yandexGPTModel     string
	chatSessions       map[uuid.UUID][]Message // In-memory cache for chat sessions
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(
	recommendationRepo repository.RecommendationRepository,
	plantRepo repository.PlantRepository,
	yandexGPTAPIKey string,
	yandexGPTModel string,
) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		plantRepo:          plantRepo,
		yandexGPTAPIKey:    yandexGPTAPIKey,
		yandexGPTModel:     yandexGPTModel,
		chatSessions:       make(map[uuid.UUID][]Message),
	}
}

// SaveQuestionnaire saves a plant questionnaire
func (s *RecommendationService) SaveQuestionnaire(ctx context.Context, userID *uuid.UUID, questionnaire *models.QuestionnaireRequest) (*models.PlantQuestionnaire, error) {
	// Create the questionnaire
	plantQuestionnaire := &models.PlantQuestionnaire{
		UserID:               userID,
		SunlightPreference:   questionnaire.SunlightPreference,
		PetFriendly:          questionnaire.PetFriendly,
		CareLevel:            questionnaire.CareLevel,
		PreferredLocation:    questionnaire.PreferredLocation,
		AdditionalPreferences: questionnaire.AdditionalPreferences,
	}

	// Save the questionnaire
	err := s.recommendationRepo.SaveQuestionnaire(ctx, plantQuestionnaire)
	if err != nil {
		return nil, fmt.Errorf("failed to save questionnaire: %w", err)
	}

	return plantQuestionnaire, nil
}

// SaveDetailedQuestionnaire saves a detailed plant questionnaire and generates recommendations
func (s *RecommendationService) SaveDetailedQuestionnaire(
	ctx context.Context, 
	userID *uuid.UUID, 
	questionnaire *models.DetailedQuestionnaireRequest,
) (*models.PlantQuestionnaire, error) {
	// Convert detailed questionnaire to standard questionnaire
	plantQuestionnaire := &models.PlantQuestionnaire{
		UserID:               userID,
		SunlightPreference:   questionnaire.SunlightPreference,
		PetFriendly:          questionnaire.PetFriendly,
		CareLevel:            questionnaire.CareLevel,
		PreferredLocation:    questionnaire.PreferredLocation,
	}

	// Create additional preferences text that includes all the detailed information
	additionalPrefs := fmt.Sprintf(
		"Размер растения: %s, Цветущее: %t, Очищающее воздух: %t, Частота полива: %s, Уровень опыта: %s, Есть дети: %t",
		questionnaire.PlantSize,
		questionnaire.FloweringPreference,
		questionnaire.AirPurifying,
		questionnaire.WateringFrequency,
		questionnaire.ExperienceLevel,
		questionnaire.HasChildren,
	)

	if questionnaire.AdditionalPreferences != nil {
		additionalPrefs += ", " + *questionnaire.AdditionalPreferences
	}

	plantQuestionnaire.AdditionalPreferences = &additionalPrefs

	// Save the questionnaire
	err := s.recommendationRepo.SaveQuestionnaire(ctx, plantQuestionnaire)
	if err != nil {
		return nil, fmt.Errorf("failed to save questionnaire: %w", err)
	}

	return plantQuestionnaire, nil
}

// generateLocalRecommendations generates plant recommendations using local matching logic
func (s *RecommendationService) generateLocalRecommendations(
	ctx context.Context,
	questionnaire *models.PlantQuestionnaire,
	allPlants []*models.Plant,
) ([]*models.PlantRecommendation, error) {
	var recommendations []*models.PlantRecommendation

	for _, plant := range allPlants {
		score := 0.0
		reasoning := ""

		// Match sunlight preference
		if plant.CareInstructions.Sunlight == questionnaire.SunlightPreference {
			score += 0.4
			reasoning += fmt.Sprintf("Уровень освещенности (%s) полностью соответствует вашим требованиям. ", plant.CareInstructions.Sunlight)
		} else if (plant.CareInstructions.Sunlight == models.SunlightLevelMedium && 
			(questionnaire.SunlightPreference == models.SunlightLevelLow || questionnaire.SunlightPreference == models.SunlightLevelHigh)) ||
			((plant.CareInstructions.Sunlight == models.SunlightLevelLow || plant.CareInstructions.Sunlight == models.SunlightLevelHigh) && 
			questionnaire.SunlightPreference == models.SunlightLevelMedium) {
			score += 0.2
			reasoning += fmt.Sprintf("Уровень освещенности (%s) частично соответствует вашим требованиям. ", plant.CareInstructions.Sunlight)
		}

		// Match care level (1-5 scale)
		careLevelDiff := float64(abs(plant.CareInstructions.FertilizerFrequency - questionnaire.CareLevel))
		if careLevelDiff == 0 {
			score += 0.3
			reasoning += "Уровень ухода полностью соответствует вашим возможностям. "
		} else if careLevelDiff == 1 {
			score += 0.15
			reasoning += "Уровень ухода близок к желаемому. "
		}

		// Match pet friendly requirement
		if questionnaire.PetFriendly {
			// For now, assume all plants are not pet friendly unless explicitly marked
			// This would need to be added to the plant model and database
			score += 0.1
			reasoning += "Растение безопасно для домашних животных. "
		}

		// Add location matching if specified
		if questionnaire.PreferredLocation != nil && plant.CareInstructions.AdditionalNotes != "" {
			if strings.Contains(strings.ToLower(plant.CareInstructions.AdditionalNotes), 
				strings.ToLower(*questionnaire.PreferredLocation)) {
				score += 0.2
				reasoning += fmt.Sprintf("Подходит для размещения в %s. ", *questionnaire.PreferredLocation)
			}
		}

		// Create recommendation if score is above threshold
		if score > 0.3 { // Minimum 30% match
			recommendations = append(recommendations, &models.PlantRecommendation{
				QuestionnaireID: questionnaire.ID,
				PlantID:        plant.ID,
				Score:          score,
				Reasoning:      strings.TrimSpace(reasoning),
			})
		}
	}

	// Sort recommendations by score in descending order
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Return top 5 recommendations or all if less than 5
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations, nil
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GenerateRecommendations generates plant recommendations based on a questionnaire
func (s *RecommendationService) GenerateRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error) {
	// Get the questionnaire
	questionnaire, err := s.recommendationRepo.GetQuestionnaire(ctx, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questionnaire: %w", err)
	}

	// Get all plants
	allPlants, err := s.plantRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get plants: %w", err)
	}

	var recommendations []*models.PlantRecommendation
	
	// Try to use Yandex GPT if API key is available
	if s.yandexGPTAPIKey != "" {
		recommendations, err = s.generateRecommendationsWithYandexGPT(ctx, questionnaire, allPlants)
		if err != nil {
			// Fallback to local recommendations if Yandex GPT fails
			recommendations, err = s.generateLocalRecommendations(ctx, questionnaire, allPlants)
			if err != nil {
				return nil, fmt.Errorf("failed to generate recommendations: %w", err)
			}
		}
	} else {
		// Use local recommendations if no API key
		recommendations, err = s.generateLocalRecommendations(ctx, questionnaire, allPlants)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recommendations: %w", err)
		}
	}

	// Save the recommendations
	for _, recommendation := range recommendations {
		err = s.recommendationRepo.SaveRecommendation(ctx, recommendation)
		if err != nil {
			return nil, fmt.Errorf("failed to save recommendation: %w", err)
		}
	}

	// Get the recommended plants
	recommendedPlants, err := s.recommendationRepo.GetRecommendedPlants(ctx, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended plants: %w", err)
	}

	return recommendedPlants, nil
}

// GetRecommendations gets all recommendations for a questionnaire
func (s *RecommendationService) GetRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error) {
	// Check if recommendations exist
	recommendations, err := s.recommendationRepo.GetRecommendations(ctx, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	// If no recommendations exist, generate them
	if len(recommendations) == 0 {
		return s.GenerateRecommendations(ctx, questionnaireID)
	}

	// Get the recommended plants
	recommendedPlants, err := s.recommendationRepo.GetRecommendedPlants(ctx, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended plants: %w", err)
	}

	return recommendedPlants, nil
}

// generateRecommendationsWithYandexGPT generates plant recommendations using Yandex GPT
func (s *RecommendationService) generateRecommendationsWithYandexGPT(
	ctx context.Context,
	questionnaire *models.PlantQuestionnaire,
	allPlants []*models.Plant,
) ([]*models.PlantRecommendation, error) {
	// Prepare the prompt
	prompt := s.preparePrompt(questionnaire, allPlants)

	// Call Yandex GPT API
	response, err := s.callYandexGPTAPI(ctx, prompt, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Yandex GPT API: %w", err)
	}

	// Parse the response
	recommendations, err := s.parseYandexGPTResponse(response, questionnaire.ID, allPlants)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Yandex GPT response: %w", err)
	}

	return recommendations, nil
}

// preparePrompt prepares the prompt for Yandex GPT
func (s *RecommendationService) preparePrompt(questionnaire *models.PlantQuestionnaire, allPlants []*models.Plant) string {
	// Convert sunlight preference to Russian
	sunlightRussian := "средний"
	switch questionnaire.SunlightPreference {
	case models.SunlightLevelLow:
		sunlightRussian = "низкий"
	case models.SunlightLevelMedium:
		sunlightRussian = "средний"
	case models.SunlightLevelHigh:
		sunlightRussian = "высокий"
	}

	// Convert pet friendly to Russian
	petFriendlyRussian := "нет"
	if questionnaire.PetFriendly {
		petFriendlyRussian = "да"
	}

	// Convert care level to Russian
	careLevelRussian := "средний"
	switch questionnaire.CareLevel {
	case 1:
		careLevelRussian = "очень низкий"
	case 2:
		careLevelRussian = "низкий"
	case 3:
		careLevelRussian = "средний"
	case 4:
		careLevelRussian = "высокий"
	case 5:
		careLevelRussian = "очень высокий"
	}

	// Prepare the plant list
	var plantList string
	for i, plant := range allPlants {
		if i > 0 {
			plantList += "\n"
		}
		plantList += fmt.Sprintf("%d. %s (научное название: %s)", i+1, plant.Name, plant.ScientificName)
	}

	// Prepare the prompt
	prompt := fmt.Sprintf(`Ты - эксперт по растениям. Помоги подобрать растения для пользователя на основе его предпочтений.

Предпочтения пользователя:
- Уровень освещенности: %s
- Безопасно для животных: %s
- Уровень ухода: %s
`, sunlightRussian, petFriendlyRussian, careLevelRussian)

	if questionnaire.PreferredLocation != nil {
		prompt += fmt.Sprintf("- Предпочтительное расположение: %s\n", *questionnaire.PreferredLocation)
	}

	if questionnaire.AdditionalPreferences != nil {
		prompt += fmt.Sprintf("- Дополнительные предпочтения: %s\n", *questionnaire.AdditionalPreferences)
	}

	prompt += fmt.Sprintf(`
Список доступных растений:
%s

Выбери 5 наиболее подходящих растений из списка и объясни, почему они подходят пользователю. Для каждого растения укажи его номер из списка, название и оценку соответствия от 0 до 1, где 1 - идеальное соответствие.

Формат ответа:
1. [Номер растения]. [Название растения] - [Оценка]
[Объяснение, почему это растение подходит]

2. [Номер растения]. [Название растения] - [Оценка]
[Объяснение, почему это растение подходит]

и так далее.`, plantList)

	return prompt
}

// callYandexGPTAPI calls the Yandex GPT API with a prompt or messages
func (s *RecommendationService) callYandexGPTAPI(ctx context.Context, prompt string, messages []Message) (string, error) {
	// Prepare the request
	requestBody := YandexGPTRequest{
		ModelURI: s.yandexGPTModel,
		CompletionOptions: CompletionOptions{
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	}

	// Use either prompt or messages
	if prompt != "" {
		requestBody.Messages = []Message{
			{
				Role: "user",
				Text: prompt,
			},
		}
	} else if messages != nil {
		requestBody.Messages = messages
	}

	// Convert the request to JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+s.yandexGPTAPIKey)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Parse the response
	var response YandexGPTResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if there are any alternatives
	if len(response.Result.Alternatives) == 0 {
		return "", fmt.Errorf("no alternatives in response")
	}

	// Return the text of the first alternative
	return response.Result.Alternatives[0].Message.Text, nil
}

// parseYandexGPTResponse parses the response from Yandex GPT
func (s *RecommendationService) parseYandexGPTResponse(
	response string,
	questionnaireID uuid.UUID,
	allPlants []*models.Plant,
) ([]*models.PlantRecommendation, error) {
	// Split the response into lines
	lines := bytes.Split([]byte(response), []byte("\n"))

	var recommendations []*models.PlantRecommendation
	var currentPlantNumber int
	var currentScore float64
	var currentReasoning string
	var parsingReasoning bool

	// Parse each line
	for _, line := range lines {
		lineStr := string(line)

		// Skip empty lines
		if len(lineStr) == 0 {
			continue
		}

		// Check if this is a new plant
		var plantNumber int
		var plantName string
		var score float64
		_, err := fmt.Sscanf(lineStr, "%d. %s - %f", &plantNumber, &plantName, &score)
		if err == nil && plantNumber > 0 && plantNumber <= len(allPlants) {
			// If we were parsing a reasoning, save the previous plant
			if parsingReasoning && currentPlantNumber > 0 {
				// Find the plant by number
				if currentPlantNumber <= len(allPlants) {
					plant := allPlants[currentPlantNumber-1]
					recommendations = append(recommendations, &models.PlantRecommendation{
						QuestionnaireID: questionnaireID,
						PlantID:         plant.ID,
						Score:           currentScore,
						Reasoning:       currentReasoning,
					})
				}
			}

			// Start parsing a new plant
			currentPlantNumber = plantNumber
			currentScore = score
			currentReasoning = ""
			parsingReasoning = true
		} else if parsingReasoning {
			// Add to the current reasoning
			if len(currentReasoning) > 0 {
				currentReasoning += "\n"
			}
			currentReasoning += lineStr
		}
	}

	// Save the last plant
	if parsingReasoning && currentPlantNumber > 0 {
		// Find the plant by number
		if currentPlantNumber <= len(allPlants) {
			plant := allPlants[currentPlantNumber-1]
			recommendations = append(recommendations, &models.PlantRecommendation{
				QuestionnaireID: questionnaireID,
				PlantID:         plant.ID,
				Score:           currentScore,
				Reasoning:       currentReasoning,
			})
		}
	}

	// If no recommendations were parsed, return an error
	if len(recommendations) == 0 {
		return nil, fmt.Errorf("failed to parse any recommendations from response")
	}

	return recommendations, nil
}

// CreateChatSession creates a new chat session
func (s *RecommendationService) CreateChatSession(ctx context.Context, userID uuid.UUID) (*models.ChatSession, error) {
	// Create a new chat session
	session, err := s.recommendationRepo.CreateChatSession(ctx, userID, "Разговор о растениях")
	if err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	// Initialize the in-memory session with a system message
	systemMessage := Message{
		Role: "system",
		Text: "Ты - эксперт по растениям. Помогай пользователям с вопросами о выращивании, уходе и выборе растений. Отвечай на русском языке.",
	}
	s.chatSessions[session.ID] = []Message{systemMessage}

	return session, nil
}

// GetChatSession gets a chat session by ID
func (s *RecommendationService) GetChatSession(ctx context.Context, id uuid.UUID) (*models.ChatSession, error) {
	return s.recommendationRepo.GetChatSession(ctx, id)
}

// GetChatSessionsByUser gets all chat sessions for a user
func (s *RecommendationService) GetChatSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.ChatSession, error) {
	return s.recommendationRepo.GetChatSessionsByUser(ctx, userID)
}

// SendChatMessage sends a message to the chat and gets a response
func (s *RecommendationService) SendChatMessage(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	message string,
) (*models.ChatMessage, error) {
	// Get the chat session
	session, err := s.recommendationRepo.GetChatSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}

	// Check if the user owns the session
	if session.UserID != userID {
		return nil, fmt.Errorf("user does not own this chat session")
	}

	// Create and save the user message
	userMessage := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   message,
		CreatedAt: time.Now(),
	}
	
	err = s.recommendationRepo.SaveChatMessage(ctx, userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get all previous messages for context
	dbMessages, err := s.recommendationRepo.GetChatMessages(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}

	// Prepare messages for the API call
	var messages []Message
	
	// Check if we have in-memory session context
	if sessionMessages, ok := s.chatSessions[sessionID]; ok {
		// Use the in-memory session which includes the system message
		messages = sessionMessages
	} else {
		// Initialize with a system message
		messages = []Message{
			{
				Role: "system",
				Text: "Ты - эксперт по растениям. Помогай пользователям с вопросами о выращивании, уходе и выборе растений. Отвечай на русском языке.",
			},
		}
	}

	// Add previous messages from the database (up to the last 10 messages)
	maxMessages := 10
	startIdx := 0
	if len(dbMessages) > maxMessages {
		startIdx = len(dbMessages) - maxMessages
	}
	
	for i := startIdx; i < len(dbMessages); i++ {
		msg := dbMessages[i]
		messages = append(messages, Message{
			Role: msg.Role,
			Text: msg.Content,
		})
	}

	// Add the current user message
	messages = append(messages, Message{
		Role: "user",
		Text: message,
	})

	// Call Yandex GPT API
	response, err := s.callYandexGPTAPI(ctx, "", messages)
	if err != nil {
		return nil, fmt.Errorf("failed to call Yandex GPT API: %w", err)
	}

	// Create and save the assistant message
	assistantMessage := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   response,
		CreatedAt: time.Now(),
	}
	
	err = s.recommendationRepo.SaveChatMessage(ctx, assistantMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// Update the in-memory session
	messages = append(messages, Message{
		Role: "assistant",
		Text: response,
	})
	s.chatSessions[sessionID] = messages

	// Update the last used timestamp
	err = s.recommendationRepo.UpdateChatSessionLastUsed(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to update chat session last used: %w", err)
	}

	return assistantMessage, nil
}

// GetChatMessages gets all messages for a chat session
func (s *RecommendationService) GetChatMessages(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) ([]*models.ChatMessage, error) {
	// Get the chat session
	session, err := s.recommendationRepo.GetChatSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}

	// Check if the user owns the session
	if session.UserID != userID {
		return nil, fmt.Errorf("user does not own this chat session")
	}

	// Get all messages for the session
	return s.recommendationRepo.GetChatMessages(ctx, sessionID)
}