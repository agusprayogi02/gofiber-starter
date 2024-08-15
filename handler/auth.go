package handler

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/service"

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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user *dto.RegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return &helper.UnprocessableEntityError{Message: err.Error()}
	}

	if err := h.userS.Register(user); err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusCreated,
		Message:    "User registered successfully",
	}, c)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var userReq *dto.LoginRequest
	if err := c.BodyParser(&userReq); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
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
