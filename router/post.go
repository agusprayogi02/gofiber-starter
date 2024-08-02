package router

import (
	"starter-gofiber/config"
	"starter-gofiber/handler"
	"starter-gofiber/middleware"
	"starter-gofiber/repository"
	"starter-gofiber/service"

	"github.com/gofiber/fiber/v2"
)

func NewPostRouter(app fiber.Router) {
	repo := repository.NewPostRepository(config.DB)
	s := service.NewPostService(repo)
	h := handler.NewPostHandler(s)

	post := app.Group("/post")

	// JWT Middleware
	post.Use(middleware.AuthMiddleware())
	post.Get("", h.All)
	post.Post("", h.Create)
	post.Put("/:id", h.Update)
	post.Delete("/:id", h.Delete)
	post.Get("/:id", h.GetByID)
}
