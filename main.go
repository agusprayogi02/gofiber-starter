package main

import (
	"starter-gofiber/config"
	"starter-gofiber/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()
	config.LoadPermissions()
	config.LoadStorage()
	config.LoadDB()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	config.App(app)
	router.AppRouter(app)

	err := app.Listen(":" + config.ENV.PORT)
	if err != nil {
		panic(err)
	}
}
