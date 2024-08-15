package router

import (
	"starter-gofiber/config"
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/variables"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var Authz *casbin.Middleware

func AppRouter(app *fiber.App) {
	Authz = casbin.New(casbin.Config{
		ModelFilePath: "./asset/rbac/model.conf",
		PolicyAdapter: fileadapter.NewAdapter("./asset/rbac/policy.csv"),
	})
	if config.ENV.ENV_TYPE == "dev" {
		app.Use(logger.New())
	}
	static := app.Group(variables.STATIC_PATH, Authz.RequiresRoles([]string{variables.ADMIN_ROLE, variables.USER_ROLE}))
	static.Static("/", "./public")
	app.Static("/favicon.ico", "./public/favicon.ico")
	app.Use(recover.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
			Code:      fiber.StatusOK,
			Message:   "pong",
			Timestamp: helper.TimeNow(),
		})
	})

	api := app.Group("/api")
	NewAuthentication(api)
	NewPostRouter(api)
}
