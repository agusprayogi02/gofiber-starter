package service

import (
	"starter-gofiber/dto"
	"starter-gofiber/repository"
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

func (s *PostService) Create(post *dto.PostRequest) (dto.PostResponse, error) {
	rest, err := s.repo.Create(post.ToEntity())
	return dto.PostResponse{}.FromEntity(rest), err
}

func (s *PostService) GetByID(id uint) (dto.PostResponse, error) {
	rest, err := s.repo.FindId(id)
	return dto.PostResponse{}.FromEntity(*rest), err
}

func (s *PostService) Update(post *dto.PostUpdateRequest) (dto.PostResponse, error) {
	rest, err := s.repo.FindId(post.ID)
	if err != nil {
		return dto.PostResponse{}, err
	}

	rest.Tweet = post.Tweet
	if post.Photo != nil {
		rest.Photo = post.Photo
	}
	rest.UserID = post.UserID

	err = s.repo.Update(*rest, post.ID)
	return dto.PostResponse{}.FromEntity(*rest), err
}

func (s *PostService) Delete(id uint) error {
	return s.repo.Delete(id)
}
