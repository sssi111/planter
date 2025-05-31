package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anpanovv/planter/internal/db"
	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
)

// PlantRepository is the implementation of the plant repository
type PlantRepository struct {
	db *db.DB
}

// NewPlantRepository creates a new plant repository
func NewPlantRepository(db *db.DB) *PlantRepository {
	return &PlantRepository{
		db: db,
	}
}

// GetAll gets all plants
func (r *PlantRepository) GetAll(ctx context.Context) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, p.price, p.shop_id,
			   p.created_at, p.updated_at,
			   c.id as "care_instructions.id", c.watering_frequency as "care_instructions.watering_frequency",
			   c.sunlight as "care_instructions.sunlight", c.min_temperature, c.max_temperature,
			   c.humidity as "care_instructions.humidity", c.soil_type as "care_instructions.soil_type",
			   c.fertilizer_frequency as "care_instructions.fertilizer_frequency",
			   c.additional_notes as "care_instructions.additional_notes"
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		ORDER BY p.name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get plants: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		var careInstructions models.CareInstructions
		var minTemp, maxTemp int

		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
			&plant.Price, &plant.ShopID, &plant.CreatedAt, &plant.UpdatedAt,
			&careInstructions.ID, &careInstructions.WateringFrequency, &careInstructions.Sunlight,
			&minTemp, &maxTemp, &careInstructions.Humidity, &careInstructions.SoilType,
			&careInstructions.FertilizerFrequency, &careInstructions.AdditionalNotes,
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

// GetByID gets a plant by ID
func (r *PlantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Plant, error) {
	var plant models.Plant
	var careInstructions models.CareInstructions
	var minTemp, maxTemp int

	err := r.db.QueryRowxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, p.price, p.shop_id,
			   p.created_at, p.updated_at,
			   c.id, c.watering_frequency, c.sunlight, c.min_temperature, c.max_temperature,
			   c.humidity, c.soil_type, c.fertilizer_frequency, c.additional_notes
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		WHERE p.id = $1
	`, id).Scan(
		&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
		&plant.Price, &plant.ShopID, &plant.CreatedAt, &plant.UpdatedAt,
		&careInstructions.ID, &careInstructions.WateringFrequency, &careInstructions.Sunlight,
		&minTemp, &maxTemp, &careInstructions.Humidity, &careInstructions.SoilType,
		&careInstructions.FertilizerFrequency, &careInstructions.AdditionalNotes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("plant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get plant: %w", err)
	}

	careInstructions.Temperature = models.TemperatureRange{
		Min: minTemp,
		Max: maxTemp,
	}
	plant.CareInstructions = careInstructions

	return &plant, nil
}

// Search searches for plants by query
func (r *PlantRepository) Search(ctx context.Context, query string) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, p.price, p.shop_id,
			   p.created_at, p.updated_at,
			   c.id as "care_instructions.id", c.watering_frequency as "care_instructions.watering_frequency",
			   c.sunlight as "care_instructions.sunlight", c.min_temperature, c.max_temperature,
			   c.humidity as "care_instructions.humidity", c.soil_type as "care_instructions.soil_type",
			   c.fertilizer_frequency as "care_instructions.fertilizer_frequency",
			   c.additional_notes as "care_instructions.additional_notes"
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		WHERE p.name ILIKE $1 OR p.scientific_name ILIKE $1 OR p.description ILIKE $1
		ORDER BY p.name
	`, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search plants: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		var careInstructions models.CareInstructions
		var minTemp, maxTemp int

		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
			&plant.Price, &plant.ShopID, &plant.CreatedAt, &plant.UpdatedAt,
			&careInstructions.ID, &careInstructions.WateringFrequency, &careInstructions.Sunlight,
			&minTemp, &maxTemp, &careInstructions.Humidity, &careInstructions.SoilType,
			&careInstructions.FertilizerFrequency, &careInstructions.AdditionalNotes,
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

// GetFavorites gets a user's favorite plants
func (r *PlantRepository) GetFavorites(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, p.price, p.shop_id,
			   p.created_at, p.updated_at,
			   c.id as "care_instructions.id", c.watering_frequency as "care_instructions.watering_frequency",
			   c.sunlight as "care_instructions.sunlight", c.min_temperature, c.max_temperature,
			   c.humidity as "care_instructions.humidity", c.soil_type as "care_instructions.soil_type",
			   c.fertilizer_frequency as "care_instructions.fertilizer_frequency",
			   c.additional_notes as "care_instructions.additional_notes"
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		JOIN user_favorite_plants ufp ON p.id = ufp.plant_id
		WHERE ufp.user_id = $1
		ORDER BY ufp.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite plants: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		var careInstructions models.CareInstructions
		var minTemp, maxTemp int

		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
			&plant.Price, &plant.ShopID, &plant.CreatedAt, &plant.UpdatedAt,
			&careInstructions.ID, &careInstructions.WateringFrequency, &careInstructions.Sunlight,
			&minTemp, &maxTemp, &careInstructions.Humidity, &careInstructions.SoilType,
			&careInstructions.FertilizerFrequency, &careInstructions.AdditionalNotes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plant: %w", err)
		}

		careInstructions.Temperature = models.TemperatureRange{
			Min: minTemp,
			Max: maxTemp,
		}
		plant.CareInstructions = careInstructions
		plant.IsFavorite = true
		plants = append(plants, &plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating plants: %w", err)
	}

	return plants, nil
}

