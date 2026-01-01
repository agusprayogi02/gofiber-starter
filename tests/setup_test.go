package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"starter-gofiber/internal/config"
	"starter-gofiber/internal/domain/post"
	"starter-gofiber/internal/domain/user"
	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/router"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testDB  *gorm.DB
	testApp *fiber.App
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	err = db.AutoMigrate(
		&user.User{},
		&post.Post{},
		&user.RefreshToken{},
		&user.PasswordReset{},
		&user.EmailVerification{},
	)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	return db
}

func SetupTestApp() *fiber.App {
	testDB = SetupTestDB()
	config.DB = testDB

	// Setup minimal ENV config for testing
	if config.ENV == nil {
		config.ENV = &config.Config{
			ENV_TYPE:      "test",
			LOCATION_CERT: "../assets/certs/certificate.pem",
		}
	}
	// Set environment variable for middleware detection
	os.Setenv("ENV_TYPE", "test")

	// Initialize private key for JWT
	if err := crypto.InitPrivateKey(config.ENV.LOCATION_CERT); err != nil {
		panic("failed to initialize private key: " + err.Error())
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: apierror.ErrorHelper,
	})

	enforcer, err := casbin.NewEnforcer("../assets/rbac/model.conf", "../assets/rbac/policy.csv")
	if err != nil {
		panic("failed to create casbin enforcer: " + err.Error())
	}
	config.InitializePermission(enforcer)
	router.AppRouter(app)

	return app
}

func CleanupTestDB() {
	sqlDB, _ := testDB.DB()
	sqlDB.Close()
}

func MakeRequest(app *fiber.App, method, url string, body interface{}, headers map[string]string) (*http.Response, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		return nil, nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, respBody, nil
}

func ParseJSON(t *testing.T, data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	assert.NoError(t, err, "Failed to parse JSON response")
}

func CreateTestUser(db *gorm.DB, email, password, role string) *user.User {
	usr := &user.User{
		Name:          "Test User",
		Email:         email,
		Password:      password,
		Role:          user.UserRole(role),
		EmailVerified: true,
	}
	db.Create(usr)
	return usr
}

func AssertSuccessResponse(t *testing.T, resp *http.Response, expectedCode int) {
	assert.Equal(t, expectedCode, resp.StatusCode, "Expected status code %d, got %d", expectedCode, resp.StatusCode)
}

func AssertErrorResponse(t *testing.T, resp *http.Response, expectedCode int, body []byte) {
	assert.Equal(t, expectedCode, resp.StatusCode)
	var errorResp map[string]interface{}
	json.Unmarshal(body, &errorResp)
	// Error responses contain "error" field
	assert.True(t, errorResp["error"] != nil || errorResp["message"] != nil, "Response should contain 'error' or 'message' field")
}
