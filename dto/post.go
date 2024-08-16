package dto

import (
	"starter-gofiber/entity"
	"starter-gofiber/variables"

	"github.com/gofiber/fiber/v2/log"
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
	log.Debug("FromEntity user 1")
	r.ID = p.ID
	r.Tweet = p.Tweet
	r.Photo = p.Photo
	r.UserID = p.UserID
	if p.User != nil {
		log.Debug("FromEntity user 2")
		user := UserResponse{}.FromEntity(*p.User)
		r.User = &user
	}
	log.Debug("FromEntity user 3")
	r.CreatedAt = p.CreatedAt.Format(variables.FORMAT_TIME)
	r.UpdatedAt = p.UpdatedAt.Format(variables.FORMAT_TIME)
	log.Debug("FromEntity user 4")
	return r
}
