package router

import (
	"starter-gofiber/internal/config"
	"starter-gofiber/internal/handler/http"
	"starter-gofiber/internal/handler/middleware"
	"starter-gofiber/internal/repository/postgres"
	"starter-gofiber/internal/service/auth"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func NewAuthentication(app fiber.Router, enforcer *casbin.Enforcer) {
	userRepo := postgres.NewUserRepository(config.DB)

	s := auth.NewAuthService(userRepo)
	h := http.NewAuthHandler(s)

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

	// Profile routes
	app.Get("/profile", authMiddleware, h.GetProfile)
	app.Put("/profile", authMiddleware, h.UpdateProfile)
	app.Post("/profile/avatar", authMiddleware, h.UpdateAvatar)

	// Preferences routes
	app.Get("/preferences", authMiddleware, h.GetPreferences)
	app.Put("/preferences", authMiddleware, h.UpdatePreferences)
}
