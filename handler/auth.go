package handler

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/service"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userS *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		userS: s,
	}
}

func (h *AuthHandler) Register(enforcer *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user *dto.RegisterRequest
		if err := c.BodyParser(&user); err != nil {
			return &helper.UnprocessableEntityError{Message: err.Error(), Order: "H1"}
		}

		if err := h.userS.Register(user); err != nil {
			return err
		}
		if ok, err := enforcer.AddRoleForUser(user.Email, user.Role); ok && err != nil {
			return &helper.UnprocessableEntityError{Message: err.Error(), Order: "H2"}
		}

		return helper.Response(dto.ResponseResult{
			StatusCode: fiber.StatusCreated,
			Message:    "User registered successfully",
		}, c)
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var userReq *dto.LoginRequest
	if err := c.BodyParser(&userReq); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	// Get IP and User Agent for session tracking
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	user, err := h.userS.Login(userReq, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Login Success",
		Data:       user,
	}, c)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req *dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	tokens, err := h.userS.RefreshToken(req, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Token refreshed successfully",
		Data:       tokens,
	}, c)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req *dto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.Logout(req.RefreshToken); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Logout successful",
	}, c)
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	user, err := helper.GetUserFromToken(c)
	if err != nil {
		return err
	}

	if err := h.userS.LogoutAll(user.ID); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Logged out from all devices",
	}, c)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req *dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ForgotPassword(req); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "If email exists, password reset link has been sent",
	}, c)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req *dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ResetPassword(req); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Password reset successful",
	}, c)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	user, err := helper.GetUserFromToken(c)
	if err != nil {
		return err
	}

	var req *dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ChangePassword(user.ID, req); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Password changed successfully",
	}, c)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var req *dto.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.VerifyEmail(req); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Email verified successfully",
	}, c)
}

func (h *AuthHandler) ResendVerificationEmail(c *fiber.Ctx) error {
	var req *dto.ForgotPasswordRequest // Reuse same DTO (only needs email)
	if err := c.BodyParser(&req); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ResendVerificationEmail(req.Email); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Verification email sent",
	}, c)
}

func (h *AuthHandler) GetActiveSessions(c *fiber.Ctx) error {
	user, err := helper.GetUserFromToken(c)
	if err != nil {
		return err
	}

	sessions, err := h.userS.GetActiveSessions(user.ID)
	if err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Active sessions retrieved",
		Data:       sessions,
	}, c)
}

func (h *AuthHandler) RevokeSession(c *fiber.Ctx) error {
	user, err := helper.GetUserFromToken(c)
	if err != nil {
		return err
	}

	sessionID, err := c.ParamsInt("sessionId")
	if err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.RevokeSession(uint(sessionID), user.ID); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Session revoked successfully",
	}, c)
}