// AddToFavorites adds a plant to a user's favorites
func (r *PlantRepository) AddToFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_favorite_plants (user_id, plant_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, plant_id) DO NOTHING
	`, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to add plant to favorites: %w", err)
	}
	return nil
}

// RemoveFromFavorites removes a plant from a user's favorites
func (r *PlantRepository) RemoveFromFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_favorite_plants
		WHERE user_id = $1 AND plant_id = $2
	`, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to remove plant from favorites: %w", err)
	}
	return nil
}

// MarkAsWatered marks a plant as watered
func (r *PlantRepository) MarkAsWatered(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	// Get the plant's watering frequency
	var wateringFrequency int
	err := r.db.QueryRowContext(ctx, `
		SELECT c.watering_frequency
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		WHERE p.id = $1
	`, plantID).Scan(&wateringFrequency)
	if err != nil {
		return fmt.Errorf("failed to get plant watering frequency: %w", err)
	}

	// Calculate the next watering date
	nextWatering := time.Now().AddDate(0, 0, wateringFrequency)

	// Update the user plant
	_, err = r.db.ExecContext(ctx, `
		UPDATE user_plants
		SET last_watered = NOW(), next_watering = $1, updated_at = NOW()
		WHERE user_id = $2 AND plant_id = $3
	`, nextWatering, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to mark plant as watered: %w", err)
	}

	return nil
}

// GetUserPlant gets a user's plant
func (r *PlantRepository) GetUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (*models.UserPlant, error) {
	var userPlant models.UserPlant
	err := r.db.GetContext(ctx, &userPlant, `
		SELECT id, user_id, plant_id, location, last_watered, next_watering, created_at, updated_at
		FROM user_plants
		WHERE user_id = $1 AND plant_id = $2
	`, userID, plantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user plant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user plant: %w", err)
	}
	return &userPlant, nil
}

// GetUserPlants gets all plants owned by a user
func (r *PlantRepository) GetUserPlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, p.price, p.shop_id,
			   p.created_at, p.updated_at,
			   c.id as "care_instructions.id", c.watering_frequency as "care_instructions.watering_frequency",
			   c.sunlight as "care_instructions.sunlight", c.min_temperature, c.max_temperature,
			   c.humidity as "care_instructions.humidity", c.soil_type as "care_instructions.soil_type",
			   c.fertilizer_frequency as "care_instructions.fertilizer_frequency",
			   c.additional_notes as "care_instructions.additional_notes",
			   up.location, up.last_watered, up.next_watering
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		JOIN user_plants up ON p.id = up.plant_id
		WHERE up.user_id = $1
		ORDER BY up.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user plants: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		var careInstructions models.CareInstructions
		var minTemp, maxTemp int

		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
			&plant.Price, &plant.ShopID, &plant.CreatedAt, &plant.UpdatedAt,
			&careInstructions.ID, &careInstructions.WateringFrequency, &careInstructions.Sunlight,
			&minTemp, &maxTemp, &careInstructions.Humidity, &careInstructions.SoilType,
			&careInstructions.FertilizerFrequency, &careInstructions.AdditionalNotes,
			&plant.Location, &plant.LastWatered, &plant.NextWatering,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plant: %w", err)
		}

		careInstructions.Temperature = models.TemperatureRange{
			Min: minTemp,
			Max: maxTemp,
		}
		plant.CareInstructions = careInstructions

		// Check if the plant is a favorite
		isFavorite, err := r.IsFavorite(ctx, userID, plant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check if plant is favorite: %w", err)
		}
		plant.IsFavorite = isFavorite

		plants = append(plants, &plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating plants: %w", err)
	}

	return plants, nil
}

