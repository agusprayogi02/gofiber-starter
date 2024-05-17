package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		return c.SendString("Hello, " + c.Params("name") + "!")
	})

	app.Listen(":3000")
}
