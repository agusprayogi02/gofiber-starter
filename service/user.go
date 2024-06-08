package service

import (
	"starter-gofiber/dto"
	"starter-gofiber/helper"
	"starter-gofiber/repository"
)

type UserService struct {
	userR *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userR: repo,
	}
}

func (s *UserService) Register(user *dto.RegisterRequest) error {
	if err := s.userR.ExistEmail(user.Email); err != nil {
		return &helper.BadRequestError{Message: "Email already exists"}
	}

	password, err := helper.HashPassword(user.Password)
	if err != nil {
		return &helper.BadRequestError{Message: "Failed to hash password"}
	}

	userEntity := user.ToEntity()
	userEntity.Password = password

	err = s.userR.Register(userEntity)
	if err != nil {
		return &helper.InternalServerError{Message: err.Error()}
	}
	return nil
}

func (s *UserService) Login(req *dto.LoginRequest) (resp *dto.LoginResponse, err error) {
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

	token, err := helper.GenerateJWT(dto.UserClaims{
		ID:    user.ID,
		Role:  user.Role.String(),
		Email: user.Email,
	})
	resp = &dto.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}
	return resp, nil
}
