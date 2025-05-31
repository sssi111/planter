-- Script to clear all plants and related data from the database
-- WARNING: This will permanently delete all plant data

DO $$
BEGIN
    -- First delete dependent records if they exist
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'notifications') THEN
        DELETE FROM notifications WHERE plant_id IS NOT NULL;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'plant_recommendations') THEN
        DELETE FROM plant_recommendations WHERE plant_id IS NOT NULL;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'shop_plants') THEN
        DELETE FROM shop_plants WHERE plant_id IS NOT NULL;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_favorite_plants') THEN
        DELETE FROM user_favorite_plants WHERE plant_id IS NOT NULL;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_plants') THEN
        DELETE FROM user_plants WHERE plant_id IS NOT NULL;
    END IF;

    -- Then delete the plants themselves
    DELETE FROM plants;
    
    RAISE NOTICE 'Successfully cleared all plant data';
EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Error clearing plant data: %', SQLERRM;
END $$;

-- Optional: Reset any sequences if needed
-- SELECT setval('plants_id_seq', 1, false);