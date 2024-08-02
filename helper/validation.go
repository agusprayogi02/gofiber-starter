package helper

import (
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

type IError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
