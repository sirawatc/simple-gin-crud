package author

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
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateAuthor(c *gin.Context) {
	logPrefix := "[AuthorHandler#CreateAuthor]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	var req CreateAuthorRequest
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

	author, code := h.service.CreateAuthor(ctx, &req)
	if code != dto.Success {
		logger.Errorf("%s Failed to create author: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusCreated, dto.BuildBaseResponse(dto.Created, author))
}

func (h *Handler) GetAuthor(c *gin.Context) {
	logPrefix := "[AuthorHandler#GetAuthor]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Errorf("%s Invalid author ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	author, code := h.service.GetAuthorByID(ctx, id)
	if code != dto.Success {
		logger.Errorf("%s Failed to get author: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Success, author))
}

func (h *Handler) GetAllAuthors(c *gin.Context) {
	logPrefix := "[AuthorHandler#GetAllAuthors]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	pagination, errors := pkgDto.NewPaginationRequest(c.Query("page"), c.Query("pageSize"))
	if len(errors) > 0 {
		logger.Errorf("%s Invalid pagination parameters: %v", logPrefix, errors)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.ValidationError, errors))
		return
	}

	authors, code := h.service.GetAllAuthors(ctx, pagination)
	if code != dto.Success {
		logger.Errorf("%s Failed to get all authors: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Success, authors))
}

func (h *Handler) UpdateAuthor(c *gin.Context) {
	logPrefix := "[AuthorHandler#UpdateAuthor]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Errorf("%s Invalid author ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	var req UpdateAuthorRequest
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

	code := h.service.UpdateAuthor(ctx, id, &req)
	if code != dto.Success {
		logger.Errorf("%s Failed to update author: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Updated, nil))
}

func (h *Handler) DeleteAuthor(c *gin.Context) {
	logPrefix := "[AuthorHandler#DeleteAuthor]"

	ctx := c.Request.Context()
	logger := logger.InjectRequestIDWithLogger(ctx, h.logger)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Errorf("%s Invalid author ID format: %v", logPrefix, err)
		c.JSON(http.StatusBadRequest, dto.BuildBaseResponse(dto.UUIDFormatInvalid, nil))
		return
	}

	code := h.service.DeleteAuthor(ctx, id)
	if code != dto.Success {
		logger.Errorf("%s Failed to delete author: %v", logPrefix, dto.CodeMessage[code])
		c.JSON(code.GetHTTPCode(), dto.BuildBaseResponse(code, nil))
		return
	}

	c.JSON(http.StatusOK, dto.BuildBaseResponse(dto.Deleted, nil))
}
