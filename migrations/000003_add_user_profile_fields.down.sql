-- Rollback: Remove user profile fields
-- Created: 2026-01-01

-- Drop user_preferences table
DROP TABLE IF EXISTS user_preferences;

-- Remove columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS avatar;
ALTER TABLE users DROP COLUMN IF EXISTS bio;

