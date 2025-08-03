package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirawatc/simple-gin-crud/pkg/logger"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo          IRepository
	authorService IAuthorService
	logger        *logrus.Logger
}

func NewService(repo IRepository, authorService IAuthorService, logger *logrus.Logger) *service {
	return &service{
		repo:          repo,
		authorService: authorService,
		logger:        logger,
	}
}

func (s *service) CreateBook(ctx context.Context, req *CreateBookRequest) (*Book, dto.Code) {
	logPrefix := "[BookService#CreateBook]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	author, code := s.authorService.GetAuthorByID(ctx, req.AuthorID)
	if code != dto.Success {
		logger.Errorf("%s Failed to get author by ID: %v", logPrefix, code)
		return nil, code
	}

	if author == nil {
		logger.Infof("%s Author not found: %v", logPrefix, req.AuthorID)
		return nil, dto.AuthorNotFound
	}

	book, err := s.repo.GetByISBN(ctx, req.ISBN)
	if err != nil {
		logger.Errorf("%s Failed to get book by ISBN: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	if book != nil {
		logger.Infof("%s Book already exists: %v", logPrefix, req.ISBN)
		return nil, dto.BookAlreadyExists
	}

	logger.Infof("%s Creating book: %+v", logPrefix, req)

	book = &Book{
		AuthorID: req.AuthorID,
		Name:     req.Name,
		ISBN:     req.ISBN,
	}

	err = s.repo.Create(ctx, book)
	if err != nil {
		logger.Errorf("%s Failed to create book: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	logger.Infof("%s Book created successfully: %v", logPrefix, book.ID)
	return book, dto.Success
}

func (s *service) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, dto.Code) {
	logPrefix := "[BookService#GetBookByID]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Getting book by ID: %v", logPrefix, id)

	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("%s Failed to get book by ID: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	if book == nil {
		logger.Infof("%s Book not found: %v", logPrefix, id)
		return nil, dto.BookNotFound
	}

	logger.Infof("%s Book retrieved successfully: %v", logPrefix, book.ID)
	return book, dto.Success
}

func (s *service) GetAllBooks(ctx context.Context, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Book], dto.Code) {
	logPrefix := "[BookService#GetAllBooks]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Getting all books: %v", logPrefix, pagination)

	books, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		logger.Errorf("%s Failed to get all books: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	if len(books.Items) == 0 {
		logger.Infof("%s No books found", logPrefix)
		return books, dto.Success
	}

	logger.Infof("%s All books retrieved successfully: %v", logPrefix, books.Pagination)
	return books, dto.Success
}

func (s *service) GetBooksByAuthorID(ctx context.Context, authorID uuid.UUID, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Book], dto.Code) {
	logPrefix := "[BookService#GetBooksByAuthorID]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Getting books by author ID: %v", logPrefix, authorID)

	books, err := s.repo.GetByAuthorID(ctx, authorID, pagination)
	if err != nil {
		logger.Errorf("%s Failed to get books by author ID: %v", logPrefix, err)
		return nil, dto.InternalError
	}

	if len(books.Items) == 0 {
		logger.Infof("%s No books found for author: %v", logPrefix, authorID)
		return books, dto.Success
	}

	logger.Infof("%s Books by author retrieved successfully: %v", logPrefix, books.Pagination)
	return books, dto.Success
}

func (s *service) UpdateBook(ctx context.Context, id uuid.UUID, req *UpdateBookRequest) dto.Code {
	logPrefix := "[BookService#UpdateBook]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("%s Failed to get book by ID: %v", logPrefix, err)
		return dto.InternalError
	}

	if book == nil {
		logger.Infof("%s Book not found: %v", logPrefix, id)
		return dto.BookNotFound
	}

	author, code := s.authorService.GetAuthorByID(ctx, req.AuthorID)
	if code != dto.Success {
		logger.Errorf("%s Failed to get author by ID: %v", logPrefix, code)
		return code
	}

	if author == nil {
		logger.Infof("%s Author not found: %v", logPrefix, req.AuthorID)
		return dto.AuthorNotFound
	}

	logger.Infof("%s Updating book %v: %+v", logPrefix, id, req)

	book = &Book{
		AuthorID: req.AuthorID,
		Name:     req.Name,
		ISBN:     req.ISBN,
	}

	err = s.repo.Update(ctx, id, book)
	if err != nil {
		logger.Errorf("%s Failed to update book: %v", logPrefix, err)
		return dto.InternalError
	}

	logger.Infof("%s Book %v updated successfully", logPrefix, id)
	return dto.Success
}

func (s *service) DeleteBook(ctx context.Context, id uuid.UUID) dto.Code {
	logPrefix := "[BookService#DeleteBook]"
	logger := logger.InjectRequestIDWithLogger(ctx, s.logger)

	logger.Infof("%s Deleting book %v", logPrefix, id)

	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Errorf("%s Failed to delete book: %v", logPrefix, err)
		return dto.InternalError
	}

	logger.Infof("%s Book deleted successfully", logPrefix)
	return dto.Success
}
