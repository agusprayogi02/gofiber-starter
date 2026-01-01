package user

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(500);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsRevoked bool      `gorm:"default:false"`
	DeviceID  string    `gorm:"type:varchar(255)"` // Optional: track device
	IPAddress string    `gorm:"type:varchar(45)"`  // Optional: track IP
	UserAgent string    `gorm:"type:text"`         // Optional: track user agent
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	gorm.Model
}
