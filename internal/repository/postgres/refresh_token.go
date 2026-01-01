package postgres

import (
	"time"

	"starter-gofiber/internal/domain/user"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *user.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepository) FindByToken(token string) (*user.RefreshToken, error) {
	var refreshToken user.RefreshToken
	err := r.db.Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		Preload("User").
		First(&refreshToken).Error
	return &refreshToken, err
}

func (r *RefreshTokenRepository) RevokeToken(token string) error {
	return r.db.Model(&user.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Model(&user.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&user.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) GetUserActiveSessions(userID uint) ([]user.RefreshToken, error) {
	var tokens []user.RefreshToken
	err := r.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

func (r *RefreshTokenRepository) RevokeSessionByID(id uint, userID uint) error {
	return r.db.Model(&user.RefreshToken{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_revoked", true).Error
}
