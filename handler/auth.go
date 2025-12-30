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

	user, err := h.userS.Login(userReq)
	if err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Login Success",
		Data:       user,
	}, c)
}
