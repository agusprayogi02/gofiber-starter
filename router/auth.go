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
	auth := app.Group("/auth")
	auth.Post("/register", h.Register(enforcer))
	auth.Post("/login", h.Login)
	auth.Post("/refresh-token", h.RefreshToken)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)
	auth.Post("/verify-email", h.VerifyEmail)
	auth.Post("/resend-verification", h.ResendVerificationEmail)

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware()
	auth.Post("/logout", authMiddleware, h.Logout)
	auth.Post("/logout-all", authMiddleware, h.LogoutAll)
	auth.Post("/change-password", authMiddleware, h.ChangePassword)
	auth.Get("/sessions", authMiddleware, h.GetActiveSessions)
	auth.Delete("/sessions/:sessionId", authMiddleware, h.RevokeSession)

	// Profile routes
	profile := app.Use(authMiddleware)
	profile.Get("/profile", h.GetProfile)
	profile.Put("/profile", h.UpdateProfile)
	profile.Post("/profile/avatar", h.UpdateAvatar)

	// Preferences routes
	profile.Get("/preferences", h.GetPreferences)
	profile.Put("/preferences", h.UpdatePreferences)
}
