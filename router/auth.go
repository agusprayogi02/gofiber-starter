package router

import (
	"starter-gofiber/config"
	"starter-gofiber/handler"
	"starter-gofiber/repository"
	"starter-gofiber/service"

	"github.com/gofiber/fiber/v2"
)

func NewAuthentication(app fiber.Router) {
	repo := repository.NewUserRepository(config.DB)
	s := service.NewAuthService(repo)
	h := handler.NewAuthHandler(s)

	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
}
