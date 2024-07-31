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
		return &helper.BadRequestError{Message: "Email already exists"}
	}

	password, err := helper.HashPassword(user.Password)
	if err != nil {
		return &helper.BadRequestError{Message: "Failed to hash password"}
	}

	userEntity := user.ToEntity()
	userEntity.Password = password

	err = s.userR.Create(userEntity)
	if err != nil {
		return &helper.InternalServerError{Message: err.Error()}
	}
	return nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (resp *dto.LoginResponse, err error) {
	user, err := s.userR.FindByEmail(req.Email)
	if err != nil {
		return nil, &helper.BadRequestError{
			Message: "Email not registered!",
		}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &helper.BadRequestError{
			Message: "Password is wrong!",
		}
	}

	token, err := helper.GenerateJWT(dto.UserClaims{}.FromEntity(*user))
	resp = &dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{}.FromEntity(*user),
	}
	return resp, nil
}
