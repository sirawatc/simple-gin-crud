package author

import (
	"github.com/google/uuid"
)

type CreateAuthorRequest struct {
	PenName   string `json:"penName" binding:"required" validate:"required,min=1,max=255"`
	BirthYear int    `json:"birthYear" binding:"required" validate:"required,min=1800,max=2600"`
}

type UpdateAuthorRequest struct {
	PenName   string `json:"penName" binding:"required" validate:"required,min=1,max=255"`
	BirthYear int    `json:"birthYear" binding:"required" validate:"required,min=1800,max=2600"`
}

type AuthorResponse struct {
	ID        uuid.UUID `json:"id"`
	PenName   string    `json:"penName"`
	BirthYear int       `json:"birthYear"`
}

type AuthorListResponse struct {
	Authors []AuthorResponse `json:"authors"`
	Total   int64            `json:"total"`
}
