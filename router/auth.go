package router

import (
	"starter-gofiber/config"
	"starter-gofiber/handler"
	"starter-gofiber/repository"
	"starter-gofiber/service"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func NewAuthentication(app fiber.Router, enforcer *casbin.Enforcer) {
	repo := repository.NewUserRepository(config.DB)
	s := service.NewAuthService(repo)
	h := handler.NewAuthHandler(s)

	app.Post("/register", h.Register(enforcer))
	app.Post("/login", h.Login)
}
