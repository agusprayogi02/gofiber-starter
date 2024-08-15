package handler

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/service"
	"starter-gofiber/variables"

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
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	posts, err := h.service.All(&params)
	if err != nil {
		return &helper.BadRequestError{
			Message: err.Error(),
		}
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
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	fileName, err := helper.UploadFile(c, file, variables.POST_PATH)
	if err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	post.Photo = fileName
	post.UserID = uint(helper.GetUserIDFormToken(c))

	if err := c.BodyParser(&post); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	rest, err := h.service.Create(&post)
	if err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		Data:       rest,
		StatusCode: fiber.StatusCreated,
		Message:    "Post created successfully",
	}, c)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	var post dto.PostUpdateRequest

	id, err := c.ParamsInt("id")
	if err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	post.ID = uint(id)
	file, err := c.FormFile("photo")
	if err != nil && err != fiber.ErrNotFound {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	if file != nil {
		fileName, err := helper.UploadFile(c, file, variables.POST_PATH)
		if err != nil {
			return &helper.UnprocessableEntityError{
				Message: err.Error(),
			}
		}
		post.Photo = &fileName
	}
	post.UserID = uint(helper.GetUserIDFormToken(c))

	if err := c.BodyParser(&post); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	rest, err := h.service.Update(&post)
	if err != nil {
		return &helper.BadRequestError{
			Message: err.Error(),
		}
	}

	return helper.Response(dto.ResponseResult{
		Data:       rest,
		StatusCode: fiber.StatusOK,
		Message:    "Post updated successfully",
	}, c)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		return &helper.BadRequestError{
			Message: err.Error(),
		}
	}
	return helper.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Post deleted successfully",
	}, c)
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
		}
	}

	rest, err := h.service.GetByID(uint(id))
	if err != nil {
		return &helper.BadRequestError{
			Message: err.Error(),
		}
	}

	return helper.Response(dto.ResponseResult{
		Data:       rest,
		StatusCode: fiber.StatusOK,
		Message:    "Post found successfully",
	}, c)
}
