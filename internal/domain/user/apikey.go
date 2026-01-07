package user

import "time"

// APIKey represents API key for authentication
type APIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	Name       string     `gorm:"size:100" json:"name"`
	KeyHash    string     `gorm:"size:255;uniqueIndex;not null" json:"-"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	LastUsedAt *time.Time `gorm:"type:timestamp" json:"last_used_at"`
	CreatedAt  time.Time  `gorm:"type:timestamp;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"type:timestamp;autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
