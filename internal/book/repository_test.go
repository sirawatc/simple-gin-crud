package book

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Transaction(fn func(tx *gorm.DB) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

func (m *MockTransactionManager) GetDB(tx ...*gorm.DB) *gorm.DB {
	args := m.Called()
	if db, ok := args.Get(0).(*gorm.DB); ok {
		return db
	}
	return nil
}

type RepositoryTestSuite struct {
	suite.Suite
	repo   IRepository
	db     *gorm.DB
	mockTM *MockTransactionManager
	mock   sqlmock.Sqlmock
}

func (suite *RepositoryTestSuite) SetupTest() {
	logger := logrus.New()
	mockTM := &MockTransactionManager{}
	db, mock := suite.mockDB()
	repo := NewRepository(mockTM, logger)
	suite.repo = repo
	suite.db = db
	suite.mock = mock
	suite.mockTM = mockTM
}

func (suite *RepositoryTestSuite) mockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	suite.NoError(err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	suite.NoError(err)

	return gormDB, mock
}

func (suite *RepositoryTestSuite) TestNewRepository() {
	logger := logrus.New()
	mockTM := &MockTransactionManager{}
	repo := NewRepository(mockTM, logger)

	suite.NotNil(repo)
	suite.IsType(&repository{}, repo)

	// Test that the repository implements the interface
	var _ IRepository = repo
	suite.Implements((*IRepository)(nil), repo)
}

func (suite *RepositoryTestSuite) TestCreate_Success() {
	book := &Book{
		AuthorID: uuid.New(),
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}
	addRow := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO \"books\" (.+)").WillReturnRows(addRow)
	suite.mock.ExpectCommit()

	err := suite.repo.Create(context.Background(), book)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestCreate_Error_DuplicateKey() {
	errMsg := "duplicate key value violates unique constraint"
	book := &Book{
		AuthorID: uuid.New(),
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO \"books\" (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Create(context.Background(), book)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestCreate_Error_ConnectionFailed() {
	errMsg := "connection failed"
	authorID := uuid.New()
	book := &Book{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO \"books\" (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Create(context.Background(), book)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_Success() {
	bookID := uuid.New()
	authorID := uuid.New()
	bookDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(bookID, nil, nil, nil, authorID, "Test Book", "978-0-7475-3269-9")
	authorDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"}).
		AddRow(authorID, nil, nil, nil, "Author 1", 1990)

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE id = (.+)").WillReturnRows(bookDataRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE \"authors\".\"id\" = (.+)").WillReturnRows(authorDataRows)

	book, err := suite.repo.GetByID(context.Background(), bookID)

	suite.NoError(err)
	suite.NotNil(book)
	suite.NotNil(book.Author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_Success_WithAuthor() {
	bookID := uuid.New()
	bookDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(bookID, nil, nil, nil, uuid.New(), "Test Book", "978-0-7475-3269-9")
	authorDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"})

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE id = (.+)").WillReturnRows(bookDataRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE \"authors\".\"id\" = (.+)").WillReturnRows(authorDataRows)

	book, err := suite.repo.GetByID(context.Background(), bookID)

	suite.NoError(err)
	suite.NotNil(book)
	suite.Nil(book.Author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_NotFound() {
	bookID := uuid.New()

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE id = (.+)").WillReturnError(gorm.ErrRecordNotFound)

	book, err := suite.repo.GetByID(context.Background(), bookID)

	suite.NoError(err)
	suite.Nil(book)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_DatabaseError() {
	bookID := uuid.New()
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE id = (.+)").WillReturnError(errors.New(errMsg))

	book, err := suite.repo.GetByID(context.Background(), bookID)

	suite.Error(err)
	suite.Nil(book)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByISBN_Success() {
	isbn := "978-0-7475-3269-9"
	bookID := uuid.New()
	authorID := uuid.New()
	bookDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(bookID, nil, nil, nil, authorID, "Test Book", isbn)
	authorDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"}).
		AddRow(authorID, nil, nil, nil, "Author 1", 1990)

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE isbn = (.+)").WillReturnRows(bookDataRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE \"authors\".\"id\" = (.+)").WillReturnRows(authorDataRows)

	book, err := suite.repo.GetByISBN(context.Background(), isbn)

	suite.NoError(err)
	suite.NotNil(book)
	suite.NotNil(book.Author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByISBN_Success_WithAuthor() {
	isbn := "978-0-7475-3269-9"
	bookID := uuid.New()
	bookDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(bookID, nil, nil, nil, uuid.New(), "Test Book", isbn)
	authorDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"})

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE isbn = (.+)").WillReturnRows(bookDataRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE \"authors\".\"id\" = (.+)").WillReturnRows(authorDataRows)

	book, err := suite.repo.GetByISBN(context.Background(), isbn)

	suite.NoError(err)
	suite.NotNil(book)
	suite.Nil(book.Author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByISBN_NotFound() {
	isbn := "978-0-7475-3269-9"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE isbn = (.+)").WillReturnError(gorm.ErrRecordNotFound)

	book, err := suite.repo.GetByISBN(context.Background(), isbn)

	suite.NoError(err)
	suite.Nil(book)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByISBN_DatabaseError() {
	isbn := "978-0-7475-3269-9"
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE isbn = (.+)").WillReturnError(errors.New(errMsg))

	book, err := suite.repo.GetByISBN(context.Background(), isbn)

	suite.Error(err)
	suite.Nil(book)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetAll_Success() {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	authorID := uuid.New()
	authorID2 := uuid.New()
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	bookDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(uuid.New(), nil, nil, nil, authorID, "Book 1", "978-0-7475-3269-9").
		AddRow(uuid.New(), nil, nil, nil, authorID2, "Book 2", "978-0-7475-3269-8")
	authorDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"}).
		AddRow(authorID, nil, nil, nil, "Author 1", 1990).
		AddRow(authorID2, nil, nil, nil, "Author 2", 1991)

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" (.+)").WillReturnRows(bookDataRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE \"authors\".\"id\" IN (.+)").WillReturnRows(authorDataRows)

	result, err := suite.repo.GetAll(context.Background(), pagination)

	suite.NoError(err)
	suite.Equal(2, len(result.Items))
	suite.NotNil(result.Items[0].Author)
	suite.NotNil(result.Items[1].Author)
	suite.Equal(pagination.Page, result.Pagination.Page)
	suite.Equal(pagination.PageSize, result.Pagination.PageSize)
	suite.Equal(int64(2), result.Pagination.TotalItems)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetAll_Success_WithAuthor() {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	bookDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(uuid.New(), nil, nil, nil, uuid.New(), "Book 1", "978-0-7475-3269-9")
	authorDataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"})

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" (.+)").WillReturnRows(bookDataRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE \"authors\".\"id\" = (.+)").WillReturnRows(authorDataRows)

	result, err := suite.repo.GetAll(context.Background(), pagination)

	suite.NoError(err)
	suite.Equal(1, len(result.Items))
	suite.Nil(result.Items[0].Author)
	suite.Equal(pagination.Page, result.Pagination.Page)
	suite.Equal(pagination.PageSize, result.Pagination.PageSize)
	suite.Equal(int64(1), result.Pagination.TotalItems)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetAll_EmptyResult() {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	dataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"})

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" (.+)").WillReturnRows(dataRows)

	result, err := suite.repo.GetAll(context.Background(), pagination)

	suite.NoError(err)
	suite.Empty(result.Items)
	suite.Equal(pagination.Page, result.Pagination.Page)
	suite.Equal(pagination.PageSize, result.Pagination.PageSize)
	suite.Equal(int64(0), result.Pagination.TotalItems)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetAll_DatabaseError() {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" (.+)").WillReturnError(errors.New(errMsg))

	result, err := suite.repo.GetAll(context.Background(), pagination)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByAuthorID_Success() {
	authorID := uuid.New()
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	dataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"}).
		AddRow(uuid.New(), nil, nil, nil, authorID, "Book 1", "978-0-7475-3269-9").
		AddRow(uuid.New(), nil, nil, nil, authorID, "Book 2", "978-0-7475-3269-8")

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" WHERE author_id = (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE author_id = (.+)").WillReturnRows(dataRows)

	result, err := suite.repo.GetByAuthorID(context.Background(), authorID, pagination)

	suite.NoError(err)
	suite.Equal(2, len(result.Items))
	suite.Equal(pagination.Page, result.Pagination.Page)
	suite.Equal(pagination.PageSize, result.Pagination.PageSize)
	suite.Equal(int64(2), result.Pagination.TotalItems)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByAuthorID_EmptyResult() {
	authorID := uuid.New()
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	dataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "author_id", "name", "isbn"})

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" WHERE author_id = (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"books\" WHERE author_id = (.+)").WillReturnRows(dataRows)

	result, err := suite.repo.GetByAuthorID(context.Background(), authorID, pagination)

	suite.NoError(err)
	suite.Empty(result.Items)
	suite.Equal(pagination.Page, result.Pagination.Page)
	suite.Equal(pagination.PageSize, result.Pagination.PageSize)
	suite.Equal(int64(0), result.Pagination.TotalItems)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByAuthorID_DatabaseError() {
	authorID := uuid.New()
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"books\" WHERE author_id = (.+)").WillReturnError(errors.New(errMsg))

	result, err := suite.repo.GetByAuthorID(context.Background(), authorID, pagination)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestUpdate_Success() {
	bookID := uuid.New()
	authorID := uuid.New()
	book := &Book{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"books\" SET (.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Update(context.Background(), bookID, book)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestUpdate_NotFound() {
	bookID := uuid.New()
	authorID := uuid.New()
	book := &Book{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"books\" SET (.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repo.Update(context.Background(), bookID, book)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestUpdate_DatabaseError() {
	bookID := uuid.New()
	authorID := uuid.New()
	book := &Book{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"books\" SET (.+) WHERE id = (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Update(context.Background(), bookID, book)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestDelete_Success() {
	bookID := uuid.New()

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"books\" SET \"deleted_at\"=(.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Delete(context.Background(), bookID)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestDelete_NotFound() {
	bookID := uuid.New()

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"books\" SET \"deleted_at\"=(.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repo.Delete(context.Background(), bookID)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestDelete_DatabaseError() {
	bookID := uuid.New()
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"books\" SET \"deleted_at\"=(.+) WHERE id = (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Delete(context.Background(), bookID)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
