package router

import (
	"starter-gofiber/config"
	"starter-gofiber/handler"
	"starter-gofiber/middleware"
	"starter-gofiber/repository"
	"starter-gofiber/service"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func NewAuthentication(app fiber.Router, enforcer *casbin.Enforcer) {
	userRepo := repository.NewUserRepository(config.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(config.DB)
	passwordResetRepo := repository.NewPasswordResetRepository(config.DB)
	emailVerifRepo := repository.NewEmailVerificationRepository(config.DB)

	s := service.NewAuthService(userRepo, refreshTokenRepo, passwordResetRepo, emailVerifRepo)
	h := handler.NewAuthHandler(s)

	// Public routes (no authentication required)
	app.Post("/register", h.Register(enforcer))
	app.Post("/login", h.Login)
	app.Post("/refresh-token", h.RefreshToken)
	app.Post("/forgot-password", h.ForgotPassword)
	app.Post("/reset-password", h.ResetPassword)
	app.Post("/verify-email", h.VerifyEmail)
	app.Post("/resend-verification", h.ResendVerificationEmail)

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware()
	app.Post("/logout", authMiddleware, h.Logout)
	app.Post("/logout-all", authMiddleware, h.LogoutAll)
	app.Post("/change-password", authMiddleware, h.ChangePassword)
	app.Get("/sessions", authMiddleware, h.GetActiveSessions)
	app.Delete("/sessions/:sessionId", authMiddleware, h.RevokeSession)
}
