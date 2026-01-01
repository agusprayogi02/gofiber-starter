package user

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (c CustomClaims) FromToken(j jwt.MapClaims) (CustomClaims, error) {
	// Validate and extract ID
	idVal, ok := j["id"]
	if !ok {
		return c, fmt.Errorf("missing 'id' claim in token")
	}
	idFloat, ok := idVal.(float64)
	if !ok {
		return c, fmt.Errorf("invalid 'id' claim type, expected float64")
	}
	c.ID = uint(idFloat)

	// Validate and extract Email
	emailVal, ok := j["email"]
	if !ok {
		return c, fmt.Errorf("missing 'email' claim in token")
	}
	emailStr, ok := emailVal.(string)
	if !ok {
		return c, fmt.Errorf("invalid 'email' claim type, expected string")
	}
	c.Email = emailStr

	// Validate and extract Role
	roleVal, ok := j["role"]
	if !ok {
		return c, fmt.Errorf("missing 'role' claim in token")
	}
	roleStr, ok := roleVal.(string)
	if !ok {
		return c, fmt.Errorf("invalid 'role' claim type, expected string")
	}
	c.Role = roleStr

	return c, nil
}
