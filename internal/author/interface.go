package author

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(ctx context.Context, author *Author, tx ...*gorm.DB) error
	GetByID(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) (*Author, error)
	GetByPenName(ctx context.Context, penName string, tx ...*gorm.DB) (*Author, error)
	GetAll(ctx context.Context, pagination *pkgDto.PaginationRequest, tx ...*gorm.DB) (*pkgDto.PaginationDataResponse[Author], error)
	Update(ctx context.Context, id uuid.UUID, author *Author, tx ...*gorm.DB) error
	Delete(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) error
}

type IService interface {
	CreateAuthor(ctx context.Context, req *CreateAuthorRequest) (*Author, dto.Code)
	GetAuthorByID(ctx context.Context, id uuid.UUID) (*Author, dto.Code)
	GetAllAuthors(ctx context.Context, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Author], dto.Code)
	UpdateAuthor(ctx context.Context, id uuid.UUID, req *UpdateAuthorRequest) dto.Code
	DeleteAuthor(ctx context.Context, id uuid.UUID) dto.Code
}
