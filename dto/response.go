package dto

type ErrorResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type SuccessResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type ResponseParams struct {
	StatusCode int
	Message    string
	Paginate   *Paginate
	Data       any
}
