-- Get user's plants and update next_watering date to yesterday
UPDATE user_plants
SET next_watering = NOW() - INTERVAL '1 day'
WHERE user_id = '1f1f1bbc-f058-4b28-97ea-7c8c83fa29bb' -- Replace with your user ID
  AND plant_id = '388393f6-6b0b-425b-9d4e-9df5b055eae8'; -- Replace with your plant ID

-- Verify the update
SELECT up.next_watering, p.name, u.email
FROM user_plants up
JOIN plants p ON p.id = up.plant_id
JOIN users u ON u.id = up.user_id
WHERE up.user_id = '1f1f1bbc-f058-4b28-97ea-7c8c83fa29bb'
  AND up.plant_id = '388393f6-6b0b-425b-9d4e-9df5b055eae8'; 