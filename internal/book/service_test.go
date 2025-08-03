package book

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/author"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	"github.com/sirawatc/simple-gin-crud/internal/shared/models"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, book *Book, tx ...*gorm.DB) error {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, book, tx)
	} else {
		args = m.Called(ctx, book)
	}
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) (*Book, error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, id, tx)
	} else {
		args = m.Called(ctx, id)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Book), args.Error(1)
}

func (m *MockRepository) GetByISBN(ctx context.Context, isbn string, tx ...*gorm.DB) (*Book, error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, isbn, tx)
	} else {
		args = m.Called(ctx, isbn)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Book), args.Error(1)
}

func (m *MockRepository) GetAll(ctx context.Context, pagination *pkgDto.PaginationRequest, tx ...*gorm.DB) (*pkgDto.PaginationDataResponse[Book], error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, pagination, tx)
	} else {
		args = m.Called(ctx, pagination)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pkgDto.PaginationDataResponse[Book]), args.Error(1)
}

func (m *MockRepository) GetByAuthorID(ctx context.Context, authorID uuid.UUID, pagination *pkgDto.PaginationRequest, tx ...*gorm.DB) (*pkgDto.PaginationDataResponse[Book], error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, authorID, pagination, tx)
	} else {
		args = m.Called(ctx, authorID, pagination)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pkgDto.PaginationDataResponse[Book]), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id uuid.UUID, book *Book, tx ...*gorm.DB) error {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, id, book, tx)
	} else {
		args = m.Called(ctx, id, book)
	}
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) error {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, id, tx)
	} else {
		args = m.Called(ctx, id)
	}
	return args.Error(0)
}

type MockAuthorService struct {
	mock.Mock
}

func (m *MockAuthorService) GetAuthorByID(ctx context.Context, id uuid.UUID) (*author.Author, dto.Code) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*author.Author), args.Get(1).(dto.Code)
}

type ServiceTestSuite struct {
	suite.Suite
	service           IService
	mockRepo          *MockRepository
	mockAuthorService *MockAuthorService
	ctx               context.Context
}

func (suite *ServiceTestSuite) SetupTest() {
	mockRepo := new(MockRepository)
	mockAuthorService := new(MockAuthorService)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	service := NewService(mockRepo, mockAuthorService, logger)

	suite.service = service
	suite.mockRepo = mockRepo
	suite.mockAuthorService = mockAuthorService
	suite.ctx = context.Background()
}

func (suite *ServiceTestSuite) TestNewService() {
	mockRepo := new(MockRepository)
	mockAuthorService := new(MockAuthorService)
	logger := logrus.New()
	service := NewService(mockRepo, mockAuthorService, logger)

	suite.NotNil(service)

	// Test that the service implements the interface
	var _ IService = service
	suite.Implements((*IService)(nil), service)
}

