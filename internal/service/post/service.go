package post

import (
	"math"

	"starter-gofiber/internal/domain/post"
	"starter-gofiber/internal/infrastructure/storage"
	"starter-gofiber/pkg/apierror"
	iValidator "starter-gofiber/pkg/validator"
	"starter-gofiber/variables"

	"github.com/go-playground/validator/v10"
)

type PostService struct {
	repo post.Repository
}

func NewPostService(repo post.Repository) post.Service {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) Create(req *post.PostRequest, userID uint) (*post.PostResponse, error) {
	var errors []*iValidator.IError

	err := iValidator.Validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el iValidator.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return nil, &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Data:    errors,
			Order:   "S1",
		}
	}

	postEntity := req.ToEntity()
	postEntity.UserID = userID

	err = s.repo.Create(&postEntity)
	if err != nil {
		return nil, &apierror.BadRequestError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	resp := post.PostResponse{}.FromEntity(postEntity)
	return &resp, nil
}

func (s *PostService) FindByID(id uint) (*post.PostResponse, error) {
	postEntity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "Post not found",
			Order:   "S1",
		}
	}

	resp := post.PostResponse{}.FromEntity(*postEntity)
	return &resp, nil
}

func (s *PostService) FindAll(page, limit int) ([]post.PostResponse, *post.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}

	var result []post.PostResponse
	for _, p := range posts {
		result = append(result, post.PostResponse{}.FromEntity(p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := &post.PaginationMeta{
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}

	return result, meta, nil
}

func (s *PostService) Update(id uint, req *post.PostUpdateRequest, userID uint) (*post.PostResponse, error) {
	var errors []*iValidator.IError

	err := iValidator.Validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el iValidator.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return nil, &apierror.UnprocessableEntityError{
			Message: err.Error(),
			Data:    errors,
			Order:   "S1",
		}
	}

	postEntity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "Post not found",
			Order:   "S2",
		}
	}

	// Check ownership
	if postEntity.UserID != userID {
		return nil, &apierror.ForbiddenError{
			Message: "You don't have permission to update this post",
			Order:   "S3",
		}
	}

	// Update fields
	postEntity.Tweet = req.Tweet
	if req.Photo != nil {
		// Delete old photo if exists
		if postEntity.Photo != nil {
			_ = storage.DeleteFile(postEntity.Photo, variables.POST_PATH)
		}
		postEntity.Photo = req.Photo
	}

	err = s.repo.Update(postEntity)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	resp := post.PostResponse{}.FromEntity(*postEntity)
	return &resp, nil
}

func (s *PostService) Delete(id uint, userID uint) error {
	postEntity, err := s.repo.FindByID(id)
	if err != nil {
		return &apierror.NotFoundError{
			Message: "Post not found",
			Order:   "S1",
		}
	}

	// Check ownership
	if postEntity.UserID != userID {
		return &apierror.ForbiddenError{
			Message: "You don't have permission to delete this post",
			Order:   "S2",
		}
	}

	// Delete photo if exists
	if postEntity.Photo != nil {
		_ = storage.DeleteFile(postEntity.Photo, variables.POST_PATH)
	}

	return s.repo.Delete(id)
}

func (s *PostService) FindByUserID(userID uint, page, limit int) ([]post.PostResponse, *post.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}

	var result []post.PostResponse
	for _, p := range posts {
		result = append(result, post.PostResponse{}.FromEntity(p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := &post.PaginationMeta{
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}

	return result, meta, nil
}
