package middleware

import (
	"starter-gofiber/config"
	"starter-gofiber/helper"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

func LoadAuthzMiddleware() *casbin.Middleware {
	modelPath := "./assets/rbac/model.conf"
	policyPath := "./assets/rbac/policy.csv"
	if config.ENV.ENV_TYPE == "test" {
		modelPath = "../assets/rbac/model.conf"
		policyPath = "../assets/rbac/policy.csv"
	}
	return casbin.New(casbin.Config{
		ModelFilePath: modelPath,
		PolicyAdapter: fileadapter.NewAdapter(policyPath),
		Lookup: func(c *fiber.Ctx) string {
			token, err := helper.GetUserFromToken(c)
			if err != nil {
				return ""
			}
			return token.Email
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return &helper.UnauthorizedError{
				Message: "Harus Login terlebih dahulu",
				Order:   "M-casbin-authz",
			}
		},
		Forbidden: func(c *fiber.Ctx) error {
			return &helper.ForbiddenError{
				Message: "Tidak ada hak akses",
				Order:   "M-casbin-authz",
			}
		},
	})
}
