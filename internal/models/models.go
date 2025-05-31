package models

import (
	"time"

	"github.com/google/uuid"
)

// SunlightLevel represents the amount of sunlight a plant needs
type SunlightLevel string

const (
	SunlightLevelLow    SunlightLevel = "LOW"
	SunlightLevelMedium SunlightLevel = "MEDIUM"
	SunlightLevelHigh   SunlightLevel = "HIGH"
)

// HumidityLevel represents the humidity level a plant needs
type HumidityLevel string

const (
	HumidityLevelLow    HumidityLevel = "LOW"
	HumidityLevelMedium HumidityLevel = "MEDIUM"
	HumidityLevelHigh   HumidityLevel = "HIGH"
)

// Language represents the user's preferred language
type Language string

const (
	LanguageRussian Language = "RUSSIAN"
	LanguageEnglish Language = "ENGLISH"
)

// User represents a user in the system
type User struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	Email               string    `json:"email" db:"email"`
	PasswordHash        string    `json:"-" db:"password_hash"`
	ProfileImageURL     *string   `json:"profileImageUrl,omitempty" db:"profile_image_url"`
	Language            Language  `json:"language" db:"language"`
	NotificationsEnabled bool      `json:"notificationsEnabled" db:"notifications_enabled"`
	Locations           []string  `json:"locations,omitempty" db:"-"`
	FavoritePlantIDs    []string  `json:"favoritePlantIds,omitempty" db:"-"`
	OwnedPlantIDs       []string  `json:"ownedPlantIds,omitempty" db:"-"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

// UserLocation represents a location associated with a user
type UserLocation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Location  string    `json:"location" db:"location"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// TemperatureRange represents the min and max temperature for a plant
type TemperatureRange struct {
	Min int `json:"min" db:"min_temperature"`
	Max int `json:"max" db:"max_temperature"`
}

// CareInstructions represents care instructions for a plant
type CareInstructions struct {
	ID                 uuid.UUID     `json:"id" db:"id"`
	WateringFrequency  int           `json:"wateringFrequency" db:"watering_frequency"`
	Sunlight           SunlightLevel `json:"sunlight" db:"sunlight"`
	Temperature        TemperatureRange `json:"temperature" db:"-"`
	Humidity           HumidityLevel `json:"humidity" db:"humidity"`
	SoilType           string        `json:"soilType" db:"soil_type"`
	FertilizerFrequency int           `json:"fertilizerFrequency" db:"fertilizer_frequency"`
	AdditionalNotes    string        `json:"additionalNotes" db:"additional_notes"`
	CreatedAt          time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time     `json:"updatedAt" db:"updated_at"`
}

// Plant represents a plant in the system
type Plant struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	Name             string          `json:"name" db:"name"`
	ScientificName   string          `json:"scientificName" db:"scientific_name"`
	Description      string          `json:"description" db:"description"`
	ImageURL         string          `json:"imageUrl" db:"image_url"`
	CareInstructions CareInstructions `json:"careInstructions" db:"-"`
	Price            *float64        `json:"price,omitempty" db:"price"`
	ShopID           *string         `json:"shopId,omitempty" db:"shop_id"`
	IsFavorite       bool            `json:"isFavorite" db:"-"`
	Location         *string         `json:"location,omitempty" db:"-"`
	LastWatered      *time.Time      `json:"lastWatered,omitempty" db:"-"`
	NextWatering     *time.Time      `json:"nextWatering,omitempty" db:"-"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
}

// UserPlant represents a plant owned by a user
type UserPlant struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"userId" db:"user_id"`
	PlantID      uuid.UUID  `json:"plantId" db:"plant_id"`
	Location     *string    `json:"location,omitempty" db:"location"`
	LastWatered  *time.Time `json:"lastWatered,omitempty" db:"last_watered"`
	NextWatering *time.Time `json:"nextWatering,omitempty" db:"next_watering"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

// UserFavoritePlant represents a plant favorited by a user
type UserFavoritePlant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	PlantID   uuid.UUID `json:"plantId" db:"plant_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// Shop represents a shop in the system
type Shop struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	Rating    float64   `json:"rating" db:"rating"`
	ImageURL  *string   `json:"imageUrl,omitempty" db:"image_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// ShopPlant represents a plant sold by a shop
type ShopPlant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ShopID    uuid.UUID `json:"shopId" db:"shop_id"`
	PlantID   uuid.UUID `json:"plantId" db:"plant_id"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// SpecialOffer represents a special offer in the system
