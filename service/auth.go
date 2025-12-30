package service

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/repository"
)

type AuthService struct {
	userR *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		userR: repo,
	}
}

func (s *AuthService) Register(user *dto.RegisterRequest) error {
	if err := s.userR.ExistEmail(user.Email); err == nil {
		return &helper.BadRequestError{Message: "Email already exists", Order: "S1"}
	}

	password, err := helper.HashPassword(user.Password)
	if err != nil {
		return &helper.BadRequestError{Message: "Failed to hash password", Order: "S2"}
	}

	userEntity := user.ToEntity()
	userEntity.Password = password

	err = s.userR.Create(userEntity)
	if err != nil {
		return &helper.InternalServerError{Message: err.Error(), Order: "S3"}
	}
	return nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (resp *dto.LoginResponse, err error) {
	user, err := s.userR.FindByEmail(req.Email)
	if err != nil {
		return nil, &helper.BadRequestError{
			Message: "Email not registered!",
			Order:   "S1",
		}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &helper.BadRequestError{
			Message: "Password is wrong!",
			Order:   "S2",
		}
	}

	token, err := helper.GenerateJWT(dto.UserClaims{}.FromEntity(*user))
	if err != nil {
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}
	return &dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{}.FromEntity(*user),
	}, nil
}
