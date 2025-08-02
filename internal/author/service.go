package author

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirawatc/simple-gin-crud/pkg/logger"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo   IRepository
	logger *logrus.Logger
}

func NewService(repo IRepository, logger *logrus.Logger) *service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateAuthor(ctx context.Context, req *CreateAuthorRequest) (*Author, dto.Code) {
	logPrefix := "[AuthorService#CreateAuthor]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	author, err := s.repo.GetByPenName(ctx, req.PenName)
	if err != nil {
		logger.Errorf("%s Failed to get author by pen name: %v", logPrefix, err)
		return nil, dto.InternalError
	}
	if author != nil {
		logger.Infof("%s Author already exists: %v", logPrefix, author.ID)
		return nil, dto.AuthorAlreadyExists
	}

	logger.Infof("%s Creating author: %+v", logPrefix, req)

	author = &Author{
		PenName:   req.PenName,
		BirthYear: req.BirthYear,
	}

	err = s.repo.Create(ctx, author)
	if err != nil {
		logger.Errorf("%s Failed to create author: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	logger.Infof("%s Author created successfully: %v", logPrefix, author.ID)
	return author, dto.Success
}

func (s *service) GetAuthorByID(ctx context.Context, id uuid.UUID) (*Author, dto.Code) {
	logPrefix := "[AuthorService#GetAuthorByID]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Getting author by ID: %v", logPrefix, id)

	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("%s Failed to get author by ID: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	if author == nil {
		logger.Infof("%s Author not found: %v", logPrefix, id)
		return nil, dto.Success
	}

	logger.Infof("%s Author retrieved successfully: %v", logPrefix, author.ID)
	return author, dto.Success
}

func (s *service) GetAllAuthors(ctx context.Context, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Author], dto.Code) {
	logPrefix := "[AuthorService#GetAllAuthors]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Getting all authors: %v", logPrefix, pagination)

	authors, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		logger.Errorf("%s Failed to get all authors: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	if len(authors.Items) == 0 {
		logger.Infof("%s No authors found", logPrefix)
		return authors, dto.Success
	}

	logger.Infof("%s All authors retrieved successfully: %v", logPrefix, authors.Pagination)
	return authors, dto.Success
}

func (s *service) UpdateAuthor(ctx context.Context, id uuid.UUID, req *UpdateAuthorRequest) dto.Code {
	logPrefix := "[AuthorService#UpdateAuthor]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("%s Failed to get author by ID: %v", logPrefix, err)
		return dto.InternalError
	}
	if author == nil {
		logger.Infof("%s Author not found: %v", logPrefix, id)
		return dto.AuthorNotFound
	}

	logger.Infof("%s Updating author %v: %+v", logPrefix, id, req)

	author = &Author{
		PenName:   req.PenName,
		BirthYear: req.BirthYear,
	}

	err = s.repo.Update(ctx, id, author)
	if err != nil {
		logger.Errorf("%s Failed to update author: %v", logPrefix, err)
		return dto.InternalError
	}

	logger.Infof("%s Author %v updated successfully", logPrefix, id)
	return dto.Success
}

func (s *service) DeleteAuthor(ctx context.Context, id uuid.UUID) dto.Code {
	logPrefix := "[AuthorService#DeleteAuthor]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Deleting author %v", logPrefix, id)

	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Errorf("%s Failed to delete author: %v", logPrefix, err)
		return dto.InternalError
	}

	logger.Infof("%s Author deleted successfully", logPrefix)
	return dto.Success
}
