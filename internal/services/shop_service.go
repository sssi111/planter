package services

import (
	"context"
	"fmt"

	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/repository"
	"github.com/google/uuid"
)

// ShopService handles shop operations
type ShopService struct {
	shopRepo repository.ShopRepository
}

// NewShopService creates a new shop service
func NewShopService(shopRepo repository.ShopRepository) *ShopService {
	return &ShopService{
		shopRepo: shopRepo,
	}
}

// GetAllShops gets all shops
func (s *ShopService) GetAllShops(ctx context.Context) ([]*models.Shop, error) {
	shops, err := s.shopRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shops: %w", err)
	}
	return shops, nil
}

// GetShop gets a shop by ID
func (s *ShopService) GetShop(ctx context.Context, shopID uuid.UUID) (*models.Shop, error) {
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return shop, nil
}

// GetShopPlants gets all plants from a shop
func (s *ShopService) GetShopPlants(ctx context.Context, shopID uuid.UUID) ([]*models.Plant, error) {
	// Check if the shop exists
	_, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}

	// Get the shop's plants
	plants, err := s.shopRepo.GetPlants(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop plants: %w", err)
	}
	return plants, nil
}

// GetSpecialOffers gets all special offers
func (s *ShopService) GetSpecialOffers(ctx context.Context) ([]*models.SpecialOffer, error) {
	offers, err := s.shopRepo.GetSpecialOffers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get special offers: %w", err)
	}
	return offers, nil
}