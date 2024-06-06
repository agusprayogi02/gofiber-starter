package main

import (
	"learn-gofiber/config"
	"learn-gofiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()
	config.LoadStorage()
	config.LoadDB()

	app := fiber.New()
	config.App(app)
	middleware.AuthMiddleware(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		return c.SendString("Hello, " + c.Params("name") + "!")
	})
	app.Static("/", "./public")

	app.Listen(":3000")
}
