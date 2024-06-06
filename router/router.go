package router

import (
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/config"
)

var Authz *casbin.Middleware

func AppRouter(app *fiber.App) {
	Authz = casbin.New(casbin.Config{
		ModelFilePath: "./asset/rbac/model.conf",
		PolicyAdapter: fileadapter.NewAdapter("./asset/rbac/policy.csv"),
	})
	static := app.Group(config.STATIC_PATH, Authz.RequiresRoles([]string{"admin"}))
	static.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server is running",
		})
	})

	NewUser(app)
}
