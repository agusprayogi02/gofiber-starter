package postgres

import (
	"starter-gofiber/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(d *gorm.DB) user.Repository {
	return &UserRepository{
		db: d,
	}
}

func (u *UserRepository) Create(m *user.User) error {
	return u.db.Create(m).Error
}

func (u *UserRepository) ExistEmail(email string) error {
	var usr user.User
	err := u.db.Where("email = ?", email).First(&usr).Error
	return err
}

func (u *UserRepository) FindByEmail(email string) (*user.User, error) {
	var usr user.User
	err := u.db.Where("email = ?", email).First(&usr).Error
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (u *UserRepository) FindByID(id uint) (*user.User, error) {
	var usr user.User
	err := u.db.Where("id = ?", id).First(&usr).Error
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func (u *UserRepository) Update(usr *user.User) error {
	return u.db.Save(usr).Error
}

// RefreshToken operations
func (u *UserRepository) CreateRefreshToken(token *user.RefreshToken) error {
	return u.db.Create(token).Error
}

func (u *UserRepository) FindRefreshTokenByToken(tokenStr string) (*user.RefreshToken, error) {
	var token user.RefreshToken
	err := u.db.Where("token = ? AND is_revoked = ?", tokenStr, false).
		Preload("User").
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (u *UserRepository) RevokeRefreshToken(tokenStr string) error {
	return u.db.Model(&user.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("is_revoked", true).Error
}

func (u *UserRepository) RevokeAllUserTokens(userID uint) error {
	return u.db.Model(&user.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

func (u *UserRepository) FindUserRefreshTokens(userID uint) ([]user.RefreshToken, error) {
	var tokens []user.RefreshToken
	err := u.db.Where("user_id = ? AND is_revoked = ?", userID, false).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

func (u *UserRepository) RevokeSessionByID(sessionID uint, userID uint) error {
	return u.db.Model(&user.RefreshToken{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Update("is_revoked", true).Error
}

// EmailVerification operations
func (u *UserRepository) CreateEmailVerification(verification *user.EmailVerification) error {
	return u.db.Create(verification).Error
}

func (u *UserRepository) FindEmailVerificationByToken(tokenStr string) (*user.EmailVerification, error) {
	var verification user.EmailVerification
	err := u.db.Where("token = ? AND is_verified = ?", tokenStr, false).
		Preload("User").
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (u *UserRepository) MarkEmailAsVerified(userID uint) error {
	return u.db.Model(&user.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}

// PasswordReset operations
func (u *UserRepository) CreatePasswordReset(reset *user.PasswordReset) error {
	return u.db.Create(reset).Error
}

func (u *UserRepository) FindPasswordResetByToken(tokenStr string) (*user.PasswordReset, error) {
	var reset user.PasswordReset
	err := u.db.Where("token = ? AND is_used = ?", tokenStr, false).
		Preload("User").
		First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (u *UserRepository) MarkPasswordResetAsUsed(tokenStr string) error {
	return u.db.Model(&user.PasswordReset{}).
		Where("token = ?", tokenStr).
		Update("is_used", true).Error
}

// APIKey operations
func (u *UserRepository) CreateAPIKey(apiKey *user.APIKey) error {
	return u.db.Create(apiKey).Error
}

func (u *UserRepository) FindAPIKeyByHash(hash string) (*user.APIKey, error) {
	var key user.APIKey
	err := u.db.Where("key_hash = ? AND is_active = ?", hash, true).
		Preload("User").
		First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (u *UserRepository) UpdateAPIKey(apiKey *user.APIKey) error {
	return u.db.Save(apiKey).Error
}
