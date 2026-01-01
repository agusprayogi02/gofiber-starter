package post

// Service defines the interface for post service operations
type Service interface {
	Create(req *PostRequest, userID uint) (*PostResponse, error)
	FindByID(id uint) (*PostResponse, error)
	FindAll(page, limit int) ([]PostResponse, *PaginationMeta, error)
	Update(id uint, req *PostUpdateRequest, userID uint) (*PostResponse, error)
	Delete(id uint, userID uint) error
	FindByUserID(userID uint, page, limit int) ([]PostResponse, *PaginationMeta, error)
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"total_pages"`
	HasNextPage bool  `json:"has_next_page"`
	HasPrevPage bool  `json:"has_prev_page"`
}
