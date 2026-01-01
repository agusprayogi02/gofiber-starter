package helper

import (
	"encoding/json"

	"starter-gofiber/dto"

	"github.com/gofiber/fiber/v2"
)

type ResponseWithData struct {
	Code     int             `json:"code"`
	Status   string          `json:"status"`
	Message  string          `json:"message"`
	Paginate *dto.Pagination `json:"paginate,omitempty"`
	Data     any             `json:"data"`
}

type ResponseWithoutData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(p dto.ResponseResult, c *fiber.Ctx) error {
	var response any
	var status string

	if p.StatusCode >= 200 && p.StatusCode <= 299 {
		status = "success"
	} else {
		status = "error"
	}

	if p.Data != nil {
		response = &ResponseWithData{
			Code:     p.StatusCode,
			Status:   status,
			Message:  p.Message,
			Paginate: p.Paginate,
			Data:     p.Data,
		}
	} else {
		response = &ResponseWithoutData{
			Code:    p.StatusCode,
			Status:  status,
			Message: p.Message,
		}
	}

	return c.Status(p.StatusCode).JSON(response)
}

// ToJSON converts interface to JSON bytes
func ToJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
