package dto

import (
	"starter-gofiber/config"
	"starter-gofiber/entity"
)

type UserClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required;email"`
	Password string `json:"password" binding:"required;min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required;min=3"`
	Email    string `json:"email" binding:"required;email"`
	Password string `json:"password" binding:"required;min=6"`
	Role     string `json:"role" binding:"oneof=admin user;default:user"`
}

func (r RegisterRequest) ToEntity() entity.User {
	return entity.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Role:     entity.UserRole(r.Role),
	}
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r UserResponse) FromEntity(u entity.User) UserResponse {
	r.ID = u.ID
	r.Name = u.Name
	r.Email = u.Email
	r.Role = u.Role.String()
	r.CreatedAt = u.CreatedAt.Format(config.FORMAT_TIME)
	r.UpdatedAt = u.UpdatedAt.Format(config.FORMAT_TIME)
	return r
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
