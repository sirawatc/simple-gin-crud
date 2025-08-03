package book

import (
	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/author"
)

type CreateBookRequest struct {
	AuthorID uuid.UUID `json:"authorId" binding:"required" validate:"required"`
	Name     string    `json:"name" binding:"required" validate:"required,min=1,max=255"`
	ISBN     string    `json:"isbn" binding:"required" validate:"required,isbn"`
}

type UpdateBookRequest struct {
	AuthorID uuid.UUID `json:"authorId" binding:"required" validate:"required"`
	Name     string    `json:"name" binding:"required" validate:"required,min=1,max=255"`
	ISBN     string    `json:"isbn" binding:"required" validate:"required,isbn"`
}

type GetBooksByAuthorRequest struct {
	AuthorID uuid.UUID `json:"authorId" uri:"authorId" binding:"required" validate:"required"`
}

type BookResponse struct {
	ID       uuid.UUID              `json:"id"`
	AuthorID uuid.UUID              `json:"authorId"`
	Name     string                 `json:"name"`
	ISBN     string                 `json:"isbn"`
	Author   *author.AuthorResponse `json:"author,omitempty"`
}

type BookListResponse struct {
	Books []BookResponse `json:"books"`
	Total int64          `json:"total"`
}
