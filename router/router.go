package router

import (
	"fmt"

	"starter-gofiber/config"
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/middleware"
	"starter-gofiber/variables"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func AppRouter(app *fiber.App) {
	adapter, err := xormadapter.NewAdapter("sqlite3", fmt.Sprintf("./asset/%s_storage.db", config.ENV.DB_NAME), true)
	if err != nil {
		panic(err)
	}
	enforcer, err := casbin.NewEnforcer("./asset/rbac/model.conf", adapter)
	if err != nil {
		panic(err)
	}
	err = config.InitializePermission(enforcer)
	if err != nil {
		panic(err)
	}

	if config.ENV.ENV_TYPE == "dev" {
		app.Use(logger.New())
	}
	static := app.Group(variables.STATIC_PATH, middleware.AuthMiddleware())
	authz := middleware.LoadAuthzMiddleware()
	static.Use(authz.RequiresPermissions([]string{"files:read"}))
	static.Static("/", "./public")
	app.Static("/favicon.ico", "./public/favicon.ico")
	if config.ENV.ENV_TYPE != "dev" {
		app.Use(recover.New())
	}

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
			Code:      fiber.StatusOK,
			Message:   "pong",
			Timestamp: helper.TimeNow(),
		})
	})

	api := app.Group("/api")
	NewAuthentication(api, enforcer)
	NewPostRouter(api)
}
