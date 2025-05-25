package repository

import (
	"context"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
)

// ShopRepository defines the interface for shop operations
type ShopRepository interface {
	// GetAll gets all shops
	GetAll(ctx context.Context) ([]*models.Shop, error)
	
	// GetByID gets a shop by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error)
	
	// GetPlants gets all plants from a shop
	GetPlants(ctx context.Context, shopID uuid.UUID) ([]*models.Plant, error)
	
	// GetSpecialOffers gets all special offers
	GetSpecialOffers(ctx context.Context) ([]*models.SpecialOffer, error)
}