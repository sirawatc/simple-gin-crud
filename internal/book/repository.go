package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirawatc/simple-gin-crud/pkg/logger"
	pkgRepo "github.com/sirawatc/simple-gin-crud/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type repository struct {
	transactionManager pkgRepo.ITransactionManager
	logger             *logrus.Logger
}

func NewRepository(transactionManager pkgRepo.ITransactionManager, logger *logrus.Logger) *repository {
	return &repository{
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (r *repository) Create(ctx context.Context, book *Book, tx ...*gorm.DB) error {
	logPrefix := "[BookRepository#Create]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)

	if err := db.Create(book).Error; err != nil {
		logger.Errorf("%s Failed to create book: %v", logPrefix, err)
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) (*Book, error) {
	logPrefix := "[BookRepository#GetByID]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var book Book

	if err := db.Preload("Author").First(&book, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s Book not found: %v", logPrefix, id)
			return nil, nil
		}
		logger.Errorf("%s Failed to get book by ID: %v", logPrefix, err)
		return nil, err
	}

	return &book, nil
}

func (r *repository) GetByISBN(ctx context.Context, isbn string, tx ...*gorm.DB) (*Book, error) {
	logPrefix := "[BookRepository#GetByISBN]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var book Book

	if err := db.Preload("Author").First(&book, "isbn = ?", isbn).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s Book not found: %v", logPrefix, isbn)
			return nil, nil
		}
		logger.Errorf("%s Failed to get book by ISBN: %v", logPrefix, err)
		return nil, err
	}

	return &book, nil
}

func (r *repository) GetByAuthorID(ctx context.Context, authorID uuid.UUID, pagination *dto.PaginationRequest, tx ...*gorm.DB) (*dto.PaginationDataResponse[Book], error) {
	logPrefix := "[BookRepository#GetByAuthorID]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var books []Book
	var total int64

	if err := db.Model(&Book{}).Where("author_id = ?", authorID).Count(&total).Error; err != nil {
		logger.Errorf("%s Failed to count total books for author: %v", logPrefix, err)
		return nil, err
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	err := db.Where("author_id = ?", authorID).Offset(offset).Limit(limit).Find(&books).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s No books found for author: %v", logPrefix, authorID)
			return dto.NewPaginationDataResponse([]Book{}, pagination, total), nil
		}
		logger.Errorf("%s Failed to get paginated books for author: %v", logPrefix, err)
		return nil, err
	}

	return dto.NewPaginationDataResponse(books, pagination, total), nil
}

func (r *repository) GetAll(ctx context.Context, pagination *dto.PaginationRequest, tx ...*gorm.DB) (*dto.PaginationDataResponse[Book], error) {
	logPrefix := "[BookRepository#GetAll]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var books []Book
	var total int64

	if err := db.Model(&Book{}).Count(&total).Error; err != nil {
		logger.Errorf("%s Failed to count total books: %v", logPrefix, err)
		return nil, err
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	err := db.Preload("Author").Offset(offset).Limit(limit).Find(&books).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s No books found", logPrefix)
			return dto.NewPaginationDataResponse([]Book{}, pagination, total), nil
		}
		logger.Errorf("%s Failed to get paginated books: %v", logPrefix, err)
		return nil, err
	}

	return dto.NewPaginationDataResponse(books, pagination, total), nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, book *Book, tx ...*gorm.DB) error {
	logPrefix := "[BookRepository#Update]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)

	if err := db.Model(&Book{}).Where("id = ?", id).Updates(book).Error; err != nil {
		logger.Errorf("%s Failed to update book: %v", logPrefix, err)
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) error {
	logPrefix := "[BookRepository#Delete]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)

	if err := db.Delete(&Book{}, "id = ?", id).Error; err != nil {
		logger.Errorf("%s Failed to delete book: %v", logPrefix, err)
		return err
	}

	return nil
}
