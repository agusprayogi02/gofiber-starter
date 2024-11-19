package middleware

import (
	"starter-gofiber/helper"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() func(*fiber.Ctx) error {
	privateKey := helper.GetPrivateKey()
	return jwtware.New(jwtware.Config{
		ContextKey: "user",
		SigningKey: jwtware.SigningKey{
			Key:    privateKey.Public(),
			JWTAlg: jwtware.RS256,
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return helper.ErrorHelper(c, &helper.UnauthorizedError{Message: err.Error()})
		},
	})
}
