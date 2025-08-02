package repository

import (
	"gorm.io/gorm"
)

type ITransactionManager interface {
	Transaction(fn func(tx *gorm.DB) error) error
	GetDB(tx ...*gorm.DB) *gorm.DB
}

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) ITransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) Transaction(fn func(tx *gorm.DB) error) error {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (tm *TransactionManager) GetDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return tm.db
}
