package repository

import (
	"starter-gofiber/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(d *gorm.DB) *UserRepository {
	return &UserRepository{
		db: d,
	}
}

func (u *UserRepository) Create(m entity.User) error {
	return u.db.Create(&m).Error
}

func (u *UserRepository) ExistEmail(email string) error {
	var user entity.User
	err := u.db.Where("email = ?", email).First(&user).Error
	return err
}

func (u *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user *entity.User
	err := u.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (u *UserRepository) FindByID(id uint) (*entity.User, error) {
	var user *entity.User
	err := u.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (u *UserRepository) Update(user *entity.User) error {
	return u.db.Save(user).Error
}
