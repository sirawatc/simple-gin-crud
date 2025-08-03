package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirawatc/simple-gin-crud/internal/author"
	"github.com/sirawatc/simple-gin-crud/internal/book"
	"github.com/sirawatc/simple-gin-crud/pkg/middleware"
	"github.com/sirawatc/simple-gin-crud/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, logger *logrus.Logger) {
	// Initialize shared dependencies
	transactionManager := repository.NewTransactionManager(db)

	// Initialize repositories
	authorRepo := author.NewRepository(transactionManager, logger)
	bookRepo := book.NewRepository(transactionManager, logger)

	// Initialize services
	authorService := author.NewService(authorRepo, logger)
	bookService := book.NewService(bookRepo, authorService, logger)

	// Initialize handlers
	authorHandler := author.NewHandler(authorService, logger)
	bookHandler := book.NewHandler(bookService, logger)

	// Add middleware
	router.Use(middleware.RequestIDMiddleware())

	// Add cache if needed ref: https://github.com/gin-contrib/cache
	// Add rate limit if needed ref: https://github.com/JGLTechnologies/gin-rate-limit
	initHealthRoutes(router, db)
	initAuthorRoutes(router, authorHandler)
	initBookRoutes(router, bookHandler)
}

func initAuthorRoutes(router *gin.Engine, authorHandler *author.Handler) {
	authors := router.Group("/author")
	{
		authors.POST("/", authorHandler.CreateAuthor)
		authors.GET("/:id", authorHandler.GetAuthor)
		authors.GET("/", authorHandler.GetAllAuthors)
		authors.PUT("/:id", authorHandler.UpdateAuthor)
		authors.DELETE("/:id", authorHandler.DeleteAuthor)
	}
}

func initBookRoutes(router *gin.Engine, bookHandler *book.Handler) {
	books := router.Group("/book")
	{
		books.POST("/", bookHandler.CreateBook)
		books.GET("/:id", bookHandler.GetBook)
		books.GET("/author/:authorId", bookHandler.GetBooksByAuthorID)
		books.GET("/", bookHandler.GetAllBooks)
		books.PUT("/:id", bookHandler.UpdateBook)
		books.DELETE("/:id", bookHandler.DeleteBook)
	}
}

func initHealthRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", func(c *gin.Context) {
		healthMsg := gin.H{
			"status": "ok",
			"checks": gin.H{
				"database": "ok",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}

		dbInstance, err := db.DB()
		if err != nil {
			healthMsg["checks"].(gin.H)["database"] = "down"
			healthMsg["status"] = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, healthMsg)
			return
		}
		err = dbInstance.Ping()
		if err != nil {
			healthMsg["checks"].(gin.H)["database"] = "down"
			healthMsg["status"] = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, healthMsg)
			return
		}

		c.JSON(http.StatusOK, healthMsg)
	})
}
