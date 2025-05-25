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

// ShopRepository is the implementation of the shop repository
type ShopRepository struct {
	db *db.DB
}

// NewShopRepository creates a new shop repository
func NewShopRepository(db *db.DB) *ShopRepository {
	return &ShopRepository{
		db: db,
	}
}

// GetAll gets all shops
func (r *ShopRepository) GetAll(ctx context.Context) ([]*models.Shop, error) {
	var shops []*models.Shop
	err := r.db.SelectContext(ctx, &shops, `
		SELECT id, name, address, rating, image_url, created_at, updated_at
		FROM shops
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get shops: %w", err)
	}
	return shops, nil
}

// GetByID gets a shop by ID
func (r *ShopRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.GetContext(ctx, &shop, `
		SELECT id, name, address, rating, image_url, created_at, updated_at
		FROM shops
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("shop not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return &shop, nil
}

// GetPlants gets all plants from a shop
func (r *ShopRepository) GetPlants(ctx context.Context, shopID uuid.UUID) ([]*models.Plant, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT p.id, p.name, p.scientific_name, p.description, p.image_url, sp.price, sp.shop_id,
			   p.created_at, p.updated_at,
			   c.id as "care_instructions.id", c.watering_frequency as "care_instructions.watering_frequency",
			   c.sunlight as "care_instructions.sunlight", c.min_temperature, c.max_temperature,
			   c.humidity as "care_instructions.humidity", c.soil_type as "care_instructions.soil_type",
			   c.fertilizer_frequency as "care_instructions.fertilizer_frequency",
			   c.additional_notes as "care_instructions.additional_notes"
		FROM plants p
		JOIN care_instructions c ON p.care_instructions_id = c.id
		JOIN shop_plants sp ON p.id = sp.plant_id
		WHERE sp.shop_id = $1
		ORDER BY p.name
	`, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop plants: %w", err)
	}
	defer rows.Close()

	var plants []*models.Plant
	for rows.Next() {
		var plant models.Plant
		var careInstructions models.CareInstructions
		var minTemp, maxTemp int
		var shopIDStr string

		err := rows.Scan(
			&plant.ID, &plant.Name, &plant.ScientificName, &plant.Description, &plant.ImageURL,
			&plant.Price, &shopIDStr, &plant.CreatedAt, &plant.UpdatedAt,
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
		plant.ShopID = &shopIDStr
		plants = append(plants, &plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating plants: %w", err)
	}

	return plants, nil
}

// GetSpecialOffers gets all special offers
func (r *ShopRepository) GetSpecialOffers(ctx context.Context) ([]*models.SpecialOffer, error) {
	var offers []*models.SpecialOffer
	err := r.db.SelectContext(ctx, &offers, `
		SELECT id, title, description, image_url, discount_percentage, valid_until, created_at, updated_at
		FROM special_offers
		WHERE valid_until > NOW()
		ORDER BY valid_until
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get special offers: %w", err)
	}
	return offers, nil
}