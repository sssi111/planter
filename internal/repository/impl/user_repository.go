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

// UserRepository is the implementation of the user repository
type UserRepository struct {
	db *db.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, name, email, profile_image_url, language, notifications_enabled, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user locations
	locations, err := r.GetLocations(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user locations: %w", err)
	}
	user.Locations = locations

	// Get favorite plant IDs
	favoritePlantIDs, err := r.GetFavoritePlantIDs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite plant IDs: %w", err)
	}
	user.FavoritePlantIDs = favoritePlantIDs

	// Get owned plant IDs
	ownedPlantIDs, err := r.GetOwnedPlantIDs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned plant IDs: %w", err)
	}
	user.OwnedPlantIDs = ownedPlantIDs

	return &user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, name, email, password_hash, profile_image_url, language, notifications_enabled, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user locations
	locations, err := r.GetLocations(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user locations: %w", err)
	}
	user.Locations = locations

	// Get favorite plant IDs
	favoritePlantIDs, err := r.GetFavoritePlantIDs(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite plant IDs: %w", err)
	}
	user.FavoritePlantIDs = favoritePlantIDs

	// Get owned plant IDs
	ownedPlantIDs, err := r.GetOwnedPlantIDs(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned plant IDs: %w", err)
	}
	user.OwnedPlantIDs = ownedPlantIDs

	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO users (name, email, password_hash, profile_image_url, language, notifications_enabled)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, user.Name, user.Email, user.PasswordHash, user.ProfileImageURL, user.Language, user.NotificationsEnabled).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Insert locations
	for _, location := range user.Locations {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_locations (user_id, location)
			VALUES ($1, $2)
		`, user.ID, location)
		if err != nil {
			return fmt.Errorf("failed to create user location: %w", err)
		}
	}

	return tx.Commit()
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user
	_, err = tx.ExecContext(ctx, `
		UPDATE users
		SET name = $1, profile_image_url = $2, language = $3, notifications_enabled = $4, updated_at = NOW()
		WHERE id = $5
	`, user.Name, user.ProfileImageURL, user.Language, user.NotificationsEnabled, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete existing locations
	_, err = tx.ExecContext(ctx, `
		DELETE FROM user_locations
		WHERE user_id = $1
	`, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user locations: %w", err)
	}

	// Insert new locations
	for _, location := range user.Locations {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_locations (user_id, location)
			VALUES ($1, $2)
		`, user.ID, location)
		if err != nil {
			return fmt.Errorf("failed to create user location: %w", err)
		}
	}

	return tx.Commit()
}

// GetLocations gets a user's locations
func (r *UserRepository) GetLocations(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var locations []string
	err := r.db.SelectContext(ctx, &locations, `
		SELECT location
		FROM user_locations
		WHERE user_id = $1
		ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user locations: %w", err)
	}
	return locations, nil
}

// AddLocation adds a location to a user
func (r *UserRepository) AddLocation(ctx context.Context, userID uuid.UUID, location string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_locations (user_id, location)
		VALUES ($1, $2)
		ON CONFLICT (user_id, location) DO NOTHING
	`, userID, location)
	if err != nil {
		return fmt.Errorf("failed to add user location: %w", err)
	}
	return nil
}

// RemoveLocation removes a location from a user
func (r *UserRepository) RemoveLocation(ctx context.Context, userID uuid.UUID, location string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_locations
		WHERE user_id = $1 AND location = $2
	`, userID, location)
	if err != nil {
		return fmt.Errorf("failed to remove user location: %w", err)
	}
	return nil
}

// GetFavoritePlantIDs gets a user's favorite plant IDs
func (r *UserRepository) GetFavoritePlantIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var plantIDs []string
	err := r.db.SelectContext(ctx, &plantIDs, `
		SELECT plant_id::text
		FROM user_favorite_plants
		WHERE user_id = $1
		ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite plant IDs: %w", err)
	}
	return plantIDs, nil
}

// GetOwnedPlantIDs gets a user's owned plant IDs
func (r *UserRepository) GetOwnedPlantIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var plantIDs []string
	err := r.db.SelectContext(ctx, &plantIDs, `
		SELECT plant_id::text
		FROM user_plants
		WHERE user_id = $1
		ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned plant IDs: %w", err)
	}
	return plantIDs, nil
}