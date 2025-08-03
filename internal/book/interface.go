package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/author"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"gorm.io/gorm"
)

type IAuthorService interface {
	GetAuthorByID(ctx context.Context, id uuid.UUID) (*author.Author, dto.Code)
}

type IRepository interface {
	Create(ctx context.Context, book *Book, tx ...*gorm.DB) error
	GetByID(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) (*Book, error)
	GetByISBN(ctx context.Context, isbn string, tx ...*gorm.DB) (*Book, error)
	GetAll(ctx context.Context, pagination *pkgDto.PaginationRequest, tx ...*gorm.DB) (*pkgDto.PaginationDataResponse[Book], error)
	Update(ctx context.Context, id uuid.UUID, book *Book, tx ...*gorm.DB) error
	Delete(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) error
	GetByAuthorID(ctx context.Context, authorID uuid.UUID, pagination *pkgDto.PaginationRequest, tx ...*gorm.DB) (*pkgDto.PaginationDataResponse[Book], error)
}

type IService interface {
	CreateBook(ctx context.Context, req *CreateBookRequest) (*Book, dto.Code)
	GetBookByID(ctx context.Context, id uuid.UUID) (*Book, dto.Code)
	GetBooksByAuthorID(ctx context.Context, authorID uuid.UUID, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Book], dto.Code)
	GetAllBooks(ctx context.Context, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Book], dto.Code)
	UpdateBook(ctx context.Context, id uuid.UUID, req *UpdateBookRequest) dto.Code
	DeleteBook(ctx context.Context, id uuid.UUID) dto.Code
}
