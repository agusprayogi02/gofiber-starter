package http

import (
	"starter-gofiber/internal/config"
	"starter-gofiber/internal/handler/middleware"
	"starter-gofiber/internal/repository/postgres"
	"starter-gofiber/internal/service/post"

	"github.com/gofiber/fiber/v2"
)

func NewPostRouter(app fiber.Router) {
	repo := postgres.NewPostRepository(config.DB)
	postService := post.NewPostService(repo)
	postHandler := NewPostHandler(postService)

	posts := app.Group("/posts")

	// Public routes
	posts.Get("", postHandler.All)
	posts.Get("/:id", postHandler.GetByID)

	// Protected routes with JWT and authorization
	authMiddleware := middleware.AuthMiddleware()
	authz := middleware.LoadAuthzMiddleware()
	posts.Post("", authMiddleware, authz.RequiresPermissions([]string{"post:create"}), postHandler.Create)
	posts.Put("/:id", authMiddleware, authz.RequiresPermissions([]string{"post:update"}), postHandler.Update)
	posts.Delete("/:id", authMiddleware, authz.RequiresPermissions([]string{"post:delete"}), postHandler.Delete)
}
