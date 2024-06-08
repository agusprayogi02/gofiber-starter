package router

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/config"
	"starter-gofiber/handler"
	"starter-gofiber/repository"
	"starter-gofiber/service"
)

func NewUser(app fiber.Router) {
	repo := repository.NewUserRepository(config.DB)
	s := service.NewUserService(repo)
	h := handler.NewUserHandler(s)

	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
}
