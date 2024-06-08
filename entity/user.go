package entity

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type UserRole string

func (role *UserRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*role = UserRole(s)
	switch *role {
	case AdminR, UserR:
		return nil
	default:
		return fmt.Errorf("invalid role: %s", s)
	}
}

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
