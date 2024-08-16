package dto

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (c CustomClaims) FromToken(j jwt.MapClaims) CustomClaims {
	c.ID = uint(j["id"].(float64))
	c.Email = j["email"].(string)
	c.Role = j["role"].(string)
	return c
}
