package router

import (
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/config"
	"starter-gofiber/dto"
	"starter-gofiber/helper"
)

var Authz *casbin.Middleware

func AppRouter(app *fiber.App) {
	Authz = casbin.New(casbin.Config{
		ModelFilePath: "./asset/rbac/model.conf",
		PolicyAdapter: fileadapter.NewAdapter("./asset/rbac/policy.csv"),
	})
	static := app.Group(config.STATIC_PATH, Authz.RequiresRoles([]string{"admin"}))
	static.Static("/", "./public")

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(dto.SuccessResponse{
			Message:   "pong",
			Timestamp: helper.TimeNow(),
		})
	})

	NewUser(app)
}
