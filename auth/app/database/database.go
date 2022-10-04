package database

import (
	"fmt"

	"github.com/khuchuz/go-clean-architecture-sql/auth/app/config"
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	return db
}
