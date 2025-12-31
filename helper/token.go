package helper

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"starter-gofiber/config"
	"starter-gofiber/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

func GetPrivateKey() *rsa.PrivateKey {
	privateKeyData, err := os.ReadFile(config.ENV.LOCATION_CERT)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
		panic(err)
	}

	PK, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
		panic(err)
	}
	return PK
}

func GetUserFromToken(c *fiber.Ctx) (*dto.CustomClaims, error) {
	token := c.Locals("user").(*jwt.Token)

	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := dto.CustomClaims{}.FromToken(claim)
		return &data, nil
	} else {
		return nil, &UnauthorizedError{
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
