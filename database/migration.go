package database

import (
	"github.com/sirawatc/simple-gin-crud/internal/author"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&author.Author{},
	)
}
