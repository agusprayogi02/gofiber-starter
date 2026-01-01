package middleware

import (
	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() func(*fiber.Ctx) error {
	privateKey := crypto.GetPrivateKey()
	return jwtware.New(jwtware.Config{
		ContextKey: "user",
		SigningKey: jwtware.SigningKey{
			Key:    privateKey.Public(),
			JWTAlg: jwtware.RS256,
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return apierror.ErrorHelper(c, &apierror.UnauthorizedError{Message: err.Error()})
		},
	})
}
