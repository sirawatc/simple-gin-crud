package author

import (
	"github.com/sirawatc/simple-gin-crud/internal/shared/models"
)

type Author struct {
	models.BaseModel
	PenName   string `json:"penName" gorm:"not null;unique"`
	BirthYear int    `json:"birthYear" gorm:"not null"`
}
