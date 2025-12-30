package dto

type ErrorResponse struct {
	Code      int     `json:"code"`
	Order     *string `json:"order"`
	Message   string  `json:"message"`
	Data      any     `json:"data"`
	Timestamp string  `json:"timestamp"`
}

type SuccessResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type ResponseResult struct {
	StatusCode int
	Message    string
	Paginate   *Pagination
	Data       any
}
