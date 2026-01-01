package postgres

import (
	"starter-gofiber/internal/domain/post"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(d *gorm.DB) post.Repository {
	return &PostRepository{
		db: d,
	}
}

func (r *PostRepository) Create(p *post.Post) error {
	err := r.db.Create(p).Error
	if err != nil {
		return err
	}
	// Preload User relation
	err = r.db.Model(p).Preload("User").First(p).Error
	return err
}

func (r *PostRepository) FindByID(id uint) (*post.Post, error) {
	var p post.Post
	err := r.db.Preload("User").Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepository) FindAll(limit, offset int) ([]post.Post, int64, error) {
	var posts []post.Post
	var total int64

	// Count total
	if err := r.db.Model(&post.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Preload("User").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) Update(p *post.Post) error {
	return r.db.Model(p).Updates(p).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&post.Post{}, id).Error
}

func (r *PostRepository) FindByUserID(userID uint, limit, offset int) ([]post.Post, int64, error) {
	var posts []post.Post
	var total int64

	// Count total for this user
	if err := r.db.Model(&post.Post{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}
