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
