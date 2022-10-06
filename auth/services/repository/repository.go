package repository

import (
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"gorm.io/gorm"
)

type UserRepositorySQL struct {
	DB *gorm.DB
}

func InitUserRepositorySQL(db *gorm.DB) *UserRepositorySQL {
	return &UserRepositorySQL{DB: db}
}

func (r *UserRepositorySQL) SQLCreateUser(user *models.User) error {
	tx := r.DB.Begin()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *UserRepositorySQL) SQLGetUser(username, password string) (*models.User, error) {
	user := new(models.User)
	err := r.DB.Where("username = ?", username).Where("password = ?", password).First(&user).Error
	return user, err
}

func (r *UserRepositorySQL) SQLUpdatePassword(username, oldpassword, password string) error {
	tx := r.DB.Begin()

	if err := tx.Error; err != nil {
		return err
	}

	result := tx.Model(&models.User{}).Where("username = ?", username).Where("password = ?", oldpassword).Updates(&models.User{Password: password})

	if err := result.Error; err != nil {
		tx.Rollback()
		return err
	}

	if result.RowsAffected < 1 || result.RowsAffected > 1 {
		tx.Rollback()
		return gorm.ErrInvalidTransaction
	}

	return tx.Commit().Error
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
	tx := r.DB.Begin()

	if err := tx.Error; err != nil {
		return err
	}

	result := tx.Where("username = ?", username).Where("password = ?", password).Delete(&models.User{})

	if err := result.Error; err != nil {
		tx.Rollback()
		return err
	}

	if result.RowsAffected < 1 || result.RowsAffected > 1 {
		tx.Rollback()
		return gorm.ErrInvalidTransaction
	}

	return tx.Commit().Error
}
