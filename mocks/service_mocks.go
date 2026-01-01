package mocks

import (
	postdomain "starter-gofiber/internal/domain/post"
	"starter-gofiber/internal/domain/user"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *user.RegisterRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) Login(req *user.LoginRequest, ipAddress, userAgent string) (*user.LoginResponse, error) {
	args := m.Called(req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(req *user.RefreshTokenRequest, ipAddress, userAgent string) (*user.RefreshTokenResponse, error) {
	args := m.Called(req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.RefreshTokenResponse), args.Error(1)
}

func (m *MockAuthService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(req *user.ForgotPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(req *user.ResetPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ChangePassword(userID uint, req *user.ChangePasswordRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *MockAuthService) VerifyEmail(req *user.VerifyEmailRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ResendVerificationEmail(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthService) GetActiveSessions(userID uint) ([]user.SessionResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.SessionResponse), args.Error(1)
}

func (m *MockAuthService) RevokeSession(sessionID, userID uint) error {
	args := m.Called(sessionID, userID)
	return args.Error(0)
}

// MockPostService is a mock implementation of PostService
type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) FindAll(page, limit int) ([]postdomain.PostResponse, *postdomain.PaginationMeta, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*postdomain.PaginationMeta), args.Error(2)
	}
	return args.Get(0).([]postdomain.PostResponse), args.Get(1).(*postdomain.PaginationMeta), args.Error(2)
}

func (m *MockPostService) Create(req *postdomain.PostRequest, userID uint) (*postdomain.PostResponse, error) {
	args := m.Called(req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postdomain.PostResponse), args.Error(1)
}

func (m *MockPostService) FindByID(id uint) (*postdomain.PostResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postdomain.PostResponse), args.Error(1)
}

func (m *MockPostService) Update(id uint, req *postdomain.PostUpdateRequest, userID uint) (*postdomain.PostResponse, error) {
	args := m.Called(id, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postdomain.PostResponse), args.Error(1)
}

func (m *MockPostService) Delete(id uint, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockPostService) FindByUserID(userID uint, page, limit int) ([]postdomain.PostResponse, *postdomain.PaginationMeta, error) {
	args := m.Called(userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*postdomain.PaginationMeta), args.Error(2)
	}
	return args.Get(0).([]postdomain.PostResponse), args.Get(1).(*postdomain.PaginationMeta), args.Error(2)
}
