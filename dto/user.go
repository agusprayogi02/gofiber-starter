package dto

import (
	"starter-gofiber/internal/domain/user"
	"starter-gofiber/variables"
)

type UserClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (r UserClaims) FromEntity(u user.User) UserClaims {
	r.ID = u.ID
	r.Email = u.Email
	r.Role = u.Role.String()
	return r
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

func (r RegisterRequest) ToEntity() user.User {
	return user.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Role:     user.UserRole(r.Role),
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

func (r UserResponse) FromEntity(u user.User) UserResponse {
	r.ID = u.ID
	r.Name = u.Name
	r.Email = u.Email
	r.Role = u.Role.String()
	r.CreatedAt = u.CreatedAt.Format(variables.FORMAT_TIME)
	r.UpdatedAt = u.UpdatedAt.Format(variables.FORMAT_TIME)
	return r
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required;email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required;min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required;min=6"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type SessionResponse struct {
	ID        uint   `json:"id"`
	DeviceID  string `json:"device_id,omitempty"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent,omitempty"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

func (r SessionResponse) FromEntity(t user.RefreshToken) SessionResponse {
	r.ID = t.ID
	r.DeviceID = t.DeviceID
	r.IPAddress = t.IPAddress
	r.UserAgent = t.UserAgent
	r.CreatedAt = t.CreatedAt.Format(variables.FORMAT_TIME)
	r.ExpiresAt = t.ExpiresAt.Format(variables.FORMAT_TIME)
	return r
}
