-- Script to clear all user plant associations from the database
-- WARNING: This will permanently delete all user plant data

DO $$
BEGIN
    -- First check and delete from notifications if it exists
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'notifications') THEN
        DELETE FROM notifications WHERE plant_id IN (SELECT plant_id FROM user_plants);
        DELETE FROM notifications WHERE plant_id IN (SELECT plant_id FROM user_favorite_plants);
    END IF;
    
    -- Delete favorites
    DELETE FROM user_favorite_plants;
    
    -- Delete user plants
    DELETE FROM user_plants;
    
    RAISE NOTICE 'Successfully cleared all user plant data';
EXCEPTION WHEN OTHERS THEN
    RAISE EXCEPTION 'Error clearing user plant data: %', SQLERRM;
END $$;

-- Optional: Reset sequences if needed
-- SELECT setval('user_plants_id_seq', 1, false);
-- SELECT setval('user_favorite_plants_id_seq', 1, false);