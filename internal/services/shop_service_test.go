package services

import (
	"context"
	"testing"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockShopRepository is a mock implementation of the ShopRepository interface
type MockShopRepository struct {
	mock.Mock
}

func (m *MockShopRepository) GetAll(ctx context.Context) ([]*models.Shop, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Shop), args.Error(1)
}

func (m *MockShopRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shop), args.Error(1)
}

func (m *MockShopRepository) GetPlants(ctx context.Context, shopID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, shopID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockShopRepository) GetSpecialOffers(ctx context.Context) ([]*models.SpecialOffer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.SpecialOffer), args.Error(1)
}

// TestShopService_GetAllShops tests the GetAllShops method of the ShopService
func TestShopService_GetAllShops(t *testing.T) {
	// Create a mock shop repository
	mockShopRepo := new(MockShopRepository)

	// Create test shops
	shop1 := &models.Shop{
		ID:      uuid.New(),
		Name:    "Shop 1",
		Address: "Address 1",
		Rating:  4.5,
	}
	shop2 := &models.Shop{
		ID:      uuid.New(),
		Name:    "Shop 2",
		Address: "Address 2",
		Rating:  4.8,
	}
	shops := []*models.Shop{shop1, shop2}

	// Set up the mock expectations
	mockShopRepo.On("GetAll", mock.Anything).Return(shops, nil)

	// Create the shop service
	shopService := NewShopService(mockShopRepo)

	// Test the GetAllShops method
	result, err := shopService.GetAllShops(context.Background())

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected shops
	assert.Equal(t, shops, result)

	// Verify that all expectations were met
	mockShopRepo.AssertExpectations(t)
}

// TestShopService_GetShop tests the GetShop method of the ShopService
func TestShopService_GetShop(t *testing.T) {
	// Create a mock shop repository
	mockShopRepo := new(MockShopRepository)

	// Create a test shop
	shopID := uuid.New()
	shop := &models.Shop{
		ID:      shopID,
		Name:    "Test Shop",
		Address: "Test Address",
		Rating:  4.7,
	}

	// Set up the mock expectations
	mockShopRepo.On("GetByID", mock.Anything, shopID).Return(shop, nil)

	// Create the shop service
	shopService := NewShopService(mockShopRepo)

	// Test the GetShop method
	result, err := shopService.GetShop(context.Background(), shopID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected shop
	assert.Equal(t, shop, result)

	// Verify that all expectations were met
	mockShopRepo.AssertExpectations(t)
}

// TestShopService_GetShopPlants tests the GetShopPlants method of the ShopService
func TestShopService_GetShopPlants(t *testing.T) {
	// Create a mock shop repository
	mockShopRepo := new(MockShopRepository)

	// Create a test shop and plants
	shopID := uuid.New()
	shop := &models.Shop{
		ID:      shopID,
		Name:    "Test Shop",
		Address: "Test Address",
		Rating:  4.7,
	}
	plant1 := &models.Plant{
		ID:          uuid.New(),
		Name:        "Plant 1",
		Description: "Description 1",
	}
	plant2 := &models.Plant{
		ID:          uuid.New(),
		Name:        "Plant 2",
		Description: "Description 2",
	}
	plants := []*models.Plant{plant1, plant2}

	// Set up the mock expectations
	mockShopRepo.On("GetByID", mock.Anything, shopID).Return(shop, nil)
	mockShopRepo.On("GetPlants", mock.Anything, shopID).Return(plants, nil)

	// Create the shop service
	shopService := NewShopService(mockShopRepo)

	// Test the GetShopPlants method
	result, err := shopService.GetShopPlants(context.Background(), shopID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected plants
	assert.Equal(t, plants, result)

	// Verify that all expectations were met
	mockShopRepo.AssertExpectations(t)
}

// TestShopService_GetSpecialOffers tests the GetSpecialOffers method of the ShopService
func TestShopService_GetSpecialOffers(t *testing.T) {
	// Create a mock shop repository
	mockShopRepo := new(MockShopRepository)

	// Create test special offers
	offer1 := &models.SpecialOffer{
		ID:                uuid.New(),
		Title:             "Offer 1",
		Description:       "Description 1",
		DiscountPercentage: 10,
	}
	offer2 := &models.SpecialOffer{
		ID:                uuid.New(),
		Title:             "Offer 2",
		Description:       "Description 2",
		DiscountPercentage: 20,
	}
	offers := []*models.SpecialOffer{offer1, offer2}

	// Set up the mock expectations
	mockShopRepo.On("GetSpecialOffers", mock.Anything).Return(offers, nil)

	// Create the shop service
	shopService := NewShopService(mockShopRepo)

	// Test the GetSpecialOffers method
	result, err := shopService.GetSpecialOffers(context.Background())

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected offers
	assert.Equal(t, offers, result)

	// Verify that all expectations were met
	mockShopRepo.AssertExpectations(t)
}