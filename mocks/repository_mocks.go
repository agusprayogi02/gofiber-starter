package mocks

import (
	postdomain "starter-gofiber/internal/domain/post"
	"starter-gofiber/internal/domain/user"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *user.User) error {
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

func (m *MockPostRepository) FindAll(limit, offset int) ([]postdomain.Post, int64, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]postdomain.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) Create(post *postdomain.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) FindByID(id uint) (*postdomain.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postdomain.Post), args.Error(1)
}

func (m *MockPostRepository) Update(post *postdomain.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostRepository) FindByUserID(userID uint, limit, offset int) ([]postdomain.Post, int64, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]postdomain.Post), args.Get(1).(int64), args.Error(2)
}

// MockRefreshTokenRepository is a mock implementation of RefreshTokenRepository
// Note: RefreshToken operations are now part of user.Repository interface
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) CreateRefreshToken(token *user.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindRefreshTokenByToken(token string) (*user.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeRefreshToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindUserRefreshTokens(userID uint) ([]user.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeSessionByID(sessionID uint, userID uint) error {
	args := m.Called(sessionID, userID)
	return args.Error(0)
}

// MockPasswordResetRepository is a mock implementation of PasswordResetRepository
// Note: PasswordReset operations are now part of user.Repository interface
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) CreatePasswordReset(reset *user.PasswordReset) error {
	args := m.Called(reset)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) FindPasswordResetByToken(token string) (*user.PasswordReset, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.PasswordReset), args.Error(1)
}

func (m *MockPasswordResetRepository) MarkPasswordResetAsUsed(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// MockEmailVerificationRepository is a mock implementation of EmailVerificationRepository
// Note: EmailVerification operations are now part of user.Repository interface
type MockEmailVerificationRepository struct {
	mock.Mock
}

func (m *MockEmailVerificationRepository) CreateEmailVerification(verification *user.EmailVerification) error {
	args := m.Called(verification)
	return args.Error(0)
}

func (m *MockEmailVerificationRepository) FindEmailVerificationByToken(token string) (*user.EmailVerification, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.EmailVerification), args.Error(1)
}

func (m *MockEmailVerificationRepository) MarkEmailAsVerified(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Helper function for GORM errors
func GormNotFoundError() error {
	return gorm.ErrRecordNotFound
}
