package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/khuchuz/go-clean-architecture-sql/auth/app/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSetupDatabase_Success(t *testing.T) {
	db, mocks, err := sqlmock.New()
	assert.NoError(t, err, "Failed to open gorm DB")
	assert.NotNil(t, db, "Mock DB is null")
	assert.NotNil(t, mocks, "SQLMock is null")

	dbmock, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	//dbmock.AutoMigrate(&models.User{})

	dbreal := SetupDatabase()

	assert.NoError(t, err)
	deep.Equal(dbmock, dbreal)
}

func TestSetupDatabase_Panic(t *testing.T) {
	config.DBPort = "1234"
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TestUserFail should have panicked!")
			}
		}()
		// This function should cause a panic
		SetupDatabase()
	}()
}
