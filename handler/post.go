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
			Order:   "H1",
		}
	}

	posts, err := h.service.All(&params)
	if err != nil {
		return err
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
			Order:   "H1",
		}
	}

	fileName, err := helper.UploadFile(c, file, variables.POST_PATH)
	if err != nil {
		return err
	}

	post.Photo = fileName
	token, err := helper.GetUserFromToken(c)
	if err != nil {
		return err
	}
	post.UserID = token.ID

	if err := c.BodyParser(&post); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H2",
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
			Order:   "H1",
		}
	}

	post.ID = uint(id)
	file, err := c.FormFile("photo")
	if err == nil {
		fileName, err := helper.UploadFile(c, file, variables.POST_PATH)
		if err != nil {
			return err
		}
		post.Photo = &fileName
	}

	token, err := helper.GetUserFromToken(c)
	if err != nil {
		return err
	}
	post.UserID = token.ID

	if err := c.BodyParser(&post); err != nil {
		return &helper.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H2",
		}
	}

	rest, err := h.service.Update(&post)
	if err != nil {
		return err
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
			Order:   "H1",
		}
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		return err
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
			Order:   "H1",
		}
	}

	rest, err := h.service.GetByID(uint(id))
	if err != nil {
		return err
	}

	return helper.Response(dto.ResponseResult{
		Data:       rest,
		StatusCode: fiber.StatusOK,
		Message:    "Post found successfully",
	}, c)
}
