package repository

import (
	"starter-gofiber/dto"
	"starter-gofiber/entity"
	"starter-gofiber/helper"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(d *gorm.DB) *PostRepository {
	return &PostRepository{
		db: d,
	}
}

func (u *PostRepository) Create(m entity.Post) (entity.Post, error) {
	err := u.db.Create(&m).Error
	return m, err
}

func (u *PostRepository) FindId(id uint) (post *entity.Post, err error) {
	err = u.db.Where("id = ?", id).First(&post).Error
	return post, err
}

func (u *PostRepository) All(paginate *dto.Pagination) (posts *[]entity.Post, err error) {
	if err := u.db.Scopes(helper.Paginate(posts, paginate, u.db)).Joins("User").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (u *PostRepository) Update(m entity.Post, id uint) error {
	return u.db.Where(&entity.Post{ID: id}).Updates(&m).Error
}

func (u *PostRepository) Delete(id uint) error {
	return u.db.Delete(&entity.Post{}, id).Error
}

func (u *PostRepository) FindByUserId(userId uint) (posts []entity.Post, err error) {
	err = u.db.Where(&entity.Post{UserID: userId}).Find(&posts).Error
	return posts, err
}
