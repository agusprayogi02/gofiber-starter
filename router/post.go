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

	posts := app.Group("/posts")

	// Public routes
	posts.Get("", h.All)
	posts.Get("/:id", h.GetByID)

	// Protected routes with JWT and authorization
	authMiddleware := middleware.AuthMiddleware()
	authz := middleware.LoadAuthzMiddleware()
	posts.Post("", authMiddleware, authz.RequiresPermissions([]string{"post:create"}), h.Create)
	posts.Put("/:id", authMiddleware, authz.RequiresPermissions([]string{"post:update"}), h.Update)
	posts.Delete("/:id", authMiddleware, authz.RequiresPermissions([]string{"post:delete"}), h.Delete)
}
