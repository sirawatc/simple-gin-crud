package author

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
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

func (m *MockRepository) Create(ctx context.Context, author *Author, tx ...*gorm.DB) error {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, author, tx)
	} else {
		args = m.Called(ctx, author)
	}
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID, tx ...*gorm.DB) (*Author, error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, id, tx)
	} else {
		args = m.Called(ctx, id)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Author), args.Error(1)
}

func (m *MockRepository) GetByPenName(ctx context.Context, penName string, tx ...*gorm.DB) (*Author, error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, penName, tx)
	} else {
		args = m.Called(ctx, penName)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Author), args.Error(1)
}

func (m *MockRepository) GetAll(ctx context.Context, pagination *pkgDto.PaginationRequest, tx ...*gorm.DB) (*pkgDto.PaginationDataResponse[Author], error) {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, pagination, tx)
	} else {
		args = m.Called(ctx, pagination)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pkgDto.PaginationDataResponse[Author]), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id uuid.UUID, author *Author, tx ...*gorm.DB) error {
	var args mock.Arguments
	if len(tx) > 0 {
		args = m.Called(ctx, id, author, tx)
	} else {
		args = m.Called(ctx, id, author)
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

type ServiceTestSuite struct {
	suite.Suite
	service  *service
	mockRepo *MockRepository
	ctx      context.Context
}

func (suite *ServiceTestSuite) SetupTest() {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	service := NewService(mockRepo, logger)

	suite.service = service
	suite.mockRepo = mockRepo
	suite.ctx = context.Background()
}

func (suite *ServiceTestSuite) TestNewService() {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	suite.NotNil(service)

	// Test that the service implements the interface
	var _ IService = service
	suite.Implements((*IService)(nil), service)
}

func (suite *ServiceTestSuite) TestCreateAuthor_Success() {
	req := &CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByPenName", suite.ctx, req.PenName).Return((*Author)(nil), nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*author.Author")).Return(nil)

	author, code := suite.service.CreateAuthor(suite.ctx, req)

	suite.Equal(dto.Success, code)
	suite.NotNil(author)
	suite.Equal(req.PenName, author.PenName)
	suite.Equal(req.BirthYear, author.BirthYear)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateAuthor_AuthorAlreadyExists() {
	authorID := uuid.New()
	req := &CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	existingAuthor := &Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByPenName", suite.ctx, req.PenName).Return(existingAuthor, nil)

	author, code := suite.service.CreateAuthor(suite.ctx, req)

	suite.Equal(dto.AuthorAlreadyExists, code)
	suite.Nil(author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateAuthor_GetByPenNameError() {
	req := &CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByPenName", suite.ctx, req.PenName).Return((*Author)(nil), errors.New("database error"))

	author, code := suite.service.CreateAuthor(suite.ctx, req)

	suite.Equal(dto.InternalError, code)
	suite.Nil(author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreateAuthor_CreateError() {
	req := &CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByPenName", suite.ctx, req.PenName).Return((*Author)(nil), nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*author.Author")).Return(errors.New("database error"))

	author, code := suite.service.CreateAuthor(suite.ctx, req)

	suite.Equal(dto.InternalError, code)
	suite.Nil(author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAuthorByID_Success() {
	authorID := uuid.New()
	expectedAuthor := &Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return(expectedAuthor, nil)

	author, code := suite.service.GetAuthorByID(suite.ctx, authorID)

	suite.Equal(dto.Success, code)
	suite.NotNil(author)
	suite.Equal(expectedAuthor.ID, author.ID)
	suite.Equal(expectedAuthor.PenName, author.PenName)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAuthorByID_NotFound() {
	authorID := uuid.New()

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return((*Author)(nil), nil)

	author, code := suite.service.GetAuthorByID(suite.ctx, authorID)

	suite.Equal(dto.Success, code)
	suite.Nil(author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAuthorByID_GetByIDError() {
	authorID := uuid.New()

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return((*Author)(nil), errors.New("database error"))

	author, code := suite.service.GetAuthorByID(suite.ctx, authorID)

	suite.Equal(dto.InternalError, code)
	suite.Nil(author)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllAuthors_Success() {
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 5}
	expectedAuthors := &pkgDto.PaginationDataResponse[Author]{
		Items: []Author{
			{BaseModel: models.BaseModel{ID: uuid.New()}, PenName: "Author 1", BirthYear: 1990},
			{BaseModel: models.BaseModel{ID: uuid.New()}, PenName: "Author 2", BirthYear: 1985},
		},
		Pagination: pkgDto.PaginationResponse{
			Page:       1,
			PageSize:   5,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	suite.mockRepo.On("GetAll", suite.ctx, pagination).Return(expectedAuthors, nil)

	authors, code := suite.service.GetAllAuthors(suite.ctx, pagination)

	suite.Equal(dto.Success, code)
	suite.NotNil(authors)
	suite.Equal(expectedAuthors.Items, authors.Items)
	suite.Equal(expectedAuthors.Pagination, authors.Pagination)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllAuthors_EmptyResult() {
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 10}
	expectedAuthors := &pkgDto.PaginationDataResponse[Author]{
		Items: []Author{},
		Pagination: pkgDto.PaginationResponse{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	suite.mockRepo.On("GetAll", suite.ctx, pagination).Return(expectedAuthors, nil)

	authors, code := suite.service.GetAllAuthors(suite.ctx, pagination)

	suite.Equal(dto.Success, code)
	suite.NotNil(authors)
	suite.Empty(authors.Items)
	suite.Equal(expectedAuthors.Pagination, authors.Pagination)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllAuthors_GetAllError() {
	pagination := &pkgDto.PaginationRequest{Page: 1, PageSize: 10}

	suite.mockRepo.On("GetAll", suite.ctx, pagination).Return((*pkgDto.PaginationDataResponse[Author])(nil), errors.New("database error"))

	authors, code := suite.service.GetAllAuthors(suite.ctx, pagination)

	suite.Equal(dto.InternalError, code)
	suite.Nil(authors)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateAuthor_Success() {
	authorID := uuid.New()
	req := &UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	existingAuthor := &Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Original Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return(existingAuthor, nil)
	suite.mockRepo.On("Update", suite.ctx, authorID, mock.AnythingOfType("*author.Author")).Return(nil)

	code := suite.service.UpdateAuthor(suite.ctx, authorID, req)

	suite.Equal(dto.Success, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateAuthor_AuthorNotFound() {
	authorID := uuid.New()
	req := &UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return((*Author)(nil), nil)

	code := suite.service.UpdateAuthor(suite.ctx, authorID, req)

	suite.Equal(dto.AuthorNotFound, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateAuthor_GetByIDError() {
	authorID := uuid.New()
	req := &UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return((*Author)(nil), errors.New("database error"))

	code := suite.service.UpdateAuthor(suite.ctx, authorID, req)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestUpdateAuthor_UpdateError() {
	authorID := uuid.New()
	req := &UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	existingAuthor := &Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Original Author",
		BirthYear: 1990,
	}

	suite.mockRepo.On("GetByID", suite.ctx, authorID).Return(existingAuthor, nil)
	suite.mockRepo.On("Update", suite.ctx, authorID, mock.AnythingOfType("*author.Author")).Return(errors.New("database error"))

	code := suite.service.UpdateAuthor(suite.ctx, authorID, req)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestDeleteAuthor_Success() {
	authorID := uuid.New()

	suite.mockRepo.On("Delete", suite.ctx, authorID).Return(nil)

	code := suite.service.DeleteAuthor(suite.ctx, authorID)

	suite.Equal(dto.Success, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestDeleteAuthor_DeleteError() {
	authorID := uuid.New()

	suite.mockRepo.On("Delete", suite.ctx, authorID).Return(errors.New("database error"))

	code := suite.service.DeleteAuthor(suite.ctx, authorID)

	suite.Equal(dto.InternalError, code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
