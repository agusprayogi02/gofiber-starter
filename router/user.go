package router

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/config"
	"starter-gofiber/handler"
	"starter-gofiber/repository"
)

func NewUser(app *fiber.App) {
	repo := repository.NewUser(config.DB)
	h := handler.NewUser(repo)

	app.Post("/register", h.Register)
}
