package handler

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/service"
	"starter-gofiber/variables"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{
		service: s,
	}
}

func (h *PostHandler) All(c *fiber.Ctx) error {
	params := dto.Pagination{
		Page:    1,
		PerPage: 10,
	}

	if err := c.ParamsParser(&params); err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
		})
	}

	posts, err := h.service.All(&params)
	if err != nil {
		return helper.ErrorHelper(c, &helper.BadRequestError{
			Message: err.Error(),
		})
	}

	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Success",
		Paginate:   &params,
		Data:       posts,
	}, c)
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	var post dto.PostRequest

	file, err := c.FormFile("photo")
	if err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
		})
	}

	fileName, err := helper.UploadFile(c, file, variables.POST_PATH)
	if err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
		})
	}

	post.Photo = fileName
	post.UserID = uint(helper.GetUserIDFormToken(c))

	var errors []*helper.IError
	if err := c.BodyParser(&post); err != nil {
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
		})
	}

	err = helper.Validator.Struct(post)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el helper.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return helper.ErrorHelper(c, &helper.UnprocessableEntityError{
			Message: err.Error(),
			Data:    errors,
		})
	}

	rest, err := h.service.Create(&post)
	if err != nil {
		return helper.ErrorHelper(c, &helper.BadRequestError{
			Message: err.Error(),
		})
	}

	return helper.Response(dto.ResponseResult{
		Data:       rest,
		StatusCode: fiber.StatusCreated,
		Message:    "Post created successfully",
	}, c)
}
