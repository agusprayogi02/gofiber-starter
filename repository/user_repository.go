package repository

import (
	"gorm.io/gorm"
	"starter-gofiber/entity"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUser(d *gorm.DB) *UserRepository {
	return &UserRepository{
		db: d,
	}
}

func (u *UserRepository) Register(m entity.User) error {
	return u.db.Create(m).Error
}

func (u *UserRepository) ExistEmail(email string) bool {
	var user entity.User
	u.db.Where("email = ?", email).First(&user)
	return user.ID != 0
}
