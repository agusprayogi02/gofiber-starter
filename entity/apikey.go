package entity

// APIKey represents API key for authentication
type APIKey struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	UserID     uint    `gorm:"not null;index" json:"user_id"`
	Name       string  `gorm:"size:100" json:"name"`
	KeyHash    string  `gorm:"size:255;uniqueIndex;not null" json:"-"`
	IsActive   bool    `gorm:"default:true" json:"is_active"`
	LastUsedAt *string `gorm:"type:timestamp" json:"last_used_at"`
	CreatedAt  string  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  string  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
