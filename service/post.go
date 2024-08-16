package service

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/repository"
	"starter-gofiber/variables"

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
		return nil, &helper.BadRequestError{
			Message: err.Error(),
			Order:   "S1",
		}
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
			Order:   "S1",
		}
	}

	rest, err := s.repo.Create(post.ToEntity())
	r := dto.PostResponse{}.FromEntity(rest)

	if err != nil {
		return nil, &helper.BadRequestError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	return &r, nil
}

func (s *PostService) GetByID(id uint) (dto.PostResponse, error) {
	rest, err := s.repo.FindId(id)
	if err != nil {
		return dto.PostResponse{}, &helper.BadRequestError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return dto.PostResponse{}.FromEntity(*rest), nil
}

func (s *PostService) Update(upp *dto.PostUpdateRequest) (*dto.PostResponse, error) {
	var errors []*helper.IError

	err := helper.Validator.Struct(upp)
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
			Order:   "S1",
		}
	}
	post, err := s.repo.FindId(upp.ID)
	if err != nil {
		return nil, &helper.BadRequestError{
			Message: "This item does not exist",
			Order:   "S2",
		}
	}

	// delete old photo
	if upp.Photo != nil {
		if err := helper.DeleteFile(post.Photo, variables.POST_PATH); err != nil {
			return nil, err
		}
		post.Photo = upp.Photo
	}
	post.Tweet = upp.Tweet
	if upp.Photo != nil {
		post.Photo = upp.Photo
	}
	post.UserID = upp.UserID

	err = s.repo.Update(*post, upp.ID)
	r := dto.PostResponse{}.FromEntity(*post)

	if err != nil {
		return nil, &helper.BadRequestError{
			Message: err.Error(),
			Order:   "S3",
		}
	}
	return &r, nil
}

func (s *PostService) Delete(id uint) error {
	post, err := s.repo.FindId(id)
	if err != nil {
		return &helper.BadRequestError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	if err := helper.DeleteFile(post.Photo, variables.POST_PATH); err != nil {
		return &helper.BadRequestError{
			Message: err.Error(),
			Order:   "S2",
		}
	}
	return s.repo.Delete(id)
}