func (suite *ServiceTestSuite) TestCreateBook_Success() {
	authorID := uuid.New()
	req := &CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	expectedAuthor := &author.Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return(expectedAuthor, dto.Success)
	suite.mockRepo.On("GetByISBN", suite.ctx, req.ISBN).Return((*Book)(nil), nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*book.Book")).Return(nil)

	book, code := suite.service.CreateBook(suite.ctx, req)

	suite.Equal(dto.Success, code)
	suite.NotNil(book)
	suite.Equal(req.Name, book.Name)
	suite.Equal(req.ISBN, book.ISBN)
	suite.Equal(req.AuthorID, book.AuthorID)
	suite.mockAuthorService.AssertExpectations(suite.T())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateBook_AuthorNotFound() {
	authorID := uuid.New()
	req := &CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return((*author.Author)(nil), dto.AuthorNotFound)

	book, code := suite.service.CreateBook(suite.ctx, req)

	suite.Equal(dto.AuthorNotFound, code)
	suite.Nil(book)
	suite.mockAuthorService.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateBook_GetAuthorByIDError() {
	authorID := uuid.New()
	req := &CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return((*author.Author)(nil), dto.InternalError)

	book, code := suite.service.CreateBook(suite.ctx, req)

	suite.Equal(dto.InternalError, code)
	suite.Nil(book)
	suite.mockAuthorService.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateBook_BookAlreadyExists() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	expectedAuthor := &author.Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	existingBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Existing Book",
		ISBN:      "978-0-7475-3269-9",
	}

	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return(expectedAuthor, dto.Success)
	suite.mockRepo.On("GetByISBN", suite.ctx, req.ISBN).Return(existingBook, nil)

	book, code := suite.service.CreateBook(suite.ctx, req)

	suite.Equal(dto.BookAlreadyExists, code)
	suite.Nil(book)
	suite.mockAuthorService.AssertExpectations(suite.T())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateBook_GetByISBNError() {
	authorID := uuid.New()
	req := &CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	expectedAuthor := &author.Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return(expectedAuthor, dto.Success)
	suite.mockRepo.On("GetByISBN", suite.ctx, req.ISBN).Return((*Book)(nil), errors.New("database error"))

	book, code := suite.service.CreateBook(suite.ctx, req)

	suite.Equal(dto.InternalError, code)
	suite.Nil(book)
	suite.mockAuthorService.AssertExpectations(suite.T())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateBook_CreateError() {
	authorID := uuid.New()
	req := &CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	expectedAuthor := &author.Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return(expectedAuthor, dto.Success)
	suite.mockRepo.On("GetByISBN", suite.ctx, req.ISBN).Return((*Book)(nil), nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*book.Book")).Return(errors.New("database error"))

	book, code := suite.service.CreateBook(suite.ctx, req)

	suite.Equal(dto.InternalError, code)
	suite.Nil(book)
	suite.mockAuthorService.AssertExpectations(suite.T())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBookByID_Success() {
	bookID := uuid.New()
	authorID := uuid.New()
	expectedBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Test Book",
		ISBN:      "1234567890123",
		Author: &author.Author{
			BaseModel: models.BaseModel{ID: authorID},
			PenName:   "Test Author",
			BirthYear: 1990,
		},
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return(expectedBook, nil)

	book, code := suite.service.GetBookByID(suite.ctx, bookID)

	suite.Equal(dto.Success, code)
	suite.NotNil(book)
	suite.Equal(expectedBook.ID, book.ID)
	suite.Equal(expectedBook.Name, book.Name)
	suite.NotNil(book.Author)
	suite.Equal(expectedBook.Author.ID, book.Author.ID)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBookByID_Success_WithoutAuthor() {
	bookID := uuid.New()
	authorID := uuid.New()
	expectedBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Test Book",
		ISBN:      "1234567890123",
		Author:    nil,
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return(expectedBook, nil)

	book, code := suite.service.GetBookByID(suite.ctx, bookID)

	suite.Equal(dto.Success, code)
	suite.NotNil(book)
	suite.Equal(expectedBook.ID, book.ID)
	suite.Equal(expectedBook.Name, book.Name)
	suite.Nil(book.Author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBookByID_NotFound() {
	bookID := uuid.New()

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return((*Book)(nil), nil)

	book, code := suite.service.GetBookByID(suite.ctx, bookID)

	suite.Equal(dto.BookNotFound, code)
	suite.Nil(book)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBookByID_GetByIDError() {
	bookID := uuid.New()

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return((*Book)(nil), errors.New("database error"))

	book, code := suite.service.GetBookByID(suite.ctx, bookID)

	suite.Equal(dto.InternalError, code)
	suite.Nil(book)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllBooks_Success() {
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 5}
	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				Name:      "Book 1",
				ISBN:      "1234567890123",
				Author: &author.Author{
					BaseModel: models.BaseModel{ID: uuid.New()},
					PenName:   "Author 1",
					BirthYear: 1990,
				},
			},
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				Name:      "Book 2",
				ISBN:      "1234567890124",
				Author: &author.Author{
					BaseModel: models.BaseModel{ID: uuid.New()},
					PenName:   "Author 2",
					BirthYear: 1985,
				},
			},
		},
		Pagination: pkgDto.PaginationResponse{
			Page:       1,
			PageSize:   5,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	suite.mockRepo.On("GetAll", suite.ctx, pagination).Return(expectedBooks, nil)

	books, code := suite.service.GetAllBooks(suite.ctx, pagination)

	suite.Equal(dto.Success, code)
	suite.NotNil(books)
	suite.Equal(expectedBooks.Items, books.Items)
	suite.Equal(expectedBooks.Pagination, books.Pagination)
	suite.NotNil(books.Items[0].Author)
	suite.NotNil(books.Items[1].Author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllBooks_EmptyResult() {
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 10}
	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{},
		Pagination: pkgDto.PaginationResponse{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	suite.mockRepo.On("GetAll", suite.ctx, pagination).Return(expectedBooks, nil)

	books, code := suite.service.GetAllBooks(suite.ctx, pagination)

	suite.Equal(dto.Success, code)
	suite.NotNil(books)
	suite.Empty(books.Items)
	suite.Equal(expectedBooks.Pagination, books.Pagination)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllBooks_GetAllError() {
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 10}

	suite.mockRepo.On("GetAll", suite.ctx, pagination).Return((*pkgDto.PaginationDataResponse[Book])(nil), errors.New("database error"))

	books, code := suite.service.GetAllBooks(suite.ctx, pagination)

	suite.Equal(dto.InternalError, code)
	suite.Nil(books)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBooksByAuthorID_Success() {
	authorID := uuid.New()
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 5}
	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				AuthorID:  authorID,
				Name:      "Book 1",
				ISBN:      "1234567890123",
			},
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				AuthorID:  authorID,
				Name:      "Book 2",
				ISBN:      "1234567890124",
			},
		},
		Pagination: pkgDto.PaginationResponse{
			Page:       1,
			PageSize:   5,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	suite.mockRepo.On("GetByAuthorID", suite.ctx, authorID, pagination).Return(expectedBooks, nil)

	books, code := suite.service.GetBooksByAuthorID(suite.ctx, authorID, pagination)

	suite.Equal(dto.Success, code)
	suite.NotNil(books)
	suite.Equal(expectedBooks.Items, books.Items)
	suite.Equal(expectedBooks.Pagination, books.Pagination)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBooksByAuthorID_EmptyResult() {
	authorID := uuid.New()
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 10}
	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{},
		Pagination: pkgDto.PaginationResponse{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	suite.mockRepo.On("GetByAuthorID", suite.ctx, authorID, pagination).Return(expectedBooks, nil)

	books, code := suite.service.GetBooksByAuthorID(suite.ctx, authorID, pagination)

	suite.Equal(dto.Success, code)
	suite.NotNil(books)
	suite.Empty(books.Items)
	suite.Equal(expectedBooks.Pagination, books.Pagination)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetBooksByAuthorID_GetByAuthorIDError() {
	authorID := uuid.New()
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 10}

	suite.mockRepo.On("GetByAuthorID", suite.ctx, authorID, pagination).Return((*pkgDto.PaginationDataResponse[Book])(nil), errors.New("database error"))

	books, code := suite.service.GetBooksByAuthorID(suite.ctx, authorID, pagination)

	suite.Equal(dto.InternalError, code)
	suite.Nil(books)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateBook_Success() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	existingBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Original Book",
		ISBN:      "1234567890123",
	}

	expectedAuthor := &author.Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return(existingBook, nil)
	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return(expectedAuthor, dto.Success)
	suite.mockRepo.On("Update", suite.ctx, bookID, mock.AnythingOfType("*book.Book")).Return(nil)

	code := suite.service.UpdateBook(suite.ctx, bookID, req)

	suite.Equal(dto.Success, code)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockAuthorService.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateBook_BookNotFound() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return((*Book)(nil), nil)

	code := suite.service.UpdateBook(suite.ctx, bookID, req)

	suite.Equal(dto.BookNotFound, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateBook_GetByIDError() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return((*Book)(nil), errors.New("database error"))

	code := suite.service.UpdateBook(suite.ctx, bookID, req)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateBook_AuthorNotFound() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	existingBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Original Book",
		ISBN:      "1234567890123",
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return(existingBook, nil)
	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return((*author.Author)(nil), dto.AuthorNotFound)

	code := suite.service.UpdateBook(suite.ctx, bookID, req)

	suite.Equal(dto.AuthorNotFound, code)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockAuthorService.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateBook_GetAuthorByIDError() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	existingBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Original Book",
		ISBN:      "1234567890123",
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return(existingBook, nil)
	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return((*author.Author)(nil), dto.InternalError)

	code := suite.service.UpdateBook(suite.ctx, bookID, req)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockAuthorService.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateBook_UpdateError() {
	bookID := uuid.New()
	authorID := uuid.New()
	req := &UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	existingBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Original Book",
		ISBN:      "1234567890123",
	}

	expectedAuthor := &author.Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByID", suite.ctx, bookID).Return(existingBook, nil)
	suite.mockAuthorService.On("GetAuthorByID", suite.ctx, authorID).Return(expectedAuthor, dto.Success)
	suite.mockRepo.On("Update", suite.ctx, bookID, mock.AnythingOfType("*book.Book")).Return(errors.New("database error"))

	code := suite.service.UpdateBook(suite.ctx, bookID, req)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockAuthorService.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestDeleteBook_Success() {
	bookID := uuid.New()

	suite.mockRepo.On("Delete", suite.ctx, bookID).Return(nil)

	code := suite.service.DeleteBook(suite.ctx, bookID)

	suite.Equal(dto.Success, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestDeleteBook_DeleteError() {
	bookID := uuid.New()

	suite.mockRepo.On("Delete", suite.ctx, bookID).Return(errors.New("database error"))

	code := suite.service.DeleteBook(suite.ctx, bookID)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
