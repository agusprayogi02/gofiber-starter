package helper

import (
	"time"

	"learn-gofiber/config"
	"learn-gofiber/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(c *fiber.Ctx, user dto.UserClaims) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString(config.PRIVATE_KEY)
}
