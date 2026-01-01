package http

import (
	"strconv"

	"starter-gofiber/dto"
	"starter-gofiber/internal/domain/post"
	"starter-gofiber/internal/infrastructure/storage"
	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/pkg/response"
	"starter-gofiber/variables"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service post.Service
}

func NewPostHandler(s post.Service) *PostHandler {
	return &PostHandler{
		service: s,
	}
}

func (h *PostHandler) All(c *fiber.Ctx) error {
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, meta, err := h.service.FindAll(page, limit)
	if err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Success",
		Data:       posts,
		Paginate: &dto.Pagination{
			Page:       page,
			PerPage:    limit,
			Total:      meta.Total,
			TotalPages: meta.TotalPages,
			NextPage:   meta.HasNextPage,
		},
	}, c)
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req post.PostRequest

	// Parse body first to get user input
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	// Get authenticated user - security critical
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	// Overwrite UserID from request with authenticated user ID (prevent authorization bypass)
	req.UserID = userClaims.ID

	// Handle file upload if present
	file, err := c.FormFile("photo")
	if err == nil {
		fileName, err := storage.UploadFile(c, file, variables.POST_PATH)
		if err != nil {
			return err
		}
		// Overwrite Photo from request with uploaded file name
		req.Photo = fileName
	}

	resp, err := h.service.Create(&req, userClaims.ID)
	if err != nil {
		return err
	}

	// Note: Cache invalidation needs to be implemented
	// helper.InvalidateCollection("posts")

	return response.Response(dto.ResponseResult{
		Data:       resp,
		StatusCode: fiber.StatusCreated,
		Message:    "Post created successfully",
	}, c)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	var req post.PostUpdateRequest

	// Parse ID from URL parameter
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	// Parse body first to get user input
	if err := c.BodyParser(&req); err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H2",
		}
	}

	// Get authenticated user - security critical
	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	// Overwrite security-critical fields from request with authenticated values
	// This prevents authorization bypass and ID manipulation
	req.ID = uint(id)
	req.UserID = userClaims.ID

	// Handle file upload if present
	file, err := c.FormFile("photo")
	if err == nil {
		fileName, err := storage.UploadFile(c, file, variables.POST_PATH)
		if err != nil {
			return err
		}
		// Overwrite Photo from request with uploaded file name
		req.Photo = &fileName
	}

	resp, err := h.service.Update(uint(id), &req, userClaims.ID)
	if err != nil {
		return err
	}

	// Note: Cache invalidation needs to be implemented
	// helper.InvalidateRelated("post", id)
	// helper.InvalidateCollection("posts")

	return response.Response(dto.ResponseResult{
		Data:       resp,
		StatusCode: fiber.StatusOK,
		Message:    "Post updated successfully",
	}, c)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	userClaims, err := crypto.GetUserFromToken(c)
	if err != nil {
		return err
	}

	err = h.service.Delete(uint(id), userClaims.ID)
	if err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		StatusCode: fiber.StatusOK,
		Message:    "Post deleted successfully",
	}, c)
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Order:   "H1",
		}
	}

	resp, err := h.service.FindByID(uint(id))
	if err != nil {
		return err
	}

	return response.Response(dto.ResponseResult{
		Data:       resp,
		StatusCode: fiber.StatusOK,
		Message:    "Post found successfully",
	}, c)
}
