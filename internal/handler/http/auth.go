package http

import (
	"strconv"

	"starter-gofiber/dto"
	"starter-gofiber/internal/domain/user"
	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userS user.Service
}

func NewAuthHandler(s user.Service) *AuthHandler {
	return &AuthHandler{
		userS: s,
	}
}

func (h *AuthHandler) Register(enforcer *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req *user.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return &apierror.UnprocessableEntityError{Message: err.Error(), Order: "H1"}
		}

		if err := h.userS.Register(req); err != nil {
			return err
		}
		ok, err := enforcer.AddRoleForUser(req.Email, req.Role)
		if !ok || err != nil {
			if err != nil {
				return &apierror.UnprocessableEntityError{Message: err.Error(), Order: "H2"}
			}
			return &apierror.InternalServerError{Message: "Failed to assign role to user", Order: "H2"}
		}

		return response.Response(dto.ResponseResult{
			StatusCode: fiber.StatusCreated,
			Message:    "User registered successfully",
		}, c)
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var userReq *user.LoginRequest
	if err := c.BodyParser(&userReq); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	// Get IP and User Agent for session tracking
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	loginResp, err := h.userS.Login(userReq, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Login Success",
		Data:       loginResp,
	}, c)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req *user.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
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

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Token refreshed successfully",
		Data:       tokens,
	}, c)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req *user.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.Logout(req.RefreshToken); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Logout successful",
	}, c)
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	if err := h.userS.LogoutAll(userClaims.ID); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Logged out from all devices",
	}, c)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req *user.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ForgotPassword(req); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "If email exists, password reset link has been sent",
	}, c)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req *user.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ResetPassword(req); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Password reset successful",
	}, c)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	var req *user.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ChangePassword(userClaims.ID, req); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Password changed successfully",
	}, c)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var req *user.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.VerifyEmail(req); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Email verified successfully",
	}, c)
}

func (h *AuthHandler) ResendVerificationEmail(c *fiber.Ctx) error {
	var req *user.ForgotPasswordRequest // Reuse same DTO (only needs email)
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.ResendVerificationEmail(req.Email); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Verification email sent",
	}, c)
}

func (h *AuthHandler) GetActiveSessions(c *fiber.Ctx) error {
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	sessions, err := h.userS.GetActiveSessions(userClaims.ID)
	if err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Active sessions retrieved",
		Data:       sessions,
	}, c)
}

func (h *AuthHandler) RevokeSession(c *fiber.Ctx) error {
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	sessionIDStr := c.Params("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	if err := h.userS.RevokeSession(uint(sessionID), userClaims.ID); err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Session revoked successfully",
	}, c)
}
