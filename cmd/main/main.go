package main

import (
	"os"

	"github.com/sirawatc/simple-gin-crud/database"
	"github.com/sirawatc/simple-gin-crud/internal/shared/config"
	"github.com/sirawatc/simple-gin-crud/pkg/logger"
	"github.com/sirawatc/simple-gin-crud/server"
)

func main() {
	cfg := config.NewConfig()

	logger := logger.NewLogger(cfg.ServiceName)

	db, err := database.NewPostgres(cfg)
	if err != nil {
		logger.Errorf("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	if cfg.Database.AutoMigrate {
		if err = database.Migrate(db); err != nil {
			logger.Errorf("Failed to migrate database: %v", err)
			os.Exit(1)
		}
	}

	server.InitServer(cfg, db, logger)
}
