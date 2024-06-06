package handler

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/dto"
	"starter-gofiber/entity"
	"starter-gofiber/helper"
	"starter-gofiber/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUser(r *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: r,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var user *dto.RegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{Message: err.Error()})
	}

	if h.userRepo.ExistEmail(user.Email) {
		return helper.ErrorHelper(c, &helper.BadRequestError{Message: "Email already exists"})
	}

	userEntity := entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := h.userRepo.Register(userEntity); err != nil {
		return helper.ErrorHelper(c, &helper.InternalServerError{Message: err.Error()})
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: fiber.StatusCreated,
		Message:    "User registered successfully",
	})

	return c.Status(fiber.StatusCreated).JSON(res)
}
