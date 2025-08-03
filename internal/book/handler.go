package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/shared/dto"
	pkgDto "github.com/sirawatc/simple-gin-crud/pkg/dto"
	"github.com/sirawatc/simple-gin-crud/pkg/logger"
	"github.com/sirawatc/simple-gin-crud/pkg/validator"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service IService
	logger  *logrus.Logger
}

func NewHandler(service IService, logger *logrus.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) CreateBook(c *gin.Context) {
	logPrefix := "[BookHandler#CreateBook]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("%s Invalid request body: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.BindingError, err.Error()))
		return
	}

	if errors := validator.NewValidator().Validate(req); errors != nil {
		logger.Errorf("%s Validation failed: %v", logPrefix, errors)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.ValidationError, errors))
		return
	}

	book, code := h.service.CreateBook(ctx, &req)
	if code != dto.Success {
		logger.Errorf("%s Failed to create book: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusCreated, dto.BuildBaseResponse(dto.Created, book))
}

func (h *Handler) GetBook(c *gin.Context) {
	logPrefix := "[BookHandler#GetBook]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Errorf("%s Invalid book ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	book, code := h.service.GetBookByID(ctx, id)
	if code != dto.Success {
		logger.Errorf("%s Failed to get book: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Success, book))
}

func (h *Handler) GetBooksByAuthorID(c *gin.Context) {
	logPrefix := "[BookHandler#GetBooksByAuthorID]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	authorID, err := uuid.Parse(c.Param("authorId"))
	if err != nil {
		logger.Errorf("%s Invalid author ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	pagination, errors := pkgDto.NewPaginationRequest(c.Query("page"), c.Query("pageSize"))
	if len(errors) > 0 {
		logger.Errorf("%s Invalid pagination parameters: %v", logPrefix, errors)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.ValidationError, errors))
		return
	}

	books, code := h.service.GetBooksByAuthorID(ctx, authorID, pagination)
	if code != dto.Success {
		logger.Errorf("%s Failed to get books by author ID: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Success, books))
}

func (h *Handler) GetAllBooks(c *gin.Context) {
	logPrefix := "[BookHandler#GetAllBooks]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	pagination, errors := pkgDto.NewPaginationRequest(c.Query("page"), c.Query("pageSize"))
	if len(errors) > 0 {
		logger.Errorf("%s Invalid pagination parameters: %v", logPrefix, errors)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.ValidationError, errors))
		return
	}

	books, code := h.service.GetAllBooks(ctx, pagination)
	if code != dto.Success {
		logger.Errorf("%s Failed to get all books: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Success, books))
}

func (h *Handler) UpdateBook(c *gin.Context) {
	logPrefix := "[BookHandler#UpdateBook]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Errorf("%s Invalid book ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("%s Invalid request body: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.BindingError, err.Error()))
		return
	}

	if errors := validator.NewValidator().Validate(req); errors != nil {
		logger.Errorf("%s Validation failed: %v", logPrefix, errors)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.ValidationError, errors))
		return
	}

	code := h.service.UpdateBook(ctx, id, &req)
	if code != dto.Success {
		logger.Errorf("%s Failed to update book: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Updated, nil))
}

func (h *Handler) DeleteBook(c *gin.Context) {
	logPrefix := "[BookHandler#DeleteBook]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Errorf("%s Invalid book ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	code := h.service.DeleteBook(ctx, id)
	if code != dto.Success {
		logger.Errorf("%s Failed to delete book: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Deleted, nil))
}
