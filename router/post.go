package router

import (
	"starter-gofiber/internal/config"
	"starter-gofiber/internal/handler/http"
	"starter-gofiber/internal/handler/middleware"
	"starter-gofiber/internal/repository/postgres"
	"starter-gofiber/internal/service/post"

	"github.com/gofiber/fiber/v2"
)

func NewPostRouter(app fiber.Router) {
	repo := postgres.NewPostRepository(config.DB)
	s := post.NewPostService(repo)
	h := http.NewPostHandler(s)

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
