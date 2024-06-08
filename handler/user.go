package handler

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/service"
)

type UserHandler struct {
	userS *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		userS: s,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var user *dto.RegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{Message: err.Error()})
	}

	if err := h.userS.Register(user); err != nil {
		return helper.ErrorHelper(c, err)
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: fiber.StatusCreated,
		Message:    "User registered successfully",
	})

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var userReq *dto.LoginRequest
	if err := c.BodyParser(&userReq); err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
		})
	}

	user, err := h.userS.Login(userReq)
	if err != nil {
		return helper.ErrorHelper(c, err)
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: fiber.StatusOK,
		Message:    "Login Success",
		Data:       user,
	})
	return c.Status(fiber.StatusOK).JSON(res)
}
