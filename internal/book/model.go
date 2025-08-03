package book

import (
	"github.com/google/uuid"
	"github.com/sirawatc/simple-gin-crud/internal/author"
	"github.com/sirawatc/simple-gin-crud/internal/shared/models"
)

type Book struct {
	models.BaseModel
	AuthorID uuid.UUID `json:"authorId" gorm:"type:uuid;not null;index"`
	Name     string    `json:"name" gorm:"not null"`
	ISBN     string    `json:"isbn" gorm:"not null;unique"`

	Author *author.Author `json:"author" gorm:"foreignKey:AuthorID"`
}
