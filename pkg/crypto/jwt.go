package crypto

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"starter-gofiber/dto"
	"starter-gofiber/pkg/apierror"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var privateKey *rsa.PrivateKey

// InitPrivateKey initializes the RSA private key
func InitPrivateKey(certLocation string) error {
	privateKeyData, err := os.ReadFile(certLocation)
	if err != nil {
		return err
	}

	pk, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return err
	}

	privateKey = pk
	return nil
}

func GetPrivateKey() *rsa.PrivateKey {
	if privateKey == nil {
		log.Fatal("Private key not initialized")
	}
	return privateKey
}

func GetUserFromToken(c *fiber.Ctx) (*dto.CustomClaims, error) {
	token := c.Locals("user").(*jwt.Token)

	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := dto.CustomClaims{}.FromToken(claim)
		return &data, nil
	} else {
		return nil, &apierror.UnauthorizedError{
			Message: errors.New("Token Tidak Valid " + token.Raw).Error(),
		}
	}
}

func GenerateJWT(user dto.UserClaims) (string, error) {
	// Create the Claims
	claims := dto.CustomClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Short-lived: 1 hour
			Issuer:    "Starter-Gofiber",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString(GetPrivateKey())
}

func GenerateRefreshToken(user dto.UserClaims) (string, error) {
	// Create the Claims for refresh token
	claims := dto.CustomClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)), // Long-lived: 30 days
			Issuer:    "Starter-Gofiber-Refresh",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	// Generate encoded token
	return token.SignedString(GetPrivateKey())
}

func GenerateRandomToken() (string, error) {
	// Create token with random claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Issuer:    "Starter-Gofiber-Token",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	return token.SignedString(GetPrivateKey())
}
