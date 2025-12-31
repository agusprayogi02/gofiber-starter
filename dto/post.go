package dto

import (
	"starter-gofiber/entity"
	"starter-gofiber/variables"
)

type PostResponse struct {
	ID        uint          `json:"id"`
	Tweet     string        `json:"tweet"`
	Photo     *string       `json:"photo"`
	UserID    uint          `json:"user_id"`
	User      *UserResponse `json:"user"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

type PostRequest struct {
	Tweet  string `json:"tweet" validate:"required,max=500"`
	Photo  string `json:"photo"`
	UserID uint   `json:"user_id" validate:"required"`
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" validate:"required,min=10"`
}

type PostUpdateRequest struct {
	ID     uint    `json:"id" validate:"required"`
	Tweet  string  `json:"tweet" validate:"required,max=500"`
	Photo  *string `json:"photo"`
	UserID uint    `json:"user_id" validate:"required"`
}

func (r PostRequest) ToEntity() entity.Post {
	return entity.Post{
		Tweet:  r.Tweet,
		Photo:  &r.Photo,
		UserID: r.UserID,
	}
}

func (r PostResponse) FromEntity(p entity.Post) PostResponse {
	r.ID = p.ID
	r.Tweet = p.Tweet
	if p.Photo != nil {
		path := variables.GenerateStatic([]string{variables.POST_PATH, *p.Photo})
		r.Photo = &path
	}
	r.UserID = p.UserID
	if p.User != nil {
		user := UserResponse{}.FromEntity(*p.User)
		r.User = &user
	}
	r.CreatedAt = p.CreatedAt.Format(variables.FORMAT_TIME)
	r.UpdatedAt = p.UpdatedAt.Format(variables.FORMAT_TIME)
	return r
}
