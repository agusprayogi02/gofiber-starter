package dto

type Pagination struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Sort       string `json:"sort,omitempty"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	NextPage   bool   `json:"next_page"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.PerPage == 0 {
		p.PerPage = 10
	}
	return p.PerPage
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "posts.id desc"
	}
	return p.Sort
}
