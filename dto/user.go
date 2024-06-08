package dto

import "starter-gofiber/entity"

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

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
