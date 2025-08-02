package author

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	"github.com/sirawatc/simple-gin-crud/internal/shared/models"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateAuthor(ctx context.Context, req *CreateAuthorRequest) (*Author, dto.Code) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*Author), args.Get(1).(dto.Code)
}

func (m *MockService) GetAuthorByID(ctx context.Context, id uuid.UUID) (*Author, dto.Code) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*Author), args.Get(1).(dto.Code)
}

func (m *MockService) GetAllAuthors(ctx context.Context, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Author], dto.Code) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*pkgDto.PaginationDataResponse[Author]), args.Get(1).(dto.Code)
}

func (m *MockService) UpdateAuthor(ctx context.Context, id uuid.UUID, req *UpdateAuthorRequest) dto.Code {
	args := m.Called(ctx, id, req)
	return args.Get(0).(dto.Code)
}

func (m *MockService) DeleteAuthor(ctx context.Context, id uuid.UUID) dto.Code {
	args := m.Called(ctx, id)
	return args.Get(0).(dto.Code)
}

type HandlerTestSuite struct {
	suite.Suite
	handler     *Handler
	mockService *MockService
	ctx         context.Context
}

func (suite *HandlerTestSuite) SetupTest() {
	mockService := new(MockService)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	handler := NewHandler(mockService, logger)

	suite.handler = handler
	suite.mockService = mockService
	suite.ctx = context.Background()
}

func (suite *HandlerTestSuite) setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func (suite *HandlerTestSuite) TestNewHandler() {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	suite.NotNil(handler)
	suite.Equal(mockService, handler.service)
	suite.Equal(logger, handler.logger)
}

