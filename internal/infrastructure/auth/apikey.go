package auth

import (
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/pkg/logger"
	"starter-gofiber/pkg/utils"

	"gorm.io/gorm"
)

// ValidateAPIKey checks if API key is valid
// Returns (isValid, userID)
func ValidateAPIKey(db *gorm.DB, apiKey string) (bool, uint) {
	type APIKey struct {
		ID         uint
		UserID     uint
		IsActive   bool
		KeyHash    string
		LastUsedAt *string
	}

	var key APIKey
	result := db.Table("api_keys").
		Where("key_hash = ? AND is_active = ?", crypto.HashString(apiKey), true).
		First(&key)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, 0
		}
		logger.Error("API key validation error: " + result.Error.Error())
		return false, 0
	}

	// Update last used timestamp async
	go func() {
		db.Table("api_keys").
			Where("id = ?", key.ID).
			Update("last_used_at", utils.TimeNow())
	}()

	return true, key.UserID
}

// GenerateAPIKey creates a new API key for user
func GenerateAPIKey(db *gorm.DB, userID uint, name string) (string, error) {
	// Generate random API key
	apiKey := crypto.GenerateRandomString(32)

	// Hash the key
	keyHash := crypto.HashString(apiKey)

	// Save to database
	type APIKey struct {
		UserID   uint
		Name     string
		KeyHash  string
		IsActive bool
	}

	key := APIKey{
		UserID:   userID,
		Name:     name,
		KeyHash:  keyHash,
		IsActive: true,
	}

	result := db.Table("api_keys").Create(&key)
	if result.Error != nil {
		return "", result.Error
	}

	// Return the plain API key (only time it's visible)
	return apiKey, nil
}

// RevokeAPIKey deactivates an API key
func RevokeAPIKey(db *gorm.DB, keyID uint) error {
	result := db.Table("api_keys").
		Where("id = ?", keyID).
		Update("is_active", false)

	return result.Error
}

// ListAPIKeys returns all API keys for a user
func ListAPIKeys(db *gorm.DB, userID uint) ([]map[string]interface{}, error) {
	var keys []map[string]interface{}
	result := db.Table("api_keys").
		Select("id", "user_id", "name", "is_active", "last_used_at", "created_at").
		Where("user_id = ?", userID).
		Find(&keys)

	return keys, result.Error
}
