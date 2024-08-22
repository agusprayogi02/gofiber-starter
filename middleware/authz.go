package middleware

import (
	"fmt"

	"starter-gofiber/config"
	"starter-gofiber/helper"

	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

func LoadAuthzMiddleware() *casbin.Middleware {
	adapter, err := xormadapter.NewAdapter("sqlite3", fmt.Sprintf("./assets/%s_storage.db", config.ENV.DB_NAME), true)
	if err != nil {
		panic(err)
	}
	return casbin.New(casbin.Config{
		ModelFilePath: "./assets/rbac/model.conf",
		PolicyAdapter: adapter,
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
