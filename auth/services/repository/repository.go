package repository

import (
	"github.com/khuchuz/go-clean-architecture-sql/models"
	"gorm.io/gorm"
)

type UserRepositorySQL struct {
	DB *gorm.DB
}

func InitUserRepositorySQL(db *gorm.DB) *UserRepositorySQL {
	return &UserRepositorySQL{DB: db}
}

func (r *UserRepositorySQL) SQLCreateUser(user *models.User) error {
	err := r.DB.Create(&user).Error
	return err
}

func (r *UserRepositorySQL) SQLGetUser(username, password string) (*models.User, error) {
	user := new(models.User)
	err := r.DB.Where("username = ?", username).Where("password = ?", password).First(&user).Error
	return user, err
}

func (r *UserRepositorySQL) SQLUpdatePassword(username, password string) error {
	err := r.DB.Model(&models.User{}).Where("username = ?", username).Update("password", password).Error
	return err
}

func (r *UserRepositorySQL) SQLIsUserExistByUsername(username string) bool {
	user := new(models.User)
	ret := r.DB.Where("username = ?", username).First(&user)
	return ret.Error == nil
}

func (r *UserRepositorySQL) SQLIsUserExistByEmail(email string) bool {
	user := new(models.User)
	ret := r.DB.Where("email = ?", email).First(&user)
	return ret.Error == nil
}

func (r *UserRepositorySQL) SQLDeleteUser(username, password string) error {
	user := new(models.User)
	err := r.DB.Where("username = ?", username).Where("password = ?", password).First(&user).Error
	if err != nil {
		return err
	}
	err = r.DB.Where("username = ?", username).Where("password = ?", password).Delete(&models.User{}).Error
	return err
}
