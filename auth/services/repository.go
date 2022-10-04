package services

import (
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
)

type UserRepositorySQL interface {
	SQLCreateUser(user *models.User) error
	SQLGetUser(username, password string) (*models.User, error)
	SQLIsUserExistByUsername(username string) bool
	SQLIsUserExistByEmail(email string) bool
	SQLUpdatePassword(username, oldpassword, password string) error
	SQLDeleteUser(username, password string) error
}
