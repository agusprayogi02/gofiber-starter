package middleware

import (
	"os"

	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

func LoadAuthzMiddleware() *casbin.Middleware {
	modelPath := "./assets/rbac/model.conf"
	policyPath := "./assets/rbac/policy.csv"

	// Check if running in test mode
	if os.Getenv("ENV_TYPE") == "test" {
		modelPath = "../assets/rbac/model.conf"
		policyPath = "../assets/rbac/policy.csv"
	}

	return casbin.New(casbin.Config{
		ModelFilePath: modelPath,
		PolicyAdapter: fileadapter.NewAdapter(policyPath),
		Lookup: func(c *fiber.Ctx) string {
			token, err := crypto.GetUserFromToken(c)
			if err != nil {
				return ""
			}
			return token.Email
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return &apierror.UnauthorizedError{
				Message: "Harus Login terlebih dahulu",
				Order:   "M-casbin-authz",
			}
		},
		Forbidden: func(c *fiber.Ctx) error {
			return &apierror.ForbiddenError{
				Message: "Tidak ada hak akses",
				Order:   "M-casbin-authz",
			}
		},
	})
}
