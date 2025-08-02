package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirawatc/simple-gin-crud/internal/shared/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitServer(cfg *config.Config, db *gorm.DB, logger *logrus.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	// Inject security headers if needed ref: https://gin-gonic.com/en/docs/examples/security-headers/

	err := router.SetTrustedProxies(nil)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to set trusted proxies")
	}

	SetupRoutes(router, db, logger)

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Starting server on %s", address)
	if err := router.Run(address); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}

	return router
}
