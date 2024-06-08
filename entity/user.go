package entity

import (
	"gorm.io/gorm"
	"starter-gofiber/config"
	"starter-gofiber/dto"
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
	ID       uint     `gorm:"type:primaryKey;autoIncrement"`
	Name     string   `gorm:"type:varchar(200);not null"`
	Email    string   `gorm:"type:varchar(200);uniqueIndex;not null"`
	Password string   `gorm:"type:varchar(150);not null"`
	Role     UserRole `gorm:"type:user_role;default:user"`
	gorm.Model
}

func (u *User) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt.Format(config.FORMAT_TIME),
		UpdatedAt: u.UpdatedAt.Format(config.FORMAT_TIME),
	}
}
