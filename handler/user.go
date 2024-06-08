package handler

import (
	"github.com/gofiber/fiber/v2"
	"starter-gofiber/config"
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

	password, err := helper.HashPassword(user.Password)
	if err != nil {
		return helper.ErrorHelper(c, &helper.BadRequestError{Message: "Failed to hash password"})
	}

	userEntity := entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: password,
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

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var userReq *dto.LoginRequest
	if err := c.BodyParser(&userReq); err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
		})
	}

	user, err := h.userRepo.FindByEmail(userReq.Email)
	if err != nil {
		return helper.ErrorHelper(c, &helper.BadRequestError{
			Message: "Email or password is wrong!",
		})
	}

	if err := helper.VerifyPassword(userReq.Password, user.Password); err != nil {
		return helper.ErrorHelper(c, &helper.BadRequestError{
			Message: err.Error(),
		})
	}

	token, err := helper.GenerateJWT(c, dto.UserClaims{
		ID:    user.ID,
		Role:  user.Role.String(),
		Email: user.Email,
	})

	res := helper.Response(dto.ResponseParams{
		StatusCode: fiber.StatusOK,
		Message:    "Login Success",
		Data: dto.LoginResponse{
			User: dto.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				Role:      user.Role.String(),
				CreatedAt: user.CreatedAt.Format(config.FORMAT_TIME),
				UpdatedAt: user.UpdatedAt.Format(config.FORMAT_TIME),
			},
			Token: token,
		},
	})
	return c.Status(fiber.StatusOK).JSON(res)
}
