package service

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/repository"

	"github.com/go-playground/validator/v10"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) All(paginate *dto.Pagination) (*[]dto.PostResponse, error) {
	posts, err := s.repo.All(paginate)
	if err != nil {
		return nil, err
	}

	var result []dto.PostResponse
	for _, post := range *posts {
		result = append(result, dto.PostResponse{}.FromEntity(post))
	}
	return &result, nil
}

func (s *PostService) Create(post *dto.PostRequest) (*dto.PostResponse, error) {
	var errors []*helper.IError

	err := helper.Validator.Struct(post)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el helper.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return nil, &helper.UnprocessableEntityError{
			Message: err.Error(),
			Data:    errors,
		}
	}

	rest, err := s.repo.Create(post.ToEntity())
	r := dto.PostResponse{}.FromEntity(rest)
	return &r, &helper.BadRequestError{
		Message: err.Error(),
	}
}

func (s *PostService) GetByID(id uint) (dto.PostResponse, error) {
	rest, err := s.repo.FindId(id)
	return dto.PostResponse{}.FromEntity(*rest), err
}

func (s *PostService) Update(post *dto.PostUpdateRequest) (*dto.PostResponse, error) {
	var errors []*helper.IError

	err := helper.Validator.Struct(post)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el helper.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return nil, &helper.UnprocessableEntityError{
			Message: err.Error(),
			Data:    errors,
		}
	}
	rest, err := s.repo.FindId(post.ID)
	if err != nil {
		return nil, &helper.BadRequestError{
			Message: "This item does not exist",
		}
	}

	rest.Tweet = post.Tweet
	if post.Photo != nil {
		rest.Photo = post.Photo
	}
	rest.UserID = post.UserID

	err = s.repo.Update(*rest, post.ID)
	r := dto.PostResponse{}.FromEntity(*rest)
	return &r, &helper.BadRequestError{
		Message: err.Error(),
	}
}

func (s *PostService) Delete(id uint) error {
	return s.repo.Delete(id)
}
