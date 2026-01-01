package post

// Repository defines the interface for post repository operations
type Repository interface {
	Create(post *Post) error
	FindByID(id uint) (*Post, error)
	FindAll(limit, offset int) ([]Post, int64, error)
	Update(post *Post) error
	Delete(id uint) error
	FindByUserID(userID uint, limit, offset int) ([]Post, int64, error)
}
