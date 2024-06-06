package helper

import "starter-gofiber/dto"

type ResponseWithData struct {
	Code     int           `json:"code"`
	Status   string        `json:"status"`
	Message  string        `json:"message"`
	Paginate *dto.Paginate `json:"paginate,omitempty"`
	Data     any           `json:"data"`
}

type ResponseWithoutData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(p dto.ResponseParams) any {
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

	return response
}
