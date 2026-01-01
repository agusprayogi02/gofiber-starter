package user

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(500);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsUsed    bool      `gorm:"default:false"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	gorm.Model
}