// AddUserPlant adds a plant to a user's collection
func (r *PlantRepository) AddUserPlant(ctx context.Context, userPlant *models.UserPlant) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_plants (user_id, plant_id, location, last_watered, next_watering)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, plant_id) DO UPDATE
		SET location = $3, last_watered = $4, next_watering = $5, updated_at = NOW()
	`, userPlant.UserID, userPlant.PlantID, userPlant.Location, userPlant.LastWatered, userPlant.NextWatering)
	if err != nil {
		return fmt.Errorf("failed to add user plant: %w", err)
	}
	return nil
}

// UpdateUserPlant updates a user's plant
func (r *PlantRepository) UpdateUserPlant(ctx context.Context, userPlant *models.UserPlant) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_plants
		SET location = $1, last_watered = $2, next_watering = $3, updated_at = NOW()
		WHERE user_id = $4 AND plant_id = $5
	`, userPlant.Location, userPlant.LastWatered, userPlant.NextWatering, userPlant.UserID, userPlant.PlantID)
	if err != nil {
		return fmt.Errorf("failed to update user plant: %w", err)
	}
	return nil
}

// RemoveUserPlant removes a plant from a user's collection
func (r *PlantRepository) RemoveUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_plants
		WHERE user_id = $1 AND plant_id = $2
	`, userID, plantID)
	if err != nil {
		return fmt.Errorf("failed to remove user plant: %w", err)
	}
	return nil
}

// IsFavorite checks if a plant is a favorite of a user
func (r *PlantRepository) IsFavorite(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
		FROM user_favorite_plants
		WHERE user_id = $1 AND plant_id = $2
	`, userID, plantID)
	if err != nil {
		return false, fmt.Errorf("failed to check if plant is favorite: %w", err)
	}
	return count > 0, nil
}

// CreatePlant creates a new plant
func (r *PlantRepository) CreatePlant(ctx context.Context, plant *models.Plant, careInstructions *models.CareInstructions) (*models.Plant, error) {
	// Begin a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create care instructions
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO care_instructions (
			watering_frequency, sunlight, min_temperature, max_temperature,
			humidity, soil_type, fertilizer_frequency, additional_notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`,
		careInstructions.WateringFrequency,
		careInstructions.Sunlight,
		careInstructions.Temperature.Min,
		careInstructions.Temperature.Max,
		careInstructions.Humidity,
		careInstructions.SoilType,
		careInstructions.FertilizerFrequency,
		careInstructions.AdditionalNotes,
	).Scan(
		&careInstructions.ID,
		&careInstructions.CreatedAt,
		&careInstructions.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create care instructions: %w", err)
	}

	// Create plant
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO plants (
			name, scientific_name, description, image_url,
			care_instructions_id, price, shop_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`,
		plant.Name,
		plant.ScientificName,
		plant.Description,
		plant.ImageURL,
		careInstructions.ID,
		plant.Price,
		plant.ShopID,
	).Scan(
		&plant.ID,
		&plant.CreatedAt,
		&plant.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create plant: %w", err)
	}

	// Set care instructions
	plant.CareInstructions = *careInstructions

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return plant, nil
}

// GetAllUserPlantsForWateringCheck gets all user plants that need to be checked for watering
func (r *PlantRepository) GetAllUserPlantsForWateringCheck(ctx context.Context) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url,
			   up.user_id, up.next_watering
		FROM plants p
		JOIN user_plants up ON p.id = up.plant_id
		WHERE up.next_watering IS NOT NULL
		ORDER BY up.next_watering ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get plants for watering check: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description,
			&plant.ImageURL, &plant.UserID, &plant.NextWatering,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plant: %w", err)
		}
		plants = append(plants, &plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating plants: %w", err)
	}

	return plants, nil
}