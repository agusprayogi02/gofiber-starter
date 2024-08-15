package main

import (
	"starter-gofiber/config"
	"starter-gofiber/helper"
	"starter-gofiber/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig() // required first, because it will load .env file

	config.LoadTimezone()
	config.LoadPermissions()
	config.LoadStorage()
	config.LoadDB()
	if config.ENV.DB_2_ENABLE {
		config.LoadDB2()
	}

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: helper.ErrorHelper,
	})
	config.App(app)
	router.AppRouter(app)

	err := app.Listen(":" + config.ENV.PORT)
	if err != nil {
		panic(err)
	}
}