func (suite *HandlerTestSuite) TestCreateAuthor_Success() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	expectedAuthor := &Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockService.On("CreateAuthor", mock.Anything, &req).Return(expectedAuthor, dto.Success)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/authors", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateAuthor(c)

	responseBody := w.Body.Bytes()
	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusCreated, w.Code)
	suite.Equal(dto.Created, response.Code)
	suite.Equal(expectedAuthor.PenName, response.Data.(map[string]interface{})["penName"])
	suite.Equal(expectedAuthor.BirthYear, int(response.Data.(map[string]interface{})["birthYear"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestCreateAuthor_InvalidJSON() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("POST", "/authors", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestCreateAuthor_BindingError() {
	c, w := suite.setupGinContext()

	req := map[string]interface{}{
		"penName":   "",
		"birthYear": false,
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/authors", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestCreateAuthor_ValidationError() {
	c, w := suite.setupGinContext()

	req := CreateAuthorRequest{
		PenName:   "penName",
		BirthYear: 1000,
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/authors", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
	suite.Equal([]interface{}{"BirthYear must be 1,800 or greater"}, response.Data)
}

func (suite *HandlerTestSuite) TestCreateAuthor_AuthorAlreadyExists() {
	c, w := suite.setupGinContext()

	req := CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockService.On("CreateAuthor", mock.Anything, &req).Return((*Author)(nil), dto.AuthorAlreadyExists)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/authors", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Equal(dto.AuthorAlreadyExists, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestCreateAuthor_ServiceError() {
	c, w := suite.setupGinContext()

	req := CreateAuthorRequest{
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockService.On("CreateAuthor", mock.Anything, &req).Return((*Author)(nil), dto.InternalError)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/authors", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAuthor_Success() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	expectedAuthor := &Author{
		BaseModel: models.BaseModel{ID: authorID},
		PenName:   "Test Author",
		BirthYear: 1990,
	}

	suite.mockService.On("GetAuthorByID", mock.Anything, authorID).Return(expectedAuthor, dto.Success)

	c.Request = httptest.NewRequest("GET", "/authors/"+authorID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.GetAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(expectedAuthor.PenName, response.Data.(map[string]interface{})["penName"])
	suite.Equal(expectedAuthor.BirthYear, int(response.Data.(map[string]interface{})["birthYear"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAuthor_InvalidUUID() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("GET", "/authors/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	suite.handler.GetAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestGetAuthor_NotFound() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()

	suite.mockService.On("GetAuthorByID", mock.Anything, authorID).Return((*Author)(nil), dto.AuthorNotFound)

	c.Request = httptest.NewRequest("GET", "/authors/"+authorID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.GetAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(dto.AuthorNotFound, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAuthor_ServiceError() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()

	suite.mockService.On("GetAuthorByID", mock.Anything, authorID).Return((*Author)(nil), dto.InternalError)

	c.Request = httptest.NewRequest("GET", "/authors/"+authorID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.GetAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAllAuthors_Success() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 5,
	}

	expectedAuthors := &pkgDto.PaginationDataResponse[Author]{
		Items: []Author{
			{BaseModel: models.BaseModel{ID: uuid.New()}, PenName: "Author 1", BirthYear: 1990},
			{BaseModel: models.BaseModel{ID: uuid.New()}, PenName: "Author 2", BirthYear: 1985},
		},
		Pagination: pkgDto.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	suite.mockService.On("GetAllAuthors", mock.Anything, pagination).Return(expectedAuthors, dto.Success)

	url := "/authors?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllAuthors(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(len(expectedAuthors.Items), len(response.Data.(map[string]interface{})["items"].([]interface{})))
	suite.Equal(expectedAuthors.Pagination.Page, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["page"].(float64)))
	suite.Equal(expectedAuthors.Pagination.PageSize, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["pageSize"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAllAuthors_EmptyResult() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	expectedAuthors := &pkgDto.PaginationDataResponse[Author]{
		Items: []Author{},
		Pagination: pkgDto.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	suite.mockService.On("GetAllAuthors", mock.Anything, pagination).Return(expectedAuthors, dto.Success)

	url := "/authors?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllAuthors(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(len(expectedAuthors.Items), len(response.Data.(map[string]interface{})["items"].([]interface{})))
	suite.Equal(expectedAuthors.Pagination.Page, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["page"].(float64)))
	suite.Equal(expectedAuthors.Pagination.PageSize, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["pageSize"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAllAuthors_InvalidPagination() {
	c, w := suite.setupGinContext()

	url := "/authors?page=invalid&pageSize=invalid"
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllAuthors(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
}

func (suite *HandlerTestSuite) TestGetAllAuthors_ServiceError() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	suite.mockService.On("GetAllAuthors", mock.Anything, pagination).Return((*pkgDto.PaginationDataResponse[Author])(nil), dto.InternalError)

	url := "/authors?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllAuthors(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateAuthor_Success() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	suite.mockService.On("UpdateAuthor", mock.Anything, authorID, &req).Return(dto.Success)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/authors/"+authorID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.UpdateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Updated, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateAuthor_InvalidUUID() {
	c, w := suite.setupGinContext()

	req := UpdateAuthorRequest{PenName: "Updated Author", BirthYear: 1985}
	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/authors/invalid-uuid", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	suite.handler.UpdateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestUpdateAuthor_InvalidJSON() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	c.Request = httptest.NewRequest("PUT", "/authors/"+authorID.String(), bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.UpdateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestUpdateAuthor_ValidationError() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := UpdateAuthorRequest{
		PenName:   "penName",
		BirthYear: 1000,
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/authors/"+authorID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.UpdateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
	suite.Equal([]interface{}{"BirthYear must be 1,800 or greater"}, response.Data)
}

func (suite *HandlerTestSuite) TestUpdateAuthor_AuthorNotFound() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	suite.mockService.On("UpdateAuthor", mock.Anything, authorID, &req).Return(dto.AuthorNotFound)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/authors/"+authorID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.UpdateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(dto.AuthorNotFound, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateAuthor_ServiceError() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := UpdateAuthorRequest{
		PenName:   "Updated Author",
		BirthYear: 1985,
	}

	suite.mockService.On("UpdateAuthor", mock.Anything, authorID, &req).Return(dto.InternalError)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/authors/"+authorID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.UpdateAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestDeleteAuthor_Success() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()

	suite.mockService.On("DeleteAuthor", mock.Anything, authorID).Return(dto.Success)

	c.Request = httptest.NewRequest("DELETE", "/authors/"+authorID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.DeleteAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Deleted, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestDeleteAuthor_InvalidUUID() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("DELETE", "/authors/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	suite.handler.DeleteAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestDeleteAuthor_ServiceError() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()

	suite.mockService.On("DeleteAuthor", mock.Anything, authorID).Return(dto.InternalError)

	c.Request = httptest.NewRequest("DELETE", "/authors/"+authorID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: authorID.String()}}

	suite.handler.DeleteAuthor(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
