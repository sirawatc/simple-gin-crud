package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestNewTransactionManager(t *testing.T) {
	gormDB, _ := setupDB(t)

	tm := NewTransactionManager(gormDB)

	assert.NotNil(t, tm)
	assert.IsType(t, &TransactionManager{}, tm)
}

func TestTransactionManager_GetDB(t *testing.T) {
	gormDB, _ := setupDB(t)
	paramDB, _ := setupDB(t)
	assert.NotEqual(t, gormDB, paramDB)

	tm := NewTransactionManager(gormDB)

	tests := []struct {
		name     string
		getTx    func() []*gorm.DB
		expected *gorm.DB
	}{
		{
			name: "no transaction provided",
			getTx: func() []*gorm.DB {
				return nil
			},
			expected: gormDB,
		},
		{
			name: "empty transaction slice",
			getTx: func() []*gorm.DB {
				return []*gorm.DB{}
			},
			expected: gormDB,
		},
		{
			name: "nil transaction",
			getTx: func() []*gorm.DB {
				return []*gorm.DB{nil}
			},
			expected: gormDB,
		},
		{
			name: "valid transaction",
			getTx: func() []*gorm.DB {
				return []*gorm.DB{paramDB}
			},
			expected: paramDB,
		},
		{
			name: "nil transaction with valid transaction",
			getTx: func() []*gorm.DB {
				return []*gorm.DB{nil, paramDB}
			},
			expected: gormDB,
		},
		{
			name: "multiple transactions (should use first)",
			getTx: func() []*gorm.DB {
				paramDB2, _ := setupDB(t)
				return []*gorm.DB{paramDB, paramDB2}
			},
			expected: paramDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.GetDB(tt.getTx()...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransactionManager_Transaction(t *testing.T) {
	gormDB, mock := setupDB(t)

	scenarios := []struct {
		scenarioName string
		tests        []struct {
			name        string
			tx          func(*gorm.DB) error
			expectedErr error
		}
	}{
		{
			scenarioName: "no transaction",
			tests: []struct {
				name        string
				tx          func(*gorm.DB) error
				expectedErr error
			}{
				{
					name: "success",
					tx: func(tx *gorm.DB) error {
						assert.NotNil(t, tx)
						return nil
					},
					expectedErr: nil,
				},
				{
					name: "failed",
					tx: func(tx *gorm.DB) error {
						assert.NotNil(t, tx)
						return errors.New("transaction failed")
					},
					expectedErr: errors.New("transaction failed"),
				},
			},
		},
		{
			scenarioName: "single transaction",
			tests: []struct {
				name        string
				tx          func(*gorm.DB) error
				expectedErr error
			}{
				{
					name: "success",
					tx: func(tx *gorm.DB) error {
						tx.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
						return nil
					},
					expectedErr: nil,
				},
				{
					name: "failed",
					tx: func(tx *gorm.DB) error {
						tx.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
						return errors.New("transaction failed")
					},
					expectedErr: errors.New("transaction failed"),
				},
			},
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.scenarioName, func(t *testing.T) {
			for _, tc := range sc.tests {
				t.Run(tc.name, func(t *testing.T) {
					tm := NewTransactionManager(gormDB)

					mock.ExpectBegin()

					if tc.expectedErr != nil {
						mock.ExpectRollback()
					} else {
						mock.ExpectCommit()
					}

					err := tm.Transaction(tc.tx)
					assert.Equal(t, tc.expectedErr, err)
					assert.NoError(t, mock.ExpectationsWereMet())
				})
			}
		})
	}
}
