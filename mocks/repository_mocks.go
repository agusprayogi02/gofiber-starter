package mocks

import (
	"starter-gofiber/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) ExistEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// MockPostRepository is a mock implementation of PostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) All(page, perPage int) ([]entity.Post, int64, error) {
	args := m.Called(page, perPage)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]entity.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) Create(post *entity.Post) (*entity.Post, error) {
	args := m.Called(post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepository) Find(id int) (*entity.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepository) Update(id int, post *entity.Post) (*entity.Post, error) {
	args := m.Called(id, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockRefreshTokenRepository is a mock implementation of RefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(refreshToken *entity.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*entity.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetUserActiveSessions(userID uint) ([]entity.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeSessionByID(id uint, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockPasswordResetRepository is a mock implementation of PasswordResetRepository
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Create(reset *entity.PasswordReset) error {
	args := m.Called(reset)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) FindByToken(token string) (*entity.PasswordReset, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PasswordReset), args.Error(1)
}

func (m *MockPasswordResetRepository) MarkAsUsed(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeleteUserResets(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeleteExpiredResets() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPasswordResetRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockEmailVerificationRepository is a mock implementation of EmailVerificationRepository
type MockEmailVerificationRepository struct {
	mock.Mock
}

func (m *MockEmailVerificationRepository) Create(verification *entity.EmailVerification) error {
	args := m.Called(verification)
	return args.Error(0)
}

func (m *MockEmailVerificationRepository) FindByToken(token string) (*entity.EmailVerification, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmailVerification), args.Error(1)
}

func (m *MockEmailVerificationRepository) MarkAsVerified(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmailVerificationRepository) DeleteUserVerifications(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockEmailVerificationRepository) DeleteExpiredVerifications() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockEmailVerificationRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Helper function for GORM errors
func GormNotFoundError() error {
	return gorm.ErrRecordNotFound
}
