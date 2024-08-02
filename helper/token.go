package helper

import (
	"crypto/rsa"
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

func GetUserIDFormToken(c *fiber.Ctx) float64 {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(float64)
}

func GenerateJWT(user dto.UserClaims) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString(GetPrivateKey())
}
