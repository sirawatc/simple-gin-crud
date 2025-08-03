package book

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

func (m *MockService) CreateBook(ctx context.Context, req *CreateBookRequest) (*Book, dto.Code) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*Book), args.Get(1).(dto.Code)
}

func (m *MockService) GetBookByID(ctx context.Context, id uuid.UUID) (*Book, dto.Code) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*Book), args.Get(1).(dto.Code)
}

func (m *MockService) GetBooksByAuthorID(ctx context.Context, authorID uuid.UUID, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Book], dto.Code) {
	args := m.Called(ctx, authorID, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*pkgDto.PaginationDataResponse[Book]), args.Get(1).(dto.Code)
}

func (m *MockService) GetAllBooks(ctx context.Context, pagination *pkgDto.PaginationRequest) (*pkgDto.PaginationDataResponse[Book], dto.Code) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.Code)
	}
	return args.Get(0).(*pkgDto.PaginationDataResponse[Book]), args.Get(1).(dto.Code)
}

func (m *MockService) UpdateBook(ctx context.Context, id uuid.UUID, req *UpdateBookRequest) dto.Code {
	args := m.Called(ctx, id, req)
	return args.Get(0).(dto.Code)
}

func (m *MockService) DeleteBook(ctx context.Context, id uuid.UUID) dto.Code {
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

func (suite *HandlerTestSuite) TestCreateBook_Success() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	authorID := uuid.New()
	req := CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	expectedBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Test Book",
		ISBN:      "978-0-7475-3269-9",
	}

	suite.mockService.On("CreateBook", mock.Anything, &req).Return(expectedBook, dto.Success)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()
	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusCreated, w.Code)
	suite.Equal(dto.Created, response.Code)
	suite.Equal(expectedBook.Name, response.Data.(map[string]interface{})["name"])
	suite.Equal(expectedBook.ISBN, response.Data.(map[string]interface{})["isbn"])
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestCreateBook_InvalidJSON() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestCreateBook_BindingError() {
	c, w := suite.setupGinContext()

	req := map[string]interface{}{
		"authorId": "",
		"name":     "Test Book",
		"isbn":     false,
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestCreateBook_ValidationError() {
	c, w := suite.setupGinContext()

	req := CreateBookRequest{
		AuthorID: uuid.New(),
		Name:     "name",
		ISBN:     "978-0-7475-3269",
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
}

func (suite *HandlerTestSuite) TestCreateBook_AuthorNotFound() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("CreateBook", mock.Anything, &req).Return((*Book)(nil), dto.AuthorNotFound)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(dto.AuthorNotFound, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestCreateBook_BookAlreadyExists() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("CreateBook", mock.Anything, &req).Return((*Book)(nil), dto.BookAlreadyExists)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Equal(dto.BookAlreadyExists, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestCreateBook_ServiceError() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()
	req := CreateBookRequest{
		AuthorID: authorID,
		Name:     "Test Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("CreateBook", mock.Anything, &req).Return((*Book)(nil), dto.InternalError)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.CreateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetBook_Success() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	authorID := uuid.New()
	expectedBook := &Book{
		BaseModel: models.BaseModel{ID: bookID},
		AuthorID:  authorID,
		Name:      "Test Book",
		ISBN:      "1234567890123",
	}

	suite.mockService.On("GetBookByID", mock.Anything, bookID).Return(expectedBook, dto.Success)

	c.Request = httptest.NewRequest("GET", "/books/"+bookID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.GetBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(expectedBook.Name, response.Data.(map[string]interface{})["name"])
	suite.Equal(expectedBook.ISBN, response.Data.(map[string]interface{})["isbn"])
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetBook_InvalidUUID() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("GET", "/books/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	suite.handler.GetBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestGetBook_NotFound() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()

	suite.mockService.On("GetBookByID", mock.Anything, bookID).Return((*Book)(nil), dto.BookNotFound)

	c.Request = httptest.NewRequest("GET", "/books/"+bookID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.GetBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(dto.BookNotFound, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetBook_ServiceError() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()

	suite.mockService.On("GetBookByID", mock.Anything, bookID).Return((*Book)(nil), dto.InternalError)

	c.Request = httptest.NewRequest("GET", "/books/"+bookID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.GetBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAllBooks_Success() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{
			{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "Book 1", ISBN: "1234567890123"},
			{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "Book 2", ISBN: "1234567890124"},
		},
		Pagination: pkgDto.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	suite.mockService.On("GetAllBooks", mock.Anything, pagination).Return(expectedBooks, dto.Success)

	url := "/books?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllBooks(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(len(expectedBooks.Items), len(response.Data.(map[string]interface{})["items"].([]interface{})))
	suite.Equal(expectedBooks.Pagination.Page, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["page"].(float64)))
	suite.Equal(expectedBooks.Pagination.PageSize, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["pageSize"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAllBooks_EmptyResult() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{},
		Pagination: pkgDto.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	suite.mockService.On("GetAllBooks", mock.Anything, pagination).Return(expectedBooks, dto.Success)

	url := "/books?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllBooks(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(len(expectedBooks.Items), len(response.Data.(map[string]interface{})["items"].([]interface{})))
	suite.Equal(expectedBooks.Pagination.Page, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["page"].(float64)))
	suite.Equal(expectedBooks.Pagination.PageSize, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["pageSize"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetAllBooks_InvalidPagination() {
	c, w := suite.setupGinContext()

	url := "/books?page=invalid&pageSize=invalid"
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllBooks(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
}

func (suite *HandlerTestSuite) TestGetAllBooks_ServiceError() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	suite.mockService.On("GetAllBooks", mock.Anything, pagination).Return((*pkgDto.PaginationDataResponse[Book])(nil), dto.InternalError)

	url := "/books?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)

	suite.handler.GetAllBooks(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetBooksByAuthorID_Success() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	authorID := uuid.New()
	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{
			{BaseModel: models.BaseModel{ID: uuid.New()}, AuthorID: authorID, Name: "Book 1", ISBN: "1234567890123"},
			{BaseModel: models.BaseModel{ID: uuid.New()}, AuthorID: authorID, Name: "Book 2", ISBN: "1234567890124"},
		},
		Pagination: pkgDto.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	suite.mockService.On("GetBooksByAuthorID", mock.Anything, authorID, pagination).Return(expectedBooks, dto.Success)

	url := "/authors/" + authorID.String() + "/books?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)
	c.Params = gin.Params{{Key: "authorId", Value: authorID.String()}}

	suite.handler.GetBooksByAuthorID(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(len(expectedBooks.Items), len(response.Data.(map[string]interface{})["items"].([]interface{})))
	suite.Equal(expectedBooks.Pagination.Page, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["page"].(float64)))
	suite.Equal(expectedBooks.Pagination.PageSize, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["pageSize"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetBooksByAuthorID_EmptyResult() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	authorID := uuid.New()
	expectedBooks := &pkgDto.PaginationDataResponse[Book]{
		Items: []Book{},
		Pagination: pkgDto.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	suite.mockService.On("GetBooksByAuthorID", mock.Anything, authorID, pagination).Return(expectedBooks, dto.Success)

	url := "/authors/" + authorID.String() + "/books?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)
	c.Params = gin.Params{{Key: "authorId", Value: authorID.String()}}

	suite.handler.GetBooksByAuthorID(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Success, response.Code)
	suite.Equal(len(expectedBooks.Items), len(response.Data.(map[string]interface{})["items"].([]interface{})))
	suite.Equal(expectedBooks.Pagination.Page, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["page"].(float64)))
	suite.Equal(expectedBooks.Pagination.PageSize, int(response.Data.(map[string]interface{})["pagination"].(map[string]interface{})["pageSize"].(float64)))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetBooksByAuthorID_InvalidAuthorUUID() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("GET", "/authors/invalid-uuid/books", nil)
	c.Params = gin.Params{{Key: "authorId", Value: "invalid-uuid"}}

	suite.handler.GetBooksByAuthorID(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestGetBooksByAuthorID_InvalidPagination() {
	c, w := suite.setupGinContext()

	authorID := uuid.New()

	url := "/authors/" + authorID.String() + "/books?page=invalid&pageSize=invalid"
	c.Request = httptest.NewRequest("GET", url, nil)
	c.Params = gin.Params{{Key: "authorId", Value: authorID.String()}}

	suite.handler.GetBooksByAuthorID(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
}

func (suite *HandlerTestSuite) TestGetBooksByAuthorID_ServiceError() {
	c, w := suite.setupGinContext()
	pagination := &pkgDto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	authorID := uuid.New()

	suite.mockService.On("GetBooksByAuthorID", mock.Anything, authorID, pagination).Return((*pkgDto.PaginationDataResponse[Book])(nil), dto.InternalError)

	url := "/authors/" + authorID.String() + "/books?page=" + strconv.Itoa(pagination.Page) + "&pageSize=" + strconv.Itoa(pagination.PageSize)
	c.Request = httptest.NewRequest("GET", url, nil)
	c.Params = gin.Params{{Key: "authorId", Value: authorID.String()}}

	suite.handler.GetBooksByAuthorID(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateBook_Success() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	authorID := uuid.New()
	req := UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("UpdateBook", mock.Anything, bookID, &req).Return(dto.Success)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Updated, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateBook_InvalidUUID() {
	c, w := suite.setupGinContext()

	req := UpdateBookRequest{Name: "Updated Book", ISBN: "1234567890123"}
	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/invalid-uuid", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestUpdateBook_InvalidJSON() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestUpdateBook_BindingError() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	req := map[string]interface{}{
		"authorId": "",
		"name":     "Test Book",
		"isbn":     false,
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.BindingError, response.Code)
}

func (suite *HandlerTestSuite) TestUpdateBook_ValidationError() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	req := UpdateBookRequest{
		AuthorID: uuid.New(),
		Name:     "name",
		ISBN:     "978-0-7475-3269",
	}

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.ValidationError, response.Code)
}

func (suite *HandlerTestSuite) TestUpdateBook_BookNotFound() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	authorID := uuid.New()
	req := UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("UpdateBook", mock.Anything, bookID, &req).Return(dto.BookNotFound)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(dto.BookNotFound, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateBook_AuthorNotFound() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	authorID := uuid.New()
	req := UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("UpdateBook", mock.Anything, bookID, &req).Return(dto.AuthorNotFound)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(dto.AuthorNotFound, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestUpdateBook_ServiceError() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()
	authorID := uuid.New()
	req := UpdateBookRequest{
		AuthorID: authorID,
		Name:     "Updated Book",
		ISBN:     "978-0-7475-3269-9",
	}

	suite.mockService.On("UpdateBook", mock.Anything, bookID, &req).Return(dto.InternalError)

	reqBody, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.UpdateBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Equal(dto.InternalError, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestDeleteBook_Success() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()

	suite.mockService.On("DeleteBook", mock.Anything, bookID).Return(dto.Success)

	c.Request = httptest.NewRequest("DELETE", "/books/"+bookID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.DeleteBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(dto.Deleted, response.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestDeleteBook_InvalidUUID() {
	c, w := suite.setupGinContext()

	c.Request = httptest.NewRequest("DELETE", "/books/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

	suite.handler.DeleteBook(c)

	responseBody := w.Body.Bytes()

	var response dto.BaseResponse
	err := json.Unmarshal(responseBody, &response)
	suite.NoError(err)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(dto.UUIDFormatInvalid, response.Code)
}

func (suite *HandlerTestSuite) TestDeleteBook_ServiceError() {
	c, w := suite.setupGinContext()

	bookID := uuid.New()

	suite.mockService.On("DeleteBook", mock.Anything, bookID).Return(dto.InternalError)

	c.Request = httptest.NewRequest("DELETE", "/books/"+bookID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: bookID.String()}}

	suite.handler.DeleteBook(c)

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
