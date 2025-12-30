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
	authz := middleware.LoadAuthzMiddleware()
	post.Get("", authz.RequiresPermissions([]string{"post:list"}), h.All)
	post.Post("", authz.RequiresPermissions([]string{"post:create"}), h.Create)
	post.Put("/:id", authz.RequiresPermissions([]string{"post:update"}), h.Update)
	post.Delete("/:id", authz.RequiresPermissions([]string{"post:delete"}), h.Delete)
	post.Get("/:id", authz.RequiresPermissions([]string{"post:read"}), h.GetByID)
}
