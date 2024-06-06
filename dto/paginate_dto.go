package dto

type Paginate struct {
	Page     int  `json:"page"`
	PerPage  int  `json:"per_page"`
	Total    int  `json:"total"`
	NextPage bool `json:"next_page"`
}

type PaginateRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
