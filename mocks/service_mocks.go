package mocks

import (
	"starter-gofiber/dto"
	"starter-gofiber/entity"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req dto.RegisterRequest, enforcer interface{}) (*entity.User, error) {
	args := m.Called(req, enforcer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthService) Login(req dto.LoginRequest, ipAddress, userAgent string) (*dto.LoginResponse, error) {
	args := m.Called(req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RefreshTokenResponse), args.Error(1)
}

func (m *MockAuthService) Logout(userID uint, refreshToken string) error {
	args := m.Called(userID, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(req dto.ForgotPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(req dto.ResetPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *MockAuthService) VerifyEmail(req dto.VerifyEmailRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ResendVerificationEmail(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) GetActiveSessions(userID uint) ([]dto.SessionResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SessionResponse), args.Error(1)
}

func (m *MockAuthService) RevokeSession(sessionID, userID uint) error {
	args := m.Called(sessionID, userID)
	return args.Error(0)
}

// MockPostService is a mock implementation of PostService
type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) All(page, perPage int) ([]entity.Post, dto.Pagination, error) {
	args := m.Called(page, perPage)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Pagination), args.Error(2)
	}
	return args.Get(0).([]entity.Post), args.Get(1).(dto.Pagination), args.Error(2)
}

func (m *MockPostService) Create(post *dto.PostRequest) (*entity.Post, error) {
	args := m.Called(post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostService) Find(id int) (*entity.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostService) Update(id int, post *dto.PostUpdateRequest) (*entity.Post, error) {
	args := m.Called(id, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
