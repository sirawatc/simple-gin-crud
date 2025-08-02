package author

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirawatc/simple-gin-crud/pkg/logger"
	repoPkg "github.com/sirawatc/simple-gin-crud/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type repository struct {
	transactionManager repoPkg.ITransactionManager
	logger             *logrus.Logger
}

func NewRepository(transactionManager repoPkg.ITransactionManager, logger *logrus.Logger) *repository {
	return &repository{
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (r *repository) Create(ctx context.Context, author *Author, tx ...*gorm.DB) error {
	logPrefix := "[AuthorRepository#Create]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)

	if err := db.Create(author).Error; err != nil {
		logger.Errorf("%s Failed to create author: %v", logPrefix, err)
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) (*Author, error) {
	logPrefix := "[AuthorRepository#GetByID]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var author Author

	if err := db.First(&author, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s Author not found: %v", logPrefix, id)
			return nil, nil
		}
		logger.Errorf("%s Failed to get author by ID: %v", logPrefix, err)
		return nil, err
	}

	return &author, nil
}

func (r *repository) GetByPenName(ctx context.Context, penName string, tx ...*gorm.DB) (*Author, error) {
	logPrefix := "[AuthorRepository#GetByPenName]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var author Author

	if err := db.First(&author, "pen_name = ?", penName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s Author not found: %v", logPrefix, penName)
			return nil, nil
		}
		logger.Errorf("%s Failed to get author by pen name: %v", logPrefix, err)
		return nil, err
	}

	return &author, nil
}

func (r *repository) GetAll(ctx context.Context, pagination *dto.PaginationRequest, tx ...*gorm.DB) (*dto.PaginationDataResponse[Author], error) {
	logPrefix := "[AuthorRepository#GetAll]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)
	var authors []Author
	var total int64

	if err := db.Model(&Author{}).Count(&total).Error; err != nil {
		logger.Errorf("%s Failed to count total authors: %v", logPrefix, err)
		return nil, err
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	err := db.Offset(offset).Limit(limit).Find(&authors).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("%s No authors found", logPrefix)
			return dto.NewPaginationDataResponse([]Author{}, pagination, total), nil
		}
		logger.Errorf("%s Failed to get paginated authors: %v", logPrefix, err)
		return nil, err
	}

	return dto.NewPaginationDataResponse(authors, pagination, total), nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, author *Author, tx ...*gorm.DB) error {
	logPrefix := "[AuthorRepository#Update]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)

	if err := db.Model(&Author{}).Where("id = ?", id).Updates(author).Error; err != nil {
		logger.Errorf("%s Failed to update author: %v", logPrefix, err)
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) error {
	logPrefix := "[AuthorRepository#Delete]"
	logger := logger.InjectRequestIDWithLogger(ctx, r.logger)

	db := r.transactionManager.GetDB(tx...)

	if err := db.Delete(&Author{}, "id = ?", id).Error; err != nil {
		logger.Errorf("%s Failed to delete author: %v", logPrefix, err)
		return err
	}

	return nil
}