type SpecialOffer struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Title             string    `json:"title" db:"title"`
	Description       string    `json:"description" db:"description"`
	ImageURL          string    `json:"imageUrl" db:"image_url"`
	DiscountPercentage int       `json:"discountPercentage" db:"discount_percentage"`
	ValidUntil        time.Time `json:"validUntil" db:"valid_until"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

// PlantQuestionnaire represents a questionnaire for plant recommendations
type PlantQuestionnaire struct {
	ID                   uuid.UUID     `json:"id" db:"id"`
	UserID               *uuid.UUID    `json:"userId,omitempty" db:"user_id"`
	SunlightPreference   SunlightLevel `json:"sunlightPreference" db:"sunlight_preference"`
	PetFriendly          bool          `json:"petFriendly" db:"pet_friendly"`
	CareLevel            int           `json:"careLevel" db:"care_level"`
	PreferredLocation    *string       `json:"preferredLocation,omitempty" db:"preferred_location"`
	AdditionalPreferences *string       `json:"additionalPreferences,omitempty" db:"additional_preferences"`
	CreatedAt            time.Time     `json:"createdAt" db:"created_at"`
}

// PlantRecommendation represents a plant recommendation based on a questionnaire
type PlantRecommendation struct {
	ID              uuid.UUID `json:"id" db:"id"`
	QuestionnaireID uuid.UUID `json:"questionnaireId" db:"questionnaire_id"`
	PlantID         uuid.UUID `json:"plantId" db:"plant_id"`
	Score           float64   `json:"score" db:"score"`
	Reasoning       string    `json:"reasoning" db:"reasoning"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// QuestionnaireRequest represents a plant questionnaire request
type QuestionnaireRequest struct {
	SunlightPreference   SunlightLevel `json:"sunlightPreference" validate:"required,oneof=LOW MEDIUM HIGH"`
	PetFriendly          bool          `json:"petFriendly"`
	CareLevel            int           `json:"careLevel" validate:"required,min=1,max=5"`
	PreferredLocation    *string       `json:"preferredLocation,omitempty"`
	AdditionalPreferences *string       `json:"additionalPreferences,omitempty"`
}

// ChatMessage represents a message in a chat session
type ChatMessage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SessionID uuid.UUID `json:"sessionId" db:"session_id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Role      string    `json:"role" db:"role"` // "user" or "assistant"
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// ChatSession represents a chat session with Yandex GPT
type ChatSession struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"userId" db:"user_id"`
	Title     string     `json:"title" db:"title"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	LastUsed  time.Time  `json:"lastUsed" db:"last_used"`
}

// ChatRequest represents a request to send a message to the chat
type ChatRequest struct {
	Message string `json:"message" validate:"required"`
}

// ChatResponse represents a response from the chat
type ChatResponse struct {
	Message ChatMessage `json:"message"`
}

// DetailedQuestionnaireRequest represents a detailed plant questionnaire request
type DetailedQuestionnaireRequest struct {
	SunlightPreference    SunlightLevel `json:"sunlightPreference" validate:"required,oneof=LOW MEDIUM HIGH"`
	PetFriendly           bool          `json:"petFriendly"`
	CareLevel             int           `json:"careLevel" validate:"required,min=1,max=5"`
	PreferredLocation     *string       `json:"preferredLocation,omitempty"`
	HasChildren           bool          `json:"hasChildren"`
	PlantSize             string        `json:"plantSize" validate:"required,oneof=SMALL MEDIUM LARGE"`
	FloweringPreference   bool          `json:"floweringPreference"`
	AirPurifying          bool          `json:"airPurifying"`
	WateringFrequency     string        `json:"wateringFrequency" validate:"required,oneof=RARE REGULAR FREQUENT"`
	ExperienceLevel       string        `json:"experienceLevel" validate:"required,oneof=BEGINNER INTERMEDIATE ADVANCED"`
	AdditionalPreferences *string       `json:"additionalPreferences,omitempty"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeWatering NotificationType = "WATERING"
)

// Notification represents a notification in the system
type Notification struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	UserID    uuid.UUID        `json:"userId" db:"user_id"`
	PlantID   uuid.UUID        `json:"plantId" db:"plant_id"`
	Type      NotificationType `json:"type" db:"type"`
	Message   string          `json:"message" db:"message"`
	IsRead    bool            `json:"isRead" db:"is_read"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time       `json:"updatedAt" db:"updated_at"`
	// Additional fields for response
	Plant     *Plant          `json:"plant,omitempty" db:"-"`
}

// NotificationResponse represents the response for notifications list
type NotificationResponse struct {
	Notifications []*Notification `json:"notifications"`
	Total         int            `json:"total"`
}