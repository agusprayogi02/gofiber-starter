package router

import (
	"starter-gofiber/config"
	"starter-gofiber/dto"
	"starter-gofiber/helper"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

var Authz *casbin.Middleware

func AppRouter(app *fiber.App) {
	Authz = casbin.New(casbin.Config{
		ModelFilePath: "./asset/rbac/model.conf",
		PolicyAdapter: fileadapter.NewAdapter("./asset/rbac/policy.csv"),
	})
	static := app.Group(config.STATIC_PATH, Authz.RequiresRoles([]string{config.ADMIN_ROLE, config.USER_ROLE}))
	static.Static("/", "./public")
	app.Static("/favicon.ico", "./public/favicon.ico")

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(dto.SuccessResponse{
			Message:   "pong",
			Timestamp: helper.TimeNow(),
		})
	})

	api := app.Group("/api")
	NewAuthentication(api)
}
