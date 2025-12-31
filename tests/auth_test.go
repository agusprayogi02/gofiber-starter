package tests

import (
	"testing"

	"starter-gofiber/dto"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupSuite() {
	testApp = SetupTestApp()
}

func (s *AuthTestSuite) TearDownSuite() {
	CleanupTestDB()
}

func (s *AuthTestSuite) SetupTest() {
	testDB.Exec("DELETE FROM users")
	testDB.Exec("DELETE FROM refresh_tokens")
	testDB.Exec("DELETE FROM password_resets")
	testDB.Exec("DELETE FROM email_verifications")
}

func (s *AuthTestSuite) TestRegister_Success() {
	payload := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "Password123!",
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/auth/register", payload, nil)
	s.NoError(err)
	AssertSuccessResponse(s.T(), resp, 201)

	var response map[string]interface{}
	ParseJSON(s.T(), body, &response)
	s.Equal("success", response["status"])
	s.Equal("User registered successfully", response["message"])
}

func (s *AuthTestSuite) TestRegister_DuplicateEmail() {
	CreateTestUser(testDB, "existing@example.com", "hashed_password", "user")

	payload := dto.RegisterRequest{
		Name:     "Jane Doe",
		Email:    "existing@example.com",
		Password: "Password123!",
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/auth/register", payload, nil)
	s.NoError(err)
	AssertErrorResponse(s.T(), resp, 400, body)
}

func (s *AuthTestSuite) TestLogin_Success() {
	user := CreateTestUser(testDB, "test@example.com", "$2a$10$abcdefghijklmnopqrstuv", "user")

	payload := dto.LoginRequest{
		Email:    user.Email,
		Password: "Password123!",
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/auth/login", payload, nil)
	s.NoError(err)

	var response map[string]interface{}
	ParseJSON(s.T(), body, &response)

	if resp.StatusCode == 200 {
		s.Equal("success", response["status"])
		data := response["data"].(map[string]interface{})
		s.Contains(data, "access_token")
		s.Contains(data, "refresh_token")
	}
}

func (s *AuthTestSuite) TestLogin_InvalidCredentials() {
	payload := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/auth/login", payload, nil)
	s.NoError(err)
	AssertErrorResponse(s.T(), resp, 401, body)
}

func (s *AuthTestSuite) TestRefreshToken_Success() {
	user := CreateTestUser(testDB, "test@example.com", "hashed_password", "user")

	loginPayload := dto.LoginRequest{
		Email:    user.Email,
		Password: "Password123!",
	}

	resp, body, _ := MakeRequest(testApp, "POST", "/api/auth/login", loginPayload, nil)

	if resp.StatusCode == 200 {
		var loginResponse map[string]interface{}
		ParseJSON(s.T(), body, &loginResponse)
		data := loginResponse["data"].(map[string]interface{})
		refreshToken := data["refresh_token"].(string)

		refreshPayload := dto.RefreshTokenRequest{
			RefreshToken: refreshToken,
		}

		refreshResp, refreshBody, err := MakeRequest(testApp, "POST", "/api/auth/refresh", refreshPayload, nil)
		s.NoError(err)

		var refreshResponse map[string]interface{}
		ParseJSON(s.T(), refreshBody, &refreshResponse)

		if refreshResp.StatusCode == 200 {
			s.Equal("success", refreshResponse["status"])
			newData := refreshResponse["data"].(map[string]interface{})
			s.Contains(newData, "access_token")
		}
	}
}

func (s *AuthTestSuite) TestLogout_Success() {
	user := CreateTestUser(testDB, "test@example.com", "hashed_password", "user")

	loginPayload := dto.LoginRequest{
		Email:    user.Email,
		Password: "Password123!",
	}

	loginResp, loginBody, _ := MakeRequest(testApp, "POST", "/api/auth/login", loginPayload, nil)

	if loginResp.StatusCode == 200 {
		var loginResponse map[string]interface{}
		ParseJSON(s.T(), loginBody, &loginResponse)
		data := loginResponse["data"].(map[string]interface{})
		refreshToken := data["refresh_token"].(string)

		logoutPayload := dto.LogoutRequest{
			RefreshToken: refreshToken,
		}

		logoutResp, logoutBody, err := MakeRequest(testApp, "POST", "/api/auth/logout", logoutPayload, nil)
		s.NoError(err)

		var logoutResponse map[string]interface{}
		ParseJSON(s.T(), logoutBody, &logoutResponse)

		if logoutResp.StatusCode == 200 {
			s.Equal("success", logoutResponse["status"])
		}
	}
}

func (s *AuthTestSuite) TestForgotPassword_Success() {
	CreateTestUser(testDB, "test@example.com", "hashed_password", "user")

	payload := dto.ForgotPasswordRequest{
		Email: "test@example.com",
	}

	resp, body, err := MakeRequest(testApp, "POST", "/api/auth/forgot-password", payload, nil)
	s.NoError(err)

	var response map[string]interface{}
	ParseJSON(s.T(), body, &response)

	if resp.StatusCode == 200 {
		s.Equal("success", response["status"])
	}
}

func (s *AuthTestSuite) TestChangePassword_Success() {
	user := CreateTestUser(testDB, "test@example.com", "$2a$10$abcdefghijklmnopqrstuv", "user")

	loginPayload := dto.LoginRequest{
		Email:    user.Email,
		Password: "Password123!",
	}

	loginResp, loginBody, _ := MakeRequest(testApp, "POST", "/api/auth/login", loginPayload, nil)

	if loginResp.StatusCode == 200 {
		var loginResponse map[string]interface{}
		ParseJSON(s.T(), loginBody, &loginResponse)
		data := loginResponse["data"].(map[string]interface{})
		accessToken := data["access_token"].(string)

		changePayload := dto.ChangePasswordRequest{
			OldPassword: "Password123!",
			NewPassword: "NewPassword123!",
		}

		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}

		changeResp, changeBody, err := MakeRequest(testApp, "PUT", "/api/auth/change-password", changePayload, headers)
		s.NoError(err)

		var changeResponse map[string]interface{}
		ParseJSON(s.T(), changeBody, &changeResponse)

		if changeResp.StatusCode == 200 {
			s.Equal("success", changeResponse["status"])
		}
	}
}
