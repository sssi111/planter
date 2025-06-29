BEGIN;

-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
DO $$
BEGIN
    -- Only create sunlight_level if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sunlight_level') THEN
        CREATE TYPE sunlight_level AS ENUM ('LOW', 'MEDIUM', 'HIGH');
    END IF;

    -- Only create humidity_level if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'humidity_level') THEN
        CREATE TYPE humidity_level AS ENUM ('LOW', 'MEDIUM', 'HIGH');
    END IF;

    -- Only create language if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'language') THEN
        CREATE TYPE language AS ENUM ('RUSSIAN', 'ENGLISH');
    END IF;
END
$$;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    profile_image_url TEXT,
    language language NOT NULL DEFAULT 'RUSSIAN',
    notifications_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create user_locations table
CREATE TABLE IF NOT EXISTS user_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create care_instructions table
CREATE TABLE IF NOT EXISTS care_instructions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    watering_frequency INTEGER NOT NULL, -- in days
    sunlight sunlight_level NOT NULL,
    min_temperature INTEGER NOT NULL,
    max_temperature INTEGER NOT NULL,
    humidity humidity_level NOT NULL,
    soil_type VARCHAR(255) NOT NULL,
    fertilizer_frequency INTEGER NOT NULL, -- in days
    additional_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create plants table
CREATE TABLE IF NOT EXISTS plants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    scientific_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    care_instructions_id UUID NOT NULL REFERENCES care_instructions(id) ON DELETE CASCADE,
    price DECIMAL(10, 2),
    shop_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create user_plants table (for owned plants)
CREATE TABLE IF NOT EXISTS user_plants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    location VARCHAR(255),
    last_watered TIMESTAMP WITH TIME ZONE,
    next_watering TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, plant_id)
);

-- Create user_favorite_plants table
CREATE TABLE IF NOT EXISTS user_favorite_plants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, plant_id)
);

-- Create shops table
CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    rating DECIMAL(2, 1) NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create shop_plants table
CREATE TABLE IF NOT EXISTS shop_plants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(shop_id, plant_id)
);

-- Create special_offers table
CREATE TABLE IF NOT EXISTS special_offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    discount_percentage INTEGER NOT NULL,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create plant_questionnaire table
CREATE TABLE IF NOT EXISTS plant_questionnaires (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    sunlight_preference sunlight_level NOT NULL,
    pet_friendly BOOLEAN NOT NULL,
    care_level INTEGER NOT NULL, -- 1-5 scale
    preferred_location VARCHAR(255),
    additional_preferences TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create plant_recommendations table
CREATE TABLE IF NOT EXISTS plant_recommendations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    questionnaire_id UUID NOT NULL REFERENCES plant_questionnaires(id) ON DELETE CASCADE,
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    score DECIMAL(3, 2) NOT NULL, -- 0-1 scale
    reasoning TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(questionnaire_id, plant_id)
);

-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index for faster notification queries
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_plants_name ON plants(name);
CREATE INDEX IF NOT EXISTS idx_plants_scientific_name ON plants(scientific_name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_user_plants_user_id ON user_plants(user_id);
CREATE INDEX IF NOT EXISTS idx_user_favorite_plants_user_id ON user_favorite_plants(user_id);
CREATE INDEX IF NOT EXISTS idx_shop_plants_shop_id ON shop_plants(shop_id);

COMMIT;