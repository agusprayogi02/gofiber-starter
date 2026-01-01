package postgres

import (
	"time"

	"starter-gofiber/internal/domain/user"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(reset *user.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *PasswordResetRepository) FindByToken(token string) (*user.PasswordReset, error) {
	var reset user.PasswordReset
	err := r.db.Where("token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).
		Preload("User").
		First(&reset).Error
	return &reset, err
}

func (r *PasswordResetRepository) MarkAsUsed(token string) error {
	return r.db.Model(&user.PasswordReset{}).
		Where("token = ?", token).
		Update("is_used", true).Error
}

func (r *PasswordResetRepository) DeleteUserResets(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&user.PasswordReset{}).Error
}

func (r *PasswordResetRepository) DeleteExpiredResets() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&user.PasswordReset{}).Error
}
