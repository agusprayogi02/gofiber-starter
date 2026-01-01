package http

import (
	"starter-gofiber/internal/config"
	"starter-gofiber/internal/handler/middleware"
	"starter-gofiber/internal/repository/postgres"
	"starter-gofiber/internal/service/auth"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func NewAuthentication(app fiber.Router, enforcer *casbin.Enforcer) {
	userRepo := postgres.NewUserRepository(config.DB)

	authService := auth.NewAuthService(userRepo)
	authHandler := NewAuthHandler(authService)

	// Public routes (no authentication required)
	app.Post("/register", authHandler.Register(enforcer))
	app.Post("/login", authHandler.Login)
	app.Post("/refresh-token", authHandler.RefreshToken)
	app.Post("/forgot-password", authHandler.ForgotPassword)
	app.Post("/reset-password", authHandler.ResetPassword)
	app.Post("/verify-email", authHandler.VerifyEmail)
	app.Post("/resend-verification", authHandler.ResendVerificationEmail)

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware()
	app.Post("/logout", authMiddleware, authHandler.Logout)
	app.Post("/logout-all", authMiddleware, authHandler.LogoutAll)
	app.Post("/change-password", authMiddleware, authHandler.ChangePassword)
	app.Get("/sessions", authMiddleware, authHandler.GetActiveSessions)
	app.Delete("/sessions/:sessionId", authMiddleware, authHandler.RevokeSession)
}
