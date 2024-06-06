package middleware

import (
	"learn-gofiber/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(app *fiber.App) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    config.PRIVATE_KEY.Public(),
			JWTAlg: jwtware.RS256,
		},
	}))
}
