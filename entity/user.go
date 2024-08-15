package entity

import (
	"gorm.io/gorm"
)

type UserRole string

func (role UserRole) String() string {
	return string(role)
}

const (
	AdminR UserRole = "admin"
	UserR  UserRole = "user"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"type:varchar(200);not null"`
	Email    string `gorm:"type:varchar(200);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(150);not null"`
	// Role     UserRole `gorm:"type:varchar(10);default:user"` // for sql server only
	Role UserRole `gorm:"type:user_role;default:user"` // for mysql and postgres
	gorm.Model
}
