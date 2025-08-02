package author

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
	author := &Author{
		PenName:   "Test Author",
		BirthYear: 1990,
	}
	addRow := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO \"authors\" (.+)").WillReturnRows(addRow)
	suite.mock.ExpectCommit()

	err := suite.repo.Create(context.Background(), author)

	suite.NoError(err)
	suite.NotNil(author.ID)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestCreate_Error_DuplicateKey() {
	errMsg := "duplicate key value violates unique constraint"
	author := &Author{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO \"authors\" (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Create(context.Background(), author)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestCreate_Error_ConnectionFailed() {
	errMsg := "connection failed"
	author := &Author{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("INSERT INTO \"authors\" (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Create(context.Background(), author)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_Success() {
	authorID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"}).
		AddRow(uuid.New(), nil, nil, nil, "Test Author", 1990)

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE id = \\$1 (.+)").WillReturnRows(rows)

	author, err := suite.repo.GetByID(context.Background(), authorID)

	suite.NoError(err)
	suite.NotNil(author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_NotFound() {
	authorID := uuid.New()

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE id = \\$1 (.+)").WillReturnError(gorm.ErrRecordNotFound)

	author, err := suite.repo.GetByID(context.Background(), authorID)

	suite.NoError(err)
	suite.Nil(author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByID_DatabaseError() {
	authorID := uuid.New()
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE id = \\$1 (.+)").WillReturnError(errors.New(errMsg))

	author, err := suite.repo.GetByID(context.Background(), authorID)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.Nil(author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByPenName_Success() {
	penName := "Test Author"
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"}).
		AddRow(uuid.New(), nil, nil, nil, "Test Author", 1990)

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE pen_name = \\$1 (.+)").WillReturnRows(rows)

	author, err := suite.repo.GetByPenName(context.Background(), penName)

	suite.NoError(err)
	suite.NotNil(author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByPenName_NotFound() {
	penName := "Non Existent Author"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE pen_name = \\$1 (.+)").WillReturnError(gorm.ErrRecordNotFound)

	author, err := suite.repo.GetByPenName(context.Background(), penName)

	suite.NoError(err)
	suite.Nil(author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetByPenName_DatabaseError() {
	penName := "Test Author"
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" WHERE pen_name = \\$1 (.+)").WillReturnError(errors.New(errMsg))

	author, err := suite.repo.GetByPenName(context.Background(), penName)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.Nil(author)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetAll_Success() {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	dataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"}).
		AddRow(uuid.New(), nil, nil, nil, "Author 1", 1990).
		AddRow(uuid.New(), nil, nil, nil, "Author 2", 1985)

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"authors\" (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" (.+)").WillReturnRows(dataRows)

	result, err := suite.repo.GetAll(context.Background(), pagination)

	suite.NoError(err)
	suite.Equal(2, len(result.Items))
	suite.Equal(pagination.Page, result.Pagination.Page)
	suite.Equal(pagination.PageSize, result.Pagination.PageSize)
	suite.Equal(int64(2), result.Pagination.TotalItems)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestGetAll_EmptyResult() {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	dataRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "pen_name", "birth_year"})

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"authors\" (.+)").WillReturnRows(countRows)
	suite.mock.ExpectQuery("SELECT \\* FROM \"authors\" (.+)").WillReturnRows(dataRows)

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

	suite.mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"authors\" (.+)").WillReturnError(errors.New(errMsg))

	result, err := suite.repo.GetAll(context.Background(), pagination)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.Nil(result)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestUpdate_Success() {
	authorID := uuid.New()
	author := &Author{
		PenName:   "Updated Author",
		BirthYear: 1995,
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"authors\" SET (.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Update(context.Background(), authorID, author)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestUpdate_NotFound() {
	authorID := uuid.New()
	author := &Author{
		PenName:   "Updated Author",
		BirthYear: 1995,
	}

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"authors\" SET (.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repo.Update(context.Background(), authorID, author)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestUpdate_DatabaseError() {
	authorID := uuid.New()
	author := &Author{
		PenName:   "Updated Author",
		BirthYear: 1995,
	}
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"authors\" SET (.+) WHERE id = (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Update(context.Background(), authorID, author)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestDelete_Success() {
	authorID := uuid.New()

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"authors\" SET \"deleted_at\"=(.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Delete(context.Background(), authorID)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestDelete_NotFound() {
	authorID := uuid.New()

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"authors\" SET \"deleted_at\"=(.+) WHERE id = (.+)").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repo.Delete(context.Background(), authorID)

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *RepositoryTestSuite) TestDelete_DatabaseError() {
	authorID := uuid.New()
	errMsg := "connection failed"

	suite.mockTM.On("GetDB").Return(suite.db)

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE \"authors\" SET \"deleted_at\"=(.+) WHERE id = (.+)").WillReturnError(errors.New(errMsg))
	suite.mock.ExpectRollback()

	err := suite.repo.Delete(context.Background(), authorID)

	suite.Error(err)
	suite.Equal(err.Error(), errMsg)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
