package middleware

import (
	"starter-gofiber/helper"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

func LoadAuthzMiddleware() *casbin.Middleware {
	return casbin.New(casbin.Config{
		ModelFilePath: "./assets/rbac/model.conf",
		PolicyAdapter: fileadapter.NewAdapter("./assets/rbac/policy.csv"),
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
