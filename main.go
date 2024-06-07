package main

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/config"
	"starter-gofiber/router"
)

func main() {
	config.LoadConfig()
	config.LoadPermissions()
	config.LoadStorage()
	config.LoadDB()

	app := fiber.New()
	config.App(app)
	router.AppRouter(app)

	err := app.Listen(":" + config.ENV.PORT)
	if err != nil {
		panic(err)
	}
}
