package postgres

import (
	"time"

	"starter-gofiber/internal/domain/user"

	"gorm.io/gorm"
)

type EmailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Create(verification *user.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *EmailVerificationRepository) FindByToken(token string) (*user.EmailVerification, error) {
	var verification user.EmailVerification
	err := r.db.Where("token = ? AND is_verified = ? AND expires_at > ?", token, false, time.Now()).
		Preload("User").
		First(&verification).Error
	return &verification, err
}

func (r *EmailVerificationRepository) MarkAsVerified(token string) error {
	return r.db.Model(&user.EmailVerification{}).
		Where("token = ?", token).
		Update("is_verified", true).Error
}

func (r *EmailVerificationRepository) DeleteUserVerifications(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&user.EmailVerification{}).Error
}

func (r *EmailVerificationRepository) DeleteExpiredVerifications() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&user.EmailVerification{}).Error
}
